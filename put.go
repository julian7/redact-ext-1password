package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/1password/onepassword-sdk-go"
	"github.com/urfave/cli/v3"
)

func cmdPut() *cli.Command {
	return &cli.Command{
		Name:        "put",
		Usage:       "Put secret to 1Password",
		ArgsUsage:   "key=op://<vault>/<item>/[<section>/]<field>",
		Description: "Writes 1Password secret from STDOUT into the specified field",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config, err := loadConfig(cmd.Args().Slice())
			if err != nil {
				return err
			}

			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}

			client, err := newOnepasswordClient(ctx, cmd)
			if err != nil {
				return fmt.Errorf("cannot create 1password client: %w", err)
			}

			vaultName, itemName, sectionName, fieldName, err := getVaultItemSectionField(config.key)
			if err != nil {
				return err
			}
			var vaultID, itemID string

			vaults, err := client.Vaults().List(ctx)
			if err != nil {
				return fmt.Errorf("cannot list vaults: %w", err)
			}
			for _, vault := range vaults {
				if vault.Title == vaultName {
					vaultID = vault.ID
					break
				}
			}
			if vaultID == "" {
				return fmt.Errorf("%q: %w", vaultName, ErrVaultNotFound)
			}
			items, err := client.Items().List(ctx, vaultID)
			if err != nil {
				return fmt.Errorf("cannot list items: %w", err)
			}
			for _, item := range items {
				if item.Title == itemName {
					itemID = item.ID
					break
				}
			}
			secret := newSecret(ctx, client.Items())

			if itemID == "" {
				params := onepassword.ItemCreateParams{
					Title:    itemName,
					Category: onepassword.ItemCategorySecureNote,
					VaultID:  vaultID,
					Fields:   []onepassword.ItemField{{ID: fieldName, Title: fieldName, Value: string(data), FieldType: onepassword.ItemFieldTypeConcealed, SectionID: DefaultSectionID}},
					Sections: []onepassword.ItemSection{{ID: "add more", Title: ""}},
				}
				if sectionName != "" {
					params.Sections = []onepassword.ItemSection{{ID: sectionName, Title: sectionName}}
					params.Fields[0].SectionID = &sectionName
				}

				if err := secret.create(params); err != nil {
					return fmt.Errorf("cannot create item %s: %w", itemName, err)
				}
			} else {
				if err := secret.get(vaultID, itemID); err != nil {
					return fmt.Errorf("cannot get item %s: %w", itemName, err)
				}
			}

			defer secret.save()
			if sectionName == "" && fieldName == "Notes" {
				secret.setNotes(string(data))
				return nil
			}

			sectionID := secret.sectionID(sectionName)

			secret.writeField(fieldName, sectionID, string(data))
			return nil
		},
	}
}
