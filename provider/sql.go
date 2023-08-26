package provider

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ProvideSQL() (*sql.DB, error) {
	var (
		host     = os.Getenv("PG_HOST")
		port     = os.Getenv("PG_PORT")
		user     = os.Getenv("PG_USER")
		password = os.Getenv("PG_PASSWORD")
		dbname   = os.Getenv("PG_DBNAME")
	)

	sql, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	))
	if err != nil {
		return nil, err
	}

	return sql, sql.Ping()
}
