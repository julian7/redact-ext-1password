package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"
)

func cmdPut() *cli.Command {
	return &cli.Command{
		Name:        "put",
		Usage:       "Put secret to 1Password",
		ArgsUsage:   "key=op://<vault>/<item>/[<section>/]<field>",
		Description: "Writes 1Password secret from STDOUT into the specified field",
		Action:      actionRun,
	}
}

func actionRun(ctx context.Context, cmd *cli.Command) (err error) {
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

	secret := newSecret(ctx, client)

	vaultID, err = secret.vaultIDbyName(vaultName)
	if err != nil {
		return fmt.Errorf("searching for vault %q: %w", vaultName, err)
	}

	itemID, err = secret.itemIDbyName(vaultID, itemName)
	if err != nil {
		if !errors.Is(err, ErrItemNotFound) {
			return fmt.Errorf("cannot list items: %w", err)
		}

		if err := secret.create(vaultID, itemName, sectionName, fieldName, string(data)); err != nil {
			return fmt.Errorf("cannot create item %s: %w", itemName, err)
		}
	}

	if itemID != "" {
		if err := secret.get(vaultID, itemID); err != nil {
			return fmt.Errorf("cannot get item %s: %w", itemName, err)
		}
	}

	defer func() {
		err = secret.save()
	}()

	if sectionName == "" && fieldName == "Notes" {
		secret.setNotes(string(data))

		return err
	}

	secret.writeField(fieldName, secret.sectionID(sectionName), string(data))

	return err
}
