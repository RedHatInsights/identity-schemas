package main

import (
    "fmt"
    "io/fs"
    "os"
    "strings"
    "github.com/xeipuuv/gojsonschema"
)

var validationErrors []string

func main() {
    dirs, _ := os.ReadDir("./")

    for _, dir := range dirs {
        validate(dir)
    }

    outputResults()
}

func validate(dir fs.DirEntry) {
    dirName := dir.Name()
    if dir.IsDir() && stringContains(gatewayDirNames(), dirName) {
        fmt.Println("Running schema validation against sample " + dirName + " identities")
        schemaLoader := gojsonschema.NewReferenceLoader("file://" + dirName + "/schema.json")

        files, _ := os.ReadDir(dirName + "/identities")
        for _, file := range files {
            fileName := file.Name()
            if !file.IsDir() {
                documentLoader := gojsonschema.NewReferenceLoader("file://" + dirName + "/identities/" + fileName)

                result, err := gojsonschema.Validate(schemaLoader, documentLoader)
                if err != nil {
                    panic(err.Error())
                }

                if !result.Valid() {
                    for _, desc := range result.Errors() {
                        e := "- " + dirName + "/" + fileName + ": " + desc.Description()
                        validationErrors = append(validationErrors, e)
                    }
                }
            }
        }
    }
}

func stringContains(slice []string, ele string) bool {
    for _, sliceEle := range slice {
        if sliceEle == ele {
            return true
        }
    }

    return false
}

func gatewayDirNames() []string {
    return []string{"3scale", "turnpike"}
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
