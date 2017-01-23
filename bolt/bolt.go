// Package bolt contains a bolt benchmark wrapper implementation
// through the use of the bolt embedded database at github.com/boltdb/bolt
package leveldb

import (
	"github.com/1lann/db-benchmark/benchmark"
	"github.com/boltdb/bolt"
)

type Wrapper struct {
	db    *bolt.DB
	table string
}

func init() {
	benchmark.RegisterWrapper("bolt", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.db, err = bolt.Open(opts.DB, 0600, nil)
	if err != nil {
		panic(err)
	}

	tx, err := w.db.Begin(true)
	if err != nil {
		panic(err)
	}

	_, err = tx.CreateBucketIfNotExists([]byte(opts.Table))
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	w.table = opts.Table
}

func (w *Wrapper) Get(id string) benchmark.Document {
	tx, err := w.db.Begin(false)
	if err != nil {
		panic(err)
	}

	val := tx.Bucket([]byte(w.table)).Get([]byte(id))

	err = tx.Rollback()
	if err != nil {
		panic(err)
	}

	return benchmark.Document{ID: id, Value: string(val)}
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var all []benchmark.Document

	tx, err := w.db.Begin(false)
	if err != nil {
		panic(err)
	}

	bucket := tx.Bucket([]byte(w.table))

	bucket.ForEach(func(k []byte, v []byte) error {
		all = append(all, benchmark.Document{
			ID:    string(k),
			Value: string(v),
		})
		return nil
	})

	err = tx.Rollback()
	if err != nil {
		panic(err)
	}

	return all
}

func (w *Wrapper) Update(doc benchmark.Document) {
	tx, err := w.db.Begin(true)
	if err != nil {
		panic(err)
	}

	err = tx.Bucket([]byte(w.table)).Put([]byte(doc.ID), []byte(doc.Value))
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	tx, err := w.db.Begin(true)
	if err != nil {
		panic(err)
	}

	err = tx.Bucket([]byte(w.table)).Put([]byte(doc.ID), []byte(doc.Value))
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	tx, err := w.db.Begin(true)
	if err != nil {
		panic(err)
	}

	err = tx.DeleteBucket([]byte(w.table))
	if err != nil {
		panic(err)
	}

	_, err = tx.CreateBucket([]byte(w.table))
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
