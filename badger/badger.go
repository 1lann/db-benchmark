// Package badger contains a badger benchmark wrapper implementation through
// the pure Go badger key-value store at: github.com/dgraph-io/badger
package badger

import (
	"os"

	"github.com/dgraph-io/badger/badger"

	"github.com/1lann/db-benchmark/benchmark"
)

type Wrapper struct {
	db *badger.KV
}

func init() {
	benchmark.RegisterWrapper("badger", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	os.Mkdir(opts.DB, 0744)
	options := badger.DefaultOptions
	options.Dir = opts.DB
	options.ValueDir = opts.DB
	options.ValueCompressionMinSize = 10000000
	options.ValueCompressionMinRatio = 10
	var err error
	w.db, err = badger.NewKV(&options)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	var item badger.KVItem
	err := w.db.Get([]byte(id), &item)
	if err != nil {
		panic(err)
	}

	return benchmark.Document{
		ID:    id,
		Value: string(item.Value()),
	}
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var documents []benchmark.Document

	it := w.db.NewIterator(badger.DefaultIteratorOptions)
	for it.Rewind(); it.Valid(); it.Next() {
		documents = append(documents, benchmark.Document{
			ID:    string(it.Item().Key()),
			Value: string(it.Item().Value()),
		})
	}

	return documents
}

func (w *Wrapper) Update(doc benchmark.Document) {
	w.Put(doc)
}

func (w *Wrapper) Put(doc benchmark.Document) {
	err := w.db.Set([]byte(doc.ID), []byte(doc.Value))
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	it := w.db.NewIterator(badger.DefaultIteratorOptions)
	for it.Rewind(); it.Valid(); it.Next() {
		go w.db.Delete(it.Item().Key())
	}
}
