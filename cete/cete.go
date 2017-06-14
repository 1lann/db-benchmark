// Package cete contains a cete benchmark wrapper implementation.
package cete

import (
	"github.com/1lann/cete"

	"github.com/1lann/db-benchmark/benchmark"
)

type Wrapper struct {
	db *cete.DB
}

func init() {
	benchmark.RegisterWrapper("cete", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.db, err = cete.OpenDatabase(opts.DB)
	if err != nil {
		panic(err)
	}

	w.db.NewTable("data")
	if err != cete.ErrAlreadyExists && err != nil {
		panic(err)
	}
	err = w.db.Table("data").NewIndex("ID")
	if err != cete.ErrAlreadyExists && err != nil {
		panic(err)
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	var doc benchmark.Document
	_, _, err := w.db.Table("data").Index("ID").One(id, &doc)
	if err != nil {
		panic(err)
	}

	return doc
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var documents []benchmark.Document

	r := w.db.Table("data").All()
	defer r.Close()

	var err error
	for {
		var doc benchmark.Document
		_, _, err = r.Next(&doc)
		if err != nil {
			return documents
		}

		documents = append(documents, doc)
	}
}

func (w *Wrapper) Update(doc benchmark.Document) {
	err := w.db.Table("data").Set(doc.ID, doc)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	err := w.db.Table("data").Set(doc.ID, doc)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	err := w.db.Table("data").Drop()
	if err != nil {
		panic(err)
	}

	err = w.db.NewTable("data")
	if err != nil {
		panic(err)
	}
	err = w.db.Table("data").NewIndex("ID")
	if err != nil {
		panic(err)
	}
}
