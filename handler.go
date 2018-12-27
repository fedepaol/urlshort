package urlshort

import (
	"net/http"

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

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedYaml, fallback), nil
}

func parseYAML(yml []byte) (map[string]string, error) {
	m := make([]yamlRecord, 0)
	err := yaml.Unmarshal([]byte(yml), &m)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for _, r := range m {
		res[r.Path] = r.URL
	}
	return res, nil
}

type yamlRecord struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
}
