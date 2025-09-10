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

var DefaultSectionIDStr = "add more"
var DefaultSectionID = &DefaultSectionIDStr

func newOnepasswordClient(ctx context.Context, cmd *cli.Command) (*onepassword.Client, error) {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo(cmd.Root().Name, version),
	)
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
	*onepassword.Client
	ctx   context.Context
	item  *onepassword.Item
	dirty bool
}

func newSecret(ctx context.Context, client *onepassword.Client) *secret {
	return &secret{
		ctx:    ctx,
		Client: client,
		dirty:  false,
	}
}

func (s *secret) create(vaultID, itemName, sectionName, fieldName, data string) error {
	params := onepassword.ItemCreateParams{
		Title:    itemName,
		Category: onepassword.ItemCategorySecureNote,
		VaultID:  vaultID,
		Fields: []onepassword.ItemField{
			{
				ID:        fieldName,
				Title:     fieldName,
				Value:     data,
				FieldType: onepassword.ItemFieldTypeConcealed,
				SectionID: DefaultSectionID,
			},
		},
		Sections: []onepassword.ItemSection{{ID: "add more", Title: ""}},
	}
	if sectionName != "" {
		params.Sections = []onepassword.ItemSection{{ID: sectionName, Title: sectionName}}
		params.Fields[0].SectionID = &sectionName
	}

	item, err := s.Items().Create(s.ctx, params)
	if err != nil {
		return err
	}

	s.item = &item

	return nil
}

func (s *secret) get(vaultID, itemID string) error {
	item, err := s.Items().Get(s.ctx, vaultID, itemID)
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
		Value:     data,
		FieldType: onepassword.ItemFieldTypeConcealed,
	}
	if sectionID != nil {
		newField.SectionID = sectionID
	}

	s.item.Fields = append(s.item.Fields, newField)
	s.dirty = true
}

func (s *secret) vaultIDbyName(vaultName string) (string, error) {
	vaults, err := s.Vaults().List(s.ctx)
	if err != nil {
		return "", err
	}

	for _, vault := range vaults {
		if vault.Title == vaultName {
			return vault.ID, nil
		}
	}

	return "", ErrVaultNotFound
}

func (s *secret) itemIDbyName(vaultID, itemName string) (string, error) {
	items, err := s.Items().List(s.ctx, vaultID)
	if err != nil {
		return "", err
	}

	for _, item := range items {
		if item.Title == itemName {
			return item.ID, nil
		}
	}

	return "", ErrItemNotFound
}

func (s *secret) save() error {
	if !s.dirty {
		return nil
	}

	item, err := s.Items().Put(s.ctx, *s.item)
	if err != nil {
		return err
	}

	s.item = &item
	s.dirty = false

	return nil
}
