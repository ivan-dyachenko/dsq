package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func describeProperties(c *cli.Context) error {
	targetKind := c.String("kind")
	if targetKind != "" {
		projectID := c.String("project")
		client, ctx, err := createClient(projectID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		props, err := properties(ctx, client, targetKind)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Type", "Indexed"})

		for _, p := range props.Names {
			prop := props.Props[p]
			typeName := fmt.Sprintf("%T", prop.Value)
			indexed := ""
			if !prop.NoIndex {
				indexed = "indexed"
			}
			table.Append([]string{prop.Name, typeName, indexed})
		}

		table.Render()
	}
	return nil
}
