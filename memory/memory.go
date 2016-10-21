// Package memory contains a in-memory benchmark wrapper implementation
// through the use of Go's maps.
package memory

import (
	"sync"

	"github.com/1lann/db-benchmark/benchmark"
)

type Wrapper struct {
	db   map[string]benchmark.Document
	lock *sync.RWMutex
}

func init() {
	benchmark.RegisterWrapper("memory", new(Wrapper))
}

func (w *Wrapper) Connect(opts benchmark.ConnectOpts) {
	w.db = make(map[string]benchmark.Document)
	w.lock = new(sync.RWMutex)
}

func (w *Wrapper) Get(id string) benchmark.Document {
	w.lock.RLock()
	val := w.db[id]
	w.lock.RUnlock()
	return val
}

func (w *Wrapper) GetAll() []benchmark.Document {
	var all []benchmark.Document
	w.lock.RLock()
	for _, value := range w.db {
		all = append(all, value)
	}
	w.lock.RUnlock()
	return all
}

func (w *Wrapper) Update(doc benchmark.Document) {
	w.lock.Lock()
	w.db[doc.ID] = doc
	w.lock.Unlock()
}

func (w *Wrapper) Put(doc benchmark.Document) {
	w.lock.Lock()
	w.db[doc.ID] = doc
	w.lock.Unlock()
}

func (w *Wrapper) Clear() {
	w.lock.Lock()
	w.db = make(map[string]benchmark.Document)
	w.lock.Unlock()
}
