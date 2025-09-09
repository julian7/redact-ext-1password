package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/1password/onepassword-sdk-go"
	"github.com/urfave/cli/v3"
)

var SECTION_ID string = "add more"
var DefaultSectionID *string = &SECTION_ID

func newOnepasswordClient(ctx context.Context, cmd *cli.Command) (*onepassword.Client, error) {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")
	client, err := onepassword.NewClient(ctx, onepassword.WithServiceAccountToken(token), onepassword.WithIntegrationInfo(cmd.Root().Name, version))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getVaultItemSectionField(secref string) (string, string, string, string, error) {
	url, err := url.Parse(secref)
	if err != nil {
		return "", "", "", "", err
	}

	if url.Scheme != "op" {
		return "", "", "", "", fmt.Errorf("%q: %w", url.Scheme, ErrInvalidScheme)
	}

	pathItems := strings.Split(url.Path, "/")
	switch pathLen := len(pathItems); {
	case pathLen < 3:
		return "", "", "", "", ErrSecRefTooShort
	case pathLen == 3:
		return url.Host, pathItems[1], "", pathItems[2], nil
	case pathLen == 4:
		return url.Host, pathItems[1], pathItems[2], pathItems[3], nil
	}
	return "", "", "", "", ErrSecRefTooLong
}

type secret struct {
	onepassword.ItemsAPI
	ctx   context.Context
	item  *onepassword.Item
	dirty bool
}

func newSecret(ctx context.Context, api onepassword.ItemsAPI) *secret {
	return &secret{
		ctx:      ctx,
		ItemsAPI: api,
		dirty:    false,
	}
}

func (s *secret) create(params onepassword.ItemCreateParams) error {
	item, err := s.ItemsAPI.Create(s.ctx, params)
	if err != nil {
		return err
	}
	s.item = &item
	return nil
}

func (s *secret) get(vaultID, itemID string) error {
	item, err := s.ItemsAPI.Get(s.ctx, vaultID, itemID)
	if err != nil {
		return err
	}
	s.item = &item
	return nil
}

func (s *secret) setNotes(notes string) {
	s.item.Notes = notes
	s.dirty = true
}

func (s *secret) sectionID(sectionName string) *string {

	if sectionName == "" {
		return DefaultSectionID
	}

	for _, section := range s.item.Sections {
		if section.Title == sectionName {
			return &section.ID
		}
	}
	s.item.Sections = append(s.item.Sections, onepassword.ItemSection{ID: sectionName, Title: sectionName})
	return &sectionName
}

func (s *secret) writeField(fieldName string, sectionID *string, data string) {
	for _, field := range s.item.Fields {
		if field.Title == fieldName {
			if sectionID != nil && field.SectionID != nil && *field.SectionID != *sectionID {
				continue
			}
			field.Value = data
			s.dirty = true
			return
		}
	}

	newField := onepassword.ItemField{
		ID:        fieldName,
		Title:     fieldName,
		Value:     string(data),
		FieldType: onepassword.ItemFieldTypeConcealed,
	}
	if sectionID != nil {
		newField.SectionID = sectionID
	}
	s.item.Fields = append(s.item.Fields, newField)
	s.dirty = true
}

func (s *secret) save() error {
	if !s.dirty {
		return nil
	}
	item, err := s.ItemsAPI.Put(s.ctx, *s.item)
	if err != nil {
		return err
	}
	s.item = &item
	s.dirty = false
	return nil
}
