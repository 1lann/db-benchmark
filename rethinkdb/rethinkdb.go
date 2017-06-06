// Package rethinkdb contains a benchmark wrapper for RethinkDB.
package rethinkdb

import (
	"log"

	"github.com/1lann/db-benchmark/benchmark"
	r "github.com/dancannon/gorethink"
)

type Wrapper struct {
	s     *r.Session
	table string
}

func init() {
	benchmark.RegisterWrapper("rethinkdb", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.s, err = r.Connect(r.ConnectOpts{
		Address:    opts.Host,
		Database:   opts.DB,
		Username:   opts.Username,
		Password:   opts.Password,
		InitialCap: 1000,
		MaxOpen:    1000,
	})
	if err != nil {
		panic(err)
	}

	w.table = opts.Table

	r.TableDrop(w.table).RunWrite(w.s)
	resp, err := r.TableCreate(w.table, r.TableCreateOpts{PrimaryKey: "ID"}).RunWrite(w.s)
	if err != nil {
		panic(err)
	}

	if resp.TablesCreated == 0 {
		panic("rethinkdb: could not create table")
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	var dst benchmark.Document
	err := r.Table(w.table).Get(id).ReadOne(&dst, w.s)
	if err != nil {
		panic(err)
	}
	return dst
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var dst []benchmark.Document
	err := r.Table(w.table).ReadAll(&dst, w.s)
	if err != nil {
		panic(err)
	}
	return dst
}

func (w *Wrapper) Update(doc benchmark.Document) {
	resp, err := r.Table(w.table).Get(doc.ID).Update(doc, r.UpdateOpts{
		Durability: "soft",
	}).RunWrite(w.s)
	if err != nil {
		panic(err)
	}
	if resp.Errors > 0 {
		log.Fatal("rethinkdb update error: ", resp.FirstError)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	resp, err := r.Table(w.table).Insert(doc, r.InsertOpts{
		Durability: "soft",
	}).RunWrite(w.s)
	if err != nil {
		panic(err)
	}
	if resp.Errors > 0 {
		log.Fatal("rethinkdb insert error: ", resp.FirstError)
	}
}

func (w *Wrapper) Clear() {
	resp, err := r.Table(w.table).Delete(r.DeleteOpts{
		Durability: "soft",
	}).RunWrite(w.s)
	if err != nil {
		panic(err)
	}
	if resp.Errors > 0 {
		log.Fatal("rethinkdb delete error: ", resp.FirstError)
	}
}
