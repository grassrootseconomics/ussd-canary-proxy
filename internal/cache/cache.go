package cache

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Cache interface {
		Get(context.Context, string) (int, error)
		Set(context.Context, string, int) error
		Update(context.Context, string, int) error
	}

	PgCacheOpts struct {
		DSN                  string
		MigrationsFolderPath string
		QueriesFolderPath    string
	}

	PgCache struct {
		db      *pgxpool.Pool
		queries *queries
	}

	queries struct {
		GetFromCache string `query:"get-cache"`
		SetToCache   string `query:"set-cache"`
		UpdateCache  string `query:"update-cache"`
	}
)

func NewPgCache(o PgCacheOpts) (Cache, error) {
	parsedConfig, err := pgxpool.ParseConfig(o.DSN)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), parsedConfig)
	if err != nil {
		return nil, err
	}

	queries, err := loadQueries(o.QueriesFolderPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(context.Background(), dbPool, o.MigrationsFolderPath); err != nil {
		return nil, err
	}

	return &PgCache{
		db:      dbPool,
		queries: queries,
	}, nil
}

func (c *PgCache) Get(ctx context.Context, phoneNumber string) (int, error) {
	var (
		version int
	)

	if err := c.db.QueryRow(
		ctx,
		c.queries.GetFromCache,
		phoneNumber,
	).Scan(&version); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return version, err
		}
	}

	return version, nil
}

func (c *PgCache) Set(ctx context.Context, phoneNumber string, version int) error {
	if _, err := c.db.Exec(
		ctx,
		c.queries.SetToCache,
		phoneNumber,
		version,
	); err != nil {
		return err
	}

	return nil
}

func (c *PgCache) Update(ctx context.Context, phoneNumber string, version int) error {
	if _, err := c.db.Exec(
		ctx,
		c.queries.UpdateCache,
		phoneNumber,
		version,
	); err != nil {
		return err
	}

	return nil
}
