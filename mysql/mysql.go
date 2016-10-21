package mysql

import (
	"database/sql"

	"github.com/1lann/db-benchmark/benchmark"
	"github.com/1lann/db-benchmark/sqlwrapper"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	sqlwrapper.RegisterSQL("mysql", Connect)
}

func Connect(opts benchmark.ConnectOpts) *sql.DB {
	db, err := sql.Open("mysql", opts.Username+":"+opts.Password+"@/"+opts.DB)
	if err != nil {
		panic(err)
	}
	return db
}
