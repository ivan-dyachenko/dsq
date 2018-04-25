package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

var (
	app *cli.App
)

func main() {
	app = cli.NewApp()
	app.Name = "dsq"
	app.Usage = "Google Datastore Queries"
	app.Version = "v0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Ivan Diachenko",
			Email: "ivan.dyachenko@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		c.App.Setup()
		return nil
	}

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "namespace, ns",
			Usage: "Namespace",
		},
		cli.StringFlag{
			Name:  "kind, k",
			Usage: "Kind",
		},
		cli.BoolFlag{
			Name:  "metadata, meta",
			Usage: "Inlude Metadata queries",
		},
		cli.StringFlag{
			Name:   "project, p",
			Usage:  "You need to set the environment variable \"DATASTORE_PROJECT_ID\"",
			EnvVar: "DATASTORE_PROJECT_ID",
		},
	}

	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{},
			Usage:   "Display one or many resources",
			Flags:   flags,
			Action: func(c *cli.Context) error {
				fmt.Printf("~~ not found %v", c.Args())
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:   "kinds",
					Usage:  "List all kinds in ps output format.",
					Action: getKinds,
					Flags:  flags,
				},
				{
					Name:    "namespace",
					Aliases: []string{"ns"},
					Usage:   "List all namespaces in ps output format.",
					Action:  getNS,
					Flags:   flags,
				},
				{
					Name:    "entities",
					Aliases: []string{"es"},
					Usage:   "List all entities in ps output format.",
					Action:  getEntities,
					Flags:   flags,
				},
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe one or many resources",
			Flags:   flags,
			Action: func(c *cli.Context) error {
				fmt.Printf("-- not found %v", c.Args())
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "properties",
					Aliases: []string{"ps"},
					Usage:   "Describe Kind's properties in ps output format.",
					Action:  describeProperties,
					Flags:   flags,
				},
			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Thar be no %q here.\n", command)
	}
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			fmt.Fprintf(c.App.Writer, "SUB WRONG: %#v\n", err)
			return err
		}

		fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
		return nil
	}

	app.Run(os.Args)
}
