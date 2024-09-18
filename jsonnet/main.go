package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

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
		panic("Must provide jsonnet root folder and file name pattern!!")
	}

	return jsonnetRootFolder, fileSearchPattern, outputFolder, generateTestData
}

func findFiles(root, pattern string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func processJsonnetFiles(jsonnetRootFolder, fileSearchPattern, outputFolder string, generateTestData bool, sugar *zap.SugaredLogger) {
	jsonnetFiles, globErr := findFiles(jsonnetRootFolder, fileSearchPattern)
	if globErr != nil {
		sugar.Panic(globErr)
		return
	}

	sugar.Infof("Found %d jsonnet files at the jsonnet path", len(jsonnetFiles))

	for _, jsonnetFile := range jsonnetFiles {
		processJsonnetFile(jsonnetFile, jsonnetRootFolder, outputFolder, generateTestData, sugar)
	}
}

func processJsonnetFile(jsonnetFile, jsonnetRootFolder, outputFolder string, generateTestData bool, sugar *zap.SugaredLogger) {
	sugar.Infof("Input jsonnet file: %s", jsonnetFile)
	tempDir := strings.ReplaceAll(filepath.Dir(jsonnetFile), jsonnetRootFolder, outputFolder)
	sugar.Infof("tempDir path: %s", tempDir)
	tempFile := strings.ReplaceAll(filepath.Base(jsonnetFile), ".jsonnet", ".json")
	createDirAndWriteFile(tempDir, tempFile, jsonnetFile, generateTestData, sugar)
}

func createDirAndWriteFile(tempDir, tempFile, jsonnetFile string, generateTestData bool, sugar *zap.SugaredLogger) {
	errDir := os.MkdirAll(tempDir, os.ModePerm)
	if errDir != nil {
		sugar.Error(errDir)
		return
	}

	outputFile := filepath.Join(tempDir, tempFile)
	sugar.Infof("Output json file: %s", outputFile)
	jsonData := Load(jsonnetFile, generateTestData)
	err := os.WriteFile(outputFile, []byte(jsonData), 0777)
	if err != nil {
		sugar.Panic(err.Error())
	}
}

func main() {
	start := time.Now()

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	jsonnetRootFolder, fileSearchPattern, outputFolder, generateTestData := parseFlags()

	sugar.Infof("jsonnet Root Folder: %s", jsonnetRootFolder)
	sugar.Infof("fileSearchPattern: %s", fileSearchPattern)
	sugar.Infof("Output Folder: %s", outputFolder)
	sugar.Infof("generateTestData: %t", generateTestData)

	processJsonnetFiles(jsonnetRootFolder, fileSearchPattern, outputFolder, generateTestData, sugar)

	duration := time.Since(start)
	sugar.Infof("Completed processing %v. Time in nano seconds: %d", os.Args, duration.Nanoseconds())
}
