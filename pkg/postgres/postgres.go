package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunzerolog"
	"time"
)

func New(pgUrl string, logger *zerolog.Logger) *bun.DB {
	dsn := fmt.Sprintf("%s?sslmode=disable", pgUrl)
	connector := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	)
	pgDb := sql.OpenDB(connector)
	db := bun.NewDB(pgDb, pgdialect.New())

	hook := bunzerolog.NewQueryHook(
		bunzerolog.WithLogger(logger),
		bunzerolog.WithQueryLogLevel(zerolog.DebugLevel),
		bunzerolog.WithSlowQueryLogLevel(zerolog.WarnLevel),
		bunzerolog.WithErrorQueryLogLevel(zerolog.ErrorLevel),
		bunzerolog.WithSlowQueryThreshold(3*time.Second),
	)
	db.AddQueryHook(hook)

	return db
}
