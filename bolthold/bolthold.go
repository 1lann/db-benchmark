// Package bolthold contains a bolthold benchmark wrapper implementation.
package bolthold

import (
	"github.com/timshannon/bolthold"

	"github.com/1lann/db-benchmark/benchmark"
)

type Wrapper struct {
	db *bolthold.Store
}

func init() {
	benchmark.RegisterWrapper("bolthold", new(Wrapper))
}

type boltholdDocument struct {
	ID    string `boltholdIndex:"ID"`
	Value string
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.db, err = bolthold.Open(opts.DB, 0666, nil)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	var doc []boltholdDocument

	err := w.db.Find(&doc, bolthold.Where("ID").Eq(id))
	if err != nil {
		panic(err)
	}

	return benchmark.Document{
		ID:    doc[0].ID,
		Value: doc[0].Value,
	}
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var documents []boltholdDocument

	err := w.db.Find(&documents, nil)
	if err != nil {
		panic(err)
	}

	returnDocs := make([]benchmark.Document, len(documents))

	for k, v := range documents {
		returnDocs[k] = benchmark.Document{
			ID: v.ID,
			Value: v.Value,
		}
	}

	return returnDocs
}

func (w *Wrapper) Update(doc benchmark.Document) {
	err := w.db.Update(doc.ID, boltholdDocument{
		ID: doc.ID,
		Value: doc.Value,
	})
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	err := w.db.Insert(doc.ID, boltholdDocument{
		ID: doc.ID,
		Value: doc.Value,
	})
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	w.db.DeleteMatching(&boltholdDocument{}, nil)
}
