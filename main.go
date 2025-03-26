package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/sgreben/flagvar"
	"github.com/xeipuuv/gojsonschema"
)

var validationErrors map[string][]string

var (
	schemaFile   string
	outputFormat flagvar.Enum
)

func main() {
	validationErrors = make(map[string][]string)
	outputFormat.Choices = []string{"text", "json"}
	outputFormat.Value = "text"

	flag.StringVar(&schemaFile, "schema", "./schema.json", "load `FILE` as the schema when validating")
	flag.Var(&outputFormat, "format", "output in `FORMAT` (json, text)")

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
			validationErrors[filepath.Base(file)] = make([]string, 0)
			validationErrors[filepath.Base(file)] = append(validationErrors[filepath.Base(file)], errors...)
		}
	}

	switch outputFormat.Value {
	case "json":
		data, err := json.Marshal(validationErrors)
		if err != nil {
			log.Fatalf("cannot marshal json: %v", err)
		}
		fmt.Println(string(data))
	case "text":
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
		fmt.Fprint(writer, "FILE\tERROR\n")
		for file, errors := range validationErrors {
			for _, err := range errors {
				fmt.Fprintf(writer, "%v\t%v\n", file, err)
			}
		}
		if err := writer.Flush(); err != nil {
			log.Fatalf("unable to flush tab writer: %v", err)
		}
	default:
		log.Fatalf("unknown format type: %v", outputFormat.Value)
	}
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
