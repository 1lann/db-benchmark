// Package buntdb contains a buntdb benchmark wrapper implementation
// through the use of the pure Go buntdb implementation at github.com/tidwall/buntdb
package buntdb

import (
	"github.com/1lann/db-benchmark/benchmark"
	"github.com/tidwall/buntdb"
)

type Wrapper struct {
	db *buntdb.DB
}

func init() {
	benchmark.RegisterWrapper("buntdb", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.db, err = buntdb.Open(opts.DB)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	var text string

	err := w.db.View(func(tx *buntdb.Tx) error {
		var err error
		text, err = tx.Get(id)
		return err
	})
	if err != nil {
		panic(err)
	}

	return benchmark.Document{
		ID:    id,
		Value: text,
	}
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var all []benchmark.Document

	err := w.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key string, value string) bool {
			all = append(all, benchmark.Document{ID: key, Value: value})
			return true
		})

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return all
}

func (w *Wrapper) Update(doc benchmark.Document) {
	err := w.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(doc.ID, doc.Value, nil)
		return err
	})

	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	err := w.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(doc.ID, doc.Value, nil)
		return err
	})

	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	err := w.db.Update(func(tx *buntdb.Tx) error {
		return tx.DeleteAll()
	})

	if err != nil {
		panic(err)
	}
}
