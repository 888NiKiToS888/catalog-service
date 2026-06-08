package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/888NiKiToS888/catalog-service/internal/app/config/section"
	"github.com/888NiKiToS888/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB
		cfg      section.RepositoryPostgres
	}
	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewConn(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	u := &url.URL{
		Scheme: "postgres",
		Host:   cfg.Address,
		Path:   cfg.Name,
	}
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	q := url.Values{}
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	dsn := u.String()

	log.Debug().Str("dsn", dsn).Msg("PostgreSQL DNS")

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(cfg.ReadTimeout.Duration),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout.Duration),
	))
	sqlDB.SetMaxOpenConns(10)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	client := &Client{
		rawBunDB: bunDB,
		cfg:      cfg,
	}
	client._bunDB = bunDB
	return client, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}
	migrator := migrate.NewMigrator(c.rawBunDB, migrations,
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable+"_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	)
	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migrator: %w", err)
	}
	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	if len(applied) > 0 {
		oldVer = applied[0].ID
	} else {
		oldVer = 0
	}
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return oldVer, 0, fmt.Errorf("failed to migrate: %w", err)
	}
	newVer = oldVer
	for _, m := range group.Migrations {
		if m.ID > newVer {
			newVer = m.ID
		}
	}
	return oldVer, newVer, nil
}
