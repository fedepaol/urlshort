package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgraph-io/badger"
	"github.com/gophercises/urlshort"
)

func main() {
	mux := defaultMux()

	filename := flag.String("mapfile", "", "the file containing the map between path and url")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		log.Fatalf("Must provide a valid file name")
	}

	res, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatalf("Failed to open %s", *filename)
	}

	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = urlshort.LoadBadgerFromYaml(db, res)
	if err != nil {
		panic(err)
	}
	badgerHandler, err := urlshort.BadgerHandler(db, mux)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", badgerHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
