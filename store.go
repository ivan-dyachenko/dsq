package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

func createClient(projectID string) (*datastore.Client, context.Context, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		return nil, nil, fmt.Errorf("Could not create datastore client: %v", err)
	}
	return client, ctx, nil
}
