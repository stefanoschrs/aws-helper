package main

import (
	"github.com/stefanoschrs/aws-helper/internal"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var version = "development"

func main() {
	app := &cli.App{
		Usage:   "Helper functions for common aws actions",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: ".env",
				Usage: "env configuration file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "invalidate",
				Usage:       "invalidate <name>",
				Description: "invalidate a CloudFront distribution's cache",
				Action:      internal.ActionInvalidate,
			},
			{
				Name: "ecs",
				Subcommands: []*cli.Command{
					{
						Name:        "deploy",
						Usage:       "deploy <name>",
						Description: "deploy a new ECS service by creating a new revision and updating the running service",
						Action:      internal.ActionECSDeploy,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
