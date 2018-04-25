package main

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/datastore"
	"github.com/urfave/cli"
)

// Prop to hold property_representation
type Prop struct {
	K    *datastore.Key
	Repr []string `datastore:"property_representation"`
}

// RawProperties wrapper around []datastore.Property
type RawProperties struct {
	Names []string
	Props map[string]datastore.Property
}

// Load implements implements PropertyLoadSaver.Load
func (x *RawProperties) Load(ps []datastore.Property) error {
	propsMap := make(map[string]datastore.Property)
	names := []string{}
	for _, p := range ps {
		propsMap[p.Name] = p
		names = append(names, p.Name)
	}

	sort.Strings(names)

	x.Props = propsMap
	x.Names = names
	return nil
}

// Save implements implements PropertyLoadSaver.Save
func (x *RawProperties) Save() ([]datastore.Property, error) {
	return []datastore.Property{}, nil
}

func kindProp(ctx context.Context, client *datastore.Client, targetKind string, orderBy string) (*RawProperties, error) {
	query := datastore.NewQuery(targetKind).Order(orderBy).Limit(1)
	var entities []*RawProperties
	_, err := client.GetAll(ctx, query, &entities)
	if err != nil {
		return nil, cli.NewExitError(err, 1)
	}

	if len(entities) == 0 {
		return nil, fmt.Errorf("%s kind doesn't have entities", targetKind)
	}

	raw := entities[0]
	return raw, nil
}

func properties(ctx context.Context, client *datastore.Client, targetKind string) (*RawProperties, error) {
	raw, err := kindProp(ctx, client, targetKind, "__key__")
	return raw, err
}

func metaProperties(ctx context.Context, client *datastore.Client, targetKind string) (map[string]*Prop, error) {
	kindKey := datastore.NameKey("__kind__", targetKind, nil)
	query := datastore.NewQuery("__property__").Ancestor(kindKey)

	var props []*Prop
	keys, err := client.GetAll(ctx, query, &props)

	if err != nil {
		return nil, cli.NewExitError(err, 1)
	}

	var propsMap = make(map[string]*Prop)
	for i, k := range keys {
		props[i].K = k
		propsMap[k.Name] = props[i]
	}
	return propsMap, err
}
