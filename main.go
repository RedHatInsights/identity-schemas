package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var validationErrors []string

var (
	schemaFile string
)

func main() {
	flag.StringVar(&schemaFile, "schema", "./schema.json", "load `FILE` as the schema when validating")

	flag.Parse()

	data, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("cannot load schema file: %v", err)
	}
	schema := gojsonschema.NewBytesLoader(data)

	for _, file := range flag.Args() {
		var data []byte
		var err error

		if file == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(file)
		}
		if err != nil {
			log.Fatalf("cannot load document: %v", err)
		}
		valid, errors, err := validateDocument(schema, string(data))
		if !valid {
			if err != nil {
				log.Fatalf("cannot validate document: %v", err)
			}
			for _, e := range errors {
				validationErrors = append(validationErrors, fmt.Sprintf("%v: %v", filepath.Base(file), e))
			}
		}
	}

	outputResults()
}

// validateDocument runs JSON schema validation using the given loader on the
// conents of the document.
func validateDocument(loader gojsonschema.JSONLoader, document string) (bool, []string, error) {
	result, err := gojsonschema.Validate(loader, gojsonschema.NewStringLoader(document))
	if err != nil {
		return false, nil, fmt.Errorf("cannot validate document: %w", err)
	}

	if !result.Valid() {
		errors := make([]string, 0)
		for _, desc := range result.Errors() {
			errors = append(errors, desc.Description())
		}
		return false, errors, nil
	}
	return true, nil, nil
}

func outputResults() {
	fmt.Println()

	if len(validationErrors) > 0 {
		fmt.Printf("Error(s) found in validation:\n")
		fmt.Println(strings.Join(validationErrors, "\n"))
		os.Exit(1)
	} else {
		fmt.Printf("No errors found.\n")
	}
}
