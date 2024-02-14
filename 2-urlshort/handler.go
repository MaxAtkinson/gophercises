package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentPath := r.URL.Path

		redirectTo, ok := pathsToUrls[currentPath]
		if ok {
			http.Redirect(w, r, redirectTo, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// Note to self: For unmarshalling, fields must be public
// (path => Path, url => Url)
type resolver struct {
	Path string
	Url  string
}

func parseYaml(yml []byte) (map[string]string, error) {
	var parseResult []resolver
	err := yaml.Unmarshal(yml, &parseResult)

	result := make(map[string]string)
	for _, value := range parseResult {
		result[value.Path] = value.Url
	}

	return result, err
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYaml(yml)
	return MapHandler(pathsToUrls, fallback), err
}
