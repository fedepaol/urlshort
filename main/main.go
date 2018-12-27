package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	yamlHandler, err := urlshort.YAMLHandler(res, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
