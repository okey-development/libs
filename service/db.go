package service

import (
	"database/sql"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

type dbConfig struct {
	Driver string `default:"pgx"`
	DSN    string `default:"postgres://postgres@127.0.0.1:5432/test"`
}

func initDB(config *dbConfig) error {
	var err error
	db, err = sql.Open(config.Driver, config.DSN)
	return err
}
func QueryDB(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func QueryRowDB(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}
