package main

import (
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func getKinds(c *cli.Context) error {
	projectID := c.String("project")
	namespace := c.String("namespace")
	showMeta := c.Bool("metadata")

	query := datastore.NewQuery("__kind__").Namespace(namespace).KeysOnly()
	client, ctx, err := createClient(projectID)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Kind", "Namespace"})

	for _, k := range keys {
		if showMeta || !strings.Contains(k.Name, "__") {
			table.Append([]string{k.Name, k.Kind, k.Namespace})
		}
	}

	table.Render()
	return nil
}

func getNS(c *cli.Context) error {
	projectID := c.String("project")

	query := datastore.NewQuery("__namespace__").KeysOnly()
	client, ctx, err := createClient(projectID)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name"})

	for _, k := range keys {
		table.Append([]string{k.Name})
	}

	table.Render()
	return nil
}

func getEntities(c *cli.Context) error {
	projectID := c.String("project")
	targetKind := c.String("kind")

	client, ctx, err := createClient(projectID)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	propsMap, err := properties(ctx, client, targetKind)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCells(false)
	table.SetRowLine(false)

	header := []string{"ID"}
	headerKes := []string{}
	for _, name := range propsMap.Names {
		header = append(header, name)
		headerKes = append(headerKes, name)
	}
	table.SetHeader(header)

	query := datastore.NewQuery(targetKind).Limit(10)
	var entities []*RawProperties
	keys, err := client.GetAll(ctx, query, &entities)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	for i, k := range keys {
		entity := entities[i]
		row := make([]string, len(headerKes)+1)
		id := string(k.ID)
		if k.Name != "" {
			id = k.Name
		}
		row[0] = id
		for i, prop := range headerKes {
			if p, ok := entity.Props[prop]; ok {
				propStr := fmt.Sprintf("%v", p.Value)
				if len(propStr) > 25 {
					propStr = fmt.Sprintf("%s...", propStr[:24])
				}
				row[i+1] = propStr
			}

		}

		table.Append(row)
	}

	table.Render()

	return nil
}
