package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "SNAPSHOT"

func main() {
	if err := app().Run(context.Background(), os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func app() *cli.Command {
	return &cli.Command{
		Name:      "redact-ext-1password",
		Usage:     "1Password store extension for react key exchange",
		ArgsUsage: "key=op://<vault>/<item>/[<section>/]<field>",
		Version:   version,
		Commands: []*cli.Command{
			cmdList(),
			cmdGet(),
			cmdPut(),
		},
	}
}
