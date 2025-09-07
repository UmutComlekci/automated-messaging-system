package database

import (
	"database/sql"
	"fmt"

	"github.com/umutcomlekci/automated-messaging-system/internal/config"
)

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Close() error
}

func NewDatabase() (Database, error) {
	if config.GetSqlDriver() == "postgres" {
		return newPostgreSqlDatabase()
	}

	return nil, fmt.Errorf("unsupported database driver: %s", config.GetSqlDriver())
}
