// Package sqlwrapper contains a benchmark wrapper for SQL based databases.
package sqlwrapper

import (
	"database/sql"

	"github.com/1lann/db-benchmark/benchmark"
)

type ConnectFunc func(opts benchmark.ConnectOpts) *sql.DB

type SQLDatabase struct {
	opts       benchmark.ConnectOpts
	connect    ConnectFunc
	db         *sql.DB
	getStmt    *sql.Stmt
	getAllStmt *sql.Stmt
	updateStmt *sql.Stmt
	putStmt    *sql.Stmt
	clearStmt  *sql.Stmt
}

func RegisterSQL(name string, connect ConnectFunc) {
	benchmark.RegisterWrapper(name, &SQLDatabase{connect: connect})
}

func (w *SQLDatabase) Connect(opts benchmark.ConnectOpts) {
	w.opts = opts
	w.db = w.connect(opts)
	w.db.SetMaxIdleConns(80)
	w.db.SetMaxOpenConns(80)

	_, err := w.db.Exec("DROP TABLE IF EXISTS " + w.opts.Table)
	if err != nil {
		panic(err)
	}

	_, err = w.db.Exec("CREATE TABLE " + w.opts.Table + ` (
	id VARCHAR(20) NOT NULL,
	value VARCHAR(1000),
	PRIMARY KEY (id)
)`)
	if err != nil {
		panic(err)
	}

	w.getStmt, err = w.db.Prepare("SELECT id, value FROM " + w.opts.Table + " WHERE id = ?")
	if err != nil {
		panic(err)
	}

	w.getAllStmt, err = w.db.Prepare("SELECT id, value FROM " + w.opts.Table)
	if err != nil {
		panic(err)
	}

	w.updateStmt, err = w.db.Prepare("UPDATE " + w.opts.Table + " SET id = ?, value = ? " +
		"WHERE id = ?")
	if err != nil {
		panic(err)
	}

	w.putStmt, err = w.db.Prepare("INSERT INTO " + w.opts.Table + " VALUES (?, ?)")
	if err != nil {
		panic(err)
	}

	w.clearStmt, err = w.db.Prepare("DELETE FROM " + w.opts.Table)
	if err != nil {
		panic(err)
	}
}

func (w *SQLDatabase) Get(id string) benchmark.Document {
	result := w.getStmt.QueryRow(id)
	var doc benchmark.Document
	err := result.Scan(&doc.ID, &doc.Value)
	if err != nil {
		panic(err)
	}

	return doc
}

func (w *SQLDatabase) GetAll() []benchmark.Document {
	results, err := w.getAllStmt.Query()
	if err != nil {
		panic(err)
	}

	var docs []benchmark.Document
	for results.Next() {
		var doc benchmark.Document
		err := results.Scan(&doc.ID, &doc.Value)
		if err != nil {
			panic(err)
		}

		docs = append(docs, doc)
	}

	return docs
}

func (w *SQLDatabase) Update(doc benchmark.Document) {
	_, err := w.updateStmt.Exec(doc.ID, doc.Value, doc.ID)
	if err != nil {
		panic(err)
	}
}

func (w *SQLDatabase) Put(doc benchmark.Document) {
	_, err := w.putStmt.Exec(doc.ID, doc.Value)
	if err != nil {
		panic(err)
	}
}

func (w *SQLDatabase) Clear() {
	_, err := w.clearStmt.Exec()
	if err != nil {
		panic(err)
	}
}
