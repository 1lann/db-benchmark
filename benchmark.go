package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/1lann/db-benchmark/benchmark"
	_ "github.com/1lann/db-benchmark/memory"
	_ "github.com/1lann/db-benchmark/mysql"
	_ "github.com/1lann/db-benchmark/rethinkdb"
)

var documents []benchmark.Document
var documentMap = make(map[string]int)

var connectOpts = map[string]benchmark.ConnectOpts{
	"rethinkdb": {
		DB:       "test",
		Table:    "benchmark",
		Username: "admin",
		Host:     "127.0.0.1:28015",
	},
	"mysql": {
		DB:       "test",
		Table:    "benchmark",
		Username: "root",
	},
	"memory": {},
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10000; i++ {
		randomData := make([]byte, 400)
		_, err := rand.Read(randomData)
		if err != nil {
			panic(err)
		}

		value := base64.StdEncoding.EncodeToString(randomData)

		documents = append(documents, benchmark.Document{
			ID:    strconv.Itoa(i) + "_object",
			Value: value,
		})

		documentMap[strconv.Itoa(i)+"_object"] = i
	}

	for db, opts := range connectOpts {
		w := benchmark.Wrappers[db]
		w.Connect(opts)
		fmt.Println("\n--- Starting benchmark for", db, "---")
		for i := 1; i <= 3; i++ {
			fmt.Println("\nWrite benchmark pass", i)
			runWrites(w)
		}
		fmt.Println("\nEnd of writes benchmark")
		for i := 1; i <= 3; i++ {
			fmt.Println("\nRead benchmark pass", i)
			runReads(w)
		}
	}

	fmt.Println("\n--- All benchmarks complete ---")
}

func runWrites(w benchmark.Wrapper) {
	w.Clear()
	fmt.Println("Putting", len(documents), "documents sequentially...")
	start := time.Now()
	for _, doc := range documents {
		w.Put(doc)
	}
	fmt.Println("Took", time.Now().Sub(start).Seconds(), "seconds")
	w.Clear()

	fmt.Println("Putting", len(documents), "documents simutaneously...")
	wg := new(sync.WaitGroup)
	wg.Add(len(documents))

	start = time.Now()
	for _, doc := range documents {
		go func(idoc benchmark.Document) {
			w.Put(idoc)
			wg.Done()
		}(doc)
	}

	wg.Wait()
	fmt.Println("Took", time.Now().Sub(start).Seconds(), "seconds")
}

func runReads(w benchmark.Wrapper) {
	fmt.Println("Getting", len(documents), "documents sequentially...")
	start := time.Now()
	for _, doc := range documents {
		dbDoc := w.Get(doc.ID)
		if dbDoc.ID != doc.ID || dbDoc.Value != doc.Value {
			panic("benchmark error: document integrity error")
		}
	}
	fmt.Println("Took", time.Now().Sub(start).Seconds(), "seconds")

	fmt.Println("Getting", len(documents), "documents simutaneously...")
	wg := new(sync.WaitGroup)
	wg.Add(len(documents))

	start = time.Now()
	for _, doc := range documents {
		go func(idoc benchmark.Document) {
			dbDoc := w.Get(idoc.ID)
			if dbDoc.ID != idoc.ID || dbDoc.Value != idoc.Value {
				panic("benchmark error: document integrity error")
			}
			wg.Done()
		}(doc)
	}

	wg.Wait()
	fmt.Println("Took", time.Now().Sub(start).Seconds(), "seconds")

	fmt.Println("Getting all documents...")
	start = time.Now()
	all := w.GetAll()
	fmt.Println("Took", time.Now().Sub(start).Seconds(), "seconds")

	for _, allDoc := range all {
		doc := documents[documentMap[allDoc.ID]]
		if allDoc.ID != doc.ID || allDoc.Value != doc.Value {
			panic("benchmark error: document integrity error")
		}
	}
}