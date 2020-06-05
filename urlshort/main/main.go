package main

import (
	"fmt"
	"net/http"
//	"github.com/gophercises/urlshort"
	"github.com/gittes/gophercises/urlshort"
	"flag"
	"io/ioutil"
	"os"
)

func main() {
	yamlFilename := flag.String("yaml", "", "a YAML file with 'url and 'path' keys in an item list")
	jsonFilename := flag.String("json", "", "a JSON file with 'url and 'path' keys in an item list")
//	debugOutput := flag.Bool("debug", false, "Output some variables for debugging")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	fmt.Println("mapHandler initialized")

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml_default := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlBytes := []byte{0}
	jsonBytes := []byte{0}


	if isFlagPassed("json") { 
		jsonBytes = readAllFile(jsonFilename)
		fmt.Println("Loaded JSON file")
	}

	if isFlagPassed("yaml") {
		yamlBytes = readAllFile(yamlFilename)
		fmt.Println("Loaded YAML file")
	} else {
		yamlBytes = []byte(yaml_default)
		fmt.Println("Loaded yaml_default")
	}

	jsonHandler, err := urlshort.DataHandler("json", jsonBytes, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON HTTP Handler initialized")
		yamlHandler, err := urlshort.DataHandler("yaml", yamlBytes, jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("YAML HTTP Handler initialized")
	
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

func readAllFile(file *string) []byte {
	reader, _ := os.Open(*file)
	bytes, _ := ioutil.ReadAll(reader)
	defer reader.Close()
	return bytes
}

func isFlagPassed(name string) bool {
    found := false
    flag.Visit(func(f *flag.Flag) {
        if f.Name == name {
            found = true
        }
    })
    return found
}