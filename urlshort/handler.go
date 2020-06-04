package urlshort

import (
	"fmt"
	"net/http"
	yaml "gopkg.in/yaml.v2"
	"encoding/json"
)

func describe(i interface{}) {
	fmt.Printf("(%+v, %T)\n", i, i)
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
//		describe(pathsToUrls)
		originalURL, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, originalURL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
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
/*
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYamlOrJson(yamlBytes, "yaml")
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYamlOrJson(jsonBytes, "json")
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}
*/
func DataHandler(dataType string, dataBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYamlOrJson(dataType, dataBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

/*
func YAMLbuildMap(pathUrls []YAMLpathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func JSONbuildMap(pathUrls []JSONpathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}
*/

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYamlOrJson(dataType string, data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	var err error
	switch dataType {
	case "json":
		err = json.Unmarshal(data, &pathUrls)
	case "yaml":
		err = yaml.Unmarshal(data, &pathUrls)
	}
	if err != nil {
		return nil, err
	}
	return pathUrls, err

}

type pathUrl struct {
	Path string `json:"path yaml:"path"`
	URL  string `json:"url  yaml:"url"`
}
