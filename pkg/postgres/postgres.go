package postgres

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunslog"
	"log/slog"
	"time"
)

func New(pgUrl string, logger *slog.Logger) *bun.DB {
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

	hook := bunslog.NewQueryHook(
		bunslog.WithLogger(logger),
		bunslog.WithQueryLogLevel(slog.LevelDebug),
		bunslog.WithSlowQueryLogLevel(slog.LevelWarn),
		bunslog.WithErrorQueryLogLevel(slog.LevelError),
		bunslog.WithSlowQueryThreshold(3*time.Second),
	)
	db.AddQueryHook(hook)

	return db
}
