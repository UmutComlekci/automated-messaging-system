package database

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/umutcomlekci/automated-messaging-system/internal/config"
	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
)

type postgreSqlDatabase struct {
	db     *sql.DB
	logger *slog.Logger
}

func newPostgreSqlDatabase() (*postgreSqlDatabase, error) {
	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		return nil, err
	}

	return &postgreSqlDatabase{
		db:     db,
		logger: logging.NewLogger("postgresql"),
	}, nil
}

func (d *postgreSqlDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *postgreSqlDatabase) QueryRow(query string, args ...any) *sql.Row {
	return d.db.QueryRow(query, args...)
}

func (d *postgreSqlDatabase) Exec(query string, args ...any) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *postgreSqlDatabase) Close() error {
	return d.db.Close()
}
