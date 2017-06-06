package postgres

import (
	"database/sql"

	"github.com/1lann/db-benchmark/benchmark"
	"github.com/1lann/db-benchmark/psqlwrapper"
	_ "github.com/lib/pq"
)

func init() {
	psqlwrapper.RegisterSQL("postgres", Connect)
}

func Connect(opts benchmark.ConnectOpts) *sql.DB {
	db, err := sql.Open("postgres", "postgres://"+opts.Username+":"+opts.Password+"@"+opts.Host+"/"+opts.DB+"?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}
