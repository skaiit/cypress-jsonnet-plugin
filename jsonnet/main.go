package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar"
)

/*
Usage

	building binary from npm index.js

	    go build -C ./node_modules/cy-jsonnet/jsonnet

	execution

	    cy-jsonnet.[sh|exe]  {jsonnet root folder} {jsonnet file pattern} {optional output folder} {optional generateTestData extvar flag}
	    eg :- cy-jsonnet.exe --jsonnetRootFolder=support/jsonnet --fileSearchPattern=*.jsonnet --outputFolder=fixtures/testData --generateTestData

	output
	    generated json file based on jsonnet definition
	    if no outputfile path mentioned, out file will be in the same path as jsonnet input file with .json extension sufix
	    if no generateTestData provided, it will be set to false
*/
func main() {
	start := time.Now()
	// Define flags
	var jsonnetRootFolder string
	var fileSearchPattern string
	var outputFolder string
	var generateTestData bool

	flag.StringVar(&jsonnetRootFolder, "jsonnetRootFolder", "", "Root folder where all the jsonnet files")
	flag.StringVar(&fileSearchPattern, "fileSearchPattern", "", "Jsonnet file(s) name patterns")
	flag.StringVar(&outputFolder, "outputFolder", jsonnetRootFolder, "Root folder to generate json files of jsonnet")
	flag.BoolVar(&generateTestData, "generateTestData", false, "Jsonnet extVar flag for controlling test data generation")

	// Parse command line arguments
	flag.Parse()

	if jsonnetRootFolder == "" || fileSearchPattern == "" {
		log.Panic("Must provide jsonnet root folder and file name pattern!!")
		flag.PrintDefaults()
		return
	}

	fmt.Printf("jsonnet Root Folder : %s\n", jsonnetRootFolder)
	fmt.Printf("fileSearchPattern : %s\n", fileSearchPattern)
	fmt.Printf("Output Folder: %s\n", outputFolder)
	fmt.Printf("generateTestData: %t\n", generateTestData)

	pathToSearch := path.Join(jsonnetRootFolder, fileSearchPattern)
	jsonnetFiles, globErr := doublestar.Glob(pathToSearch) // Using doublestar dependency because Go's native glob method does not support ** in grep patterns
	if globErr != nil {
		log.Panic(globErr)
		return
	}

	fmt.Println("Found", len(jsonnetFiles), "jsonnet files at the jsonnet path")

	for _, jsonnetFile := range jsonnetFiles {
		fmt.Printf("Input jsonnet file: %s\n", jsonnetFile)
		tempDir := strings.ReplaceAll(filepath.Dir(jsonnetFile), jsonnetRootFolder, outputFolder)
		fmt.Printf("tempDir path : %s\n", tempDir)
		tempFile := strings.ReplaceAll(filepath.Base(jsonnetFile), ".jsonnet", ".json")
		errDir := os.MkdirAll(tempDir, os.ModePerm)
		if errDir != nil {
			fmt.Println(errDir)
			return
		}
		outputFile := filepath.Join(tempDir, tempFile)
		fmt.Printf("Output json file: %s\n", outputFile)
		jsonData := Load(jsonnetFile, generateTestData)
		err := os.WriteFile(outputFile, []byte(jsonData), 0777)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	duration := time.Since(start)
	fmt.Println("Completed processing", os.Args, ". Time in nano seconds", duration.Nanoseconds())
}
