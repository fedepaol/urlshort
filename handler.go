package urlshort

import (
	"net/http"

	"github.com/dgraph-io/badger"
	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirect, ok := pathsToUrls[path]
		if !ok {
			fallback.ServeHTTP(w, r)
		}
		http.Redirect(w, r, redirect, 301)
	}
}

func LoadBadgerFromYaml(db *badger.DB, yml []byte) error {
	m := make([]yamlRecord, 0)
	err := yaml.Unmarshal([]byte(yml), &m)
	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		for _, r := range m {
			err := txn.Set([]byte(r.Path), []byte(r.URL))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func BadgerHandler(db *badger.DB, fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var redirect string
		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(path))
			if err != nil {
				return err
			}

			v, err := item.Value()
			redirect = string(v)
			return nil
		})

		if err != nil {
			fallback.ServeHTTP(w, r)
		}
		http.Redirect(w, r, redirect, 301)
	}, nil

}

type yamlRecord struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
}
