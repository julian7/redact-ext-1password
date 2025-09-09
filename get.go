package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func cmdGet() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "Get secret from 1Password",
		ArgsUsage:   "key=op://<vault>/<item>/[<section>/]<field>",
		Description: "Prints 1Password secret to STDOUT",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config, err := loadConfig(cmd.Args().Slice())
			if err != nil {
				return err
			}
			client, err := newOnepasswordClient(ctx, cmd)
			if err != nil {
				return err
			}

			secret, err := client.Secrets().Resolve(ctx, config.key)
			if err != nil {
				return fmt.Errorf("cannot resolve secret: %w", err)
			}
			fmt.Println(secret)

			return nil
		},
	}
}
