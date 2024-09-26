package main

import (
	"encoding/json"
	"flag"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/brianvoe/gofakeit/v7"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

// Logger is the global logger for the application.
func parseFlags() (string, string, string, bool) {
	var jsonnetRootFolder string
	var fileSearchPattern string
	var outputFolder string
	var generateTestData bool

	flag.StringVar(&jsonnetRootFolder, "jsonnetRootFolder", "", "Root folder where all the jsonnet files")
	flag.StringVar(&fileSearchPattern, "fileSearchPattern", "", "Jsonnet file(s) name patterns")
	flag.StringVar(&outputFolder, "outputFolder", jsonnetRootFolder, "Root folder to generate json files of jsonnet")
	flag.BoolVar(&generateTestData, "generateTestData", false, "Jsonnet extVar flag for controlling test data generation")

	flag.Parse()

	if jsonnetRootFolder == "" || fileSearchPattern == "" {
		flag.PrintDefaults()
		Logger.Panic("Must provide jsonnet root folder and file name pattern!!")
	}

	return jsonnetRootFolder, fileSearchPattern, outputFolder, generateTestData
}

// processJsonnetFiles processes all the jsonnet files in the given root folder and generates JSON files in the output folder.
func processJsonnetFiles(jsonnetRootFolder, fileSearchPattern, outputFolder string, generateTestData bool) {
	fileSearchPath := path.Join(jsonnetRootFolder, fileSearchPattern)
	jsonnetFiles, globErr := doublestar.Glob(fileSearchPath)
	if globErr != nil {
		Logger.Panic(globErr)
		return
	}

	Logger.Infof("Found %d jsonnet files at the jsonnet path", len(jsonnetFiles))

	for _, jsonnetFile := range jsonnetFiles {
		processJsonnetFile(jsonnetFile, jsonnetRootFolder, outputFolder, generateTestData)
	}
}

// processJsonnetFile processes a single jsonnet file and generates a JSON file in the output folder.
func processJsonnetFile(jsonnetFile, jsonnetRootFolder, outputFolder string, generateTestData bool) {
	Logger.Infof("Input jsonnet file: %s", jsonnetFile)
	tempDir := strings.ReplaceAll(filepath.Dir(jsonnetFile), jsonnetRootFolder, outputFolder)
	Logger.Infof("tempDir path: %s", tempDir)
	tempFile := strings.ReplaceAll(filepath.Base(jsonnetFile), ".jsonnet", ".json")
	createDirAndWriteFile(tempDir, tempFile, jsonnetFile, generateTestData)
}

// createDirAndWriteFile creates a directory and writes the generated JSON file.
func createDirAndWriteFile(tempDir, tempFile, jsonnetFile string, generateTestData bool) {
	errDir := os.MkdirAll(tempDir, os.ModePerm)
	if errDir != nil {
		Logger.Error(errDir)
		return
	}

	outputFile := filepath.Join(tempDir, tempFile)
	Logger.Infof("Output json file: %s", outputFile)
	jsonData := GenerateJsonFromTemplate(jsonnetFile, generateTestData)
	err := os.WriteFile(outputFile, []byte(jsonData), 0777)
	if err != nil {
		Logger.Panic(err.Error())
	}
}

// createFakeFunction creates a jsonnet native function for generating fake data.
func createFakeFunction() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "fake",
		Params: ast.Identifiers{"x"},
		Func: func(x []interface{}) (interface{}, error) {
			bytes, err := json.Marshal(x[0])
			if err != nil {
				return nil, err
			}
			return CallGoFakeIt(strings.Trim(string(bytes), "\""))
		},
	}
}

// CallGoFakeIt generates fake data using gofakeit library.
func CallGoFakeIt(pattern string) (string, error) {
	return gofakeit.Generate(pattern)
}

// GenerateJsonFromTemplate processes a jsonnet template file and returns the generated JSON string.
func GenerateJsonFromTemplate(templateFilePath string, generateTestData bool) string {
	vm := jsonnet.MakeVM()

	// Set up the jsonnet VM
	vm.Importer(&jsonnet.FileImporter{})
	vm.ExtVar("generateTestData", strconv.FormatBool(generateTestData))
	vm.NativeFunction(createFakeFunction())

	// Evaluate the jsonnet file
	jsonStr, err := vm.EvaluateFile(templateFilePath)
	if err != nil {
		Logger.Panic(err.Error())
	}
	return jsonStr
}
