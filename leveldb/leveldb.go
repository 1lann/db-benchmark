// Package leveldb contains a leveldb benchmark wrapper implementation
// through the use of the pure Go leveldb implementation at github.com/syndtr/goleveldb/leveldb
package leveldb

import (
	"github.com/1lann/db-benchmark/benchmark"
	"github.com/syndtr/goleveldb/leveldb"
)

type Wrapper struct {
	db *leveldb.DB
}

func init() {
	benchmark.RegisterWrapper("leveldb", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	var err error
	w.db, err = leveldb.OpenFile(opts.DB, nil)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Get(id string) benchmark.Document {
	val, err := w.db.Get([]byte(id), nil)
	if err != nil {
		panic(err)
	}
	return benchmark.Document{id, string(val)}
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var all []benchmark.Document

	iter := w.db.NewIterator(nil, nil)
	for iter.Next() {
		all = append(all, benchmark.Document{
			ID:    string(iter.Key()),
			Value: string(iter.Value()),
		})
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		panic(err)
	}

	return all
}

func (w *Wrapper) Update(doc benchmark.Document) {
	err := w.db.Put([]byte(doc.ID), []byte(doc.Value), nil)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Put(doc benchmark.Document) {
	err := w.db.Put([]byte(doc.ID), []byte(doc.Value), nil)
	if err != nil {
		panic(err)
	}
}

func (w *Wrapper) Clear() {
	all := w.GetAll()

	trans, err := w.db.OpenTransaction()
	if err != nil {
		panic(err)
	}

	for _, doc := range all {
		err = trans.Delete([]byte(doc.ID), nil)
		if err != nil {
			panic(err)
		}
	}

	err = trans.Commit()
	if err != nil {
		panic(err)
	}
}
