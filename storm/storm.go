// Package storm contains a storm benchmark wrapper implementation.
package storm

import (
	"github.com/asdine/storm"

	"github.com/1lann/db-benchmark/benchmark"
)

type Wrapper struct {
	db *storm.DB
}

func init() {
	benchmark.RegisterWrapper("storm", new(Wrapper))
}

type stormDocument struct {
	ID    string `storm:"id"`
	Index string `storm:"index"`
	Value string
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.db, err = storm.Open(opts.DB, storm.Batch())
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	var doc []stormDocument
	err := w.db.Find("Index", id, &doc)
	if err != nil {
		panic(err)
	}

	return benchmark.Document{
		ID:    doc[0].ID,
		Value: doc[0].Value,
	}
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var documents []stormDocument

	err := w.db.All(&documents)
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
	err := w.db.Update(&stormDocument{
		ID:    doc.ID,
		Index: doc.ID,
		Value: doc.Value,
	})
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	err := w.db.Save(&stormDocument{
		ID:    doc.ID,
		Index: doc.ID,
		Value: doc.Value,
	})
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	w.db.Drop(&stormDocument{})
}
