package store

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mastercactapus/quicklink/pkg/store/pgstore"
)

type Postgres struct {
	s *pgstore.Queries
}

var _ Store = (*Postgres)(nil)

//go:embed pg-schema.sql
var schema string

func NewPostgres(url string) (*Postgres, error) {
	p, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	_, err = p.Exec(context.Background(), schema)
	if err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &Postgres{
		s: pgstore.New(p),
	}, nil
}

func (p *Postgres) Get(ctx context.Context, key string) (string, error) {
	return p.s.FindLink(ctx, key)
}

func (p *Postgres) Set(ctx context.Context, key, value string) error {
	if value == "" {
		return p.s.DeleteLink(ctx, key)
	}

	return p.s.UpsertLink(ctx, pgstore.UpsertLinkParams{
		Base: key,
		Url:  value,
	})
}

func (p *Postgres) Scanner(ctx context.Context) (Scanner, error) {
	rows, err := p.s.AllLinks(ctx)
	if err != nil {
		return nil, err
	}

	var results []entry
	for _, r := range rows {
		results = append(results, entry{r.Base, r.Url})
	}

	return &entryScanner{
		results: results,
	}, nil
}
