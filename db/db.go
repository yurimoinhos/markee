package db

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/ent"
	_ "github.com/aggi-tech/aggipay/ent/runtime"
	_ "github.com/lib/pq"
)

func NewClient(ctx context.Context, databaseURL string) (*ent.Client, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL nao informado")
	}

	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return client, nil
}
