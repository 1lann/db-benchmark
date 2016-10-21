package sqlite

import (
	"database/sql"

	"github.com/1lann/db-benchmark/benchmark"
	"github.com/1lann/db-benchmark/sqlwrapper"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	sqlwrapper.RegisterSQL("sqlite", Connect)
}

func Connect(opts benchmark.ConnectOpts) *sql.DB {
	db, err := sql.Open("sqlite3", opts.DB)
	if err != nil {
		panic(err)
	}
	return db
}
