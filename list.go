package main

import (
	"context"
	"fmt"

	"github.com/1password/onepassword-sdk-go"
	"github.com/urfave/cli/v3"
)

func cmdList() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "Shows 1Password configuration",
		ArgsUsage: "key=op://<vault>/<item>/[<section>/]<field>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config, err := loadConfig(cmd.Args().Slice())
			if err != nil {
				return err
			}

			if err := onepassword.Secrets.ValidateSecretReference(ctx, config.key); err != nil {
				return err
			}
			fmt.Printf("1Password key %s\n", config.key)

			return nil
		},
	}
}
