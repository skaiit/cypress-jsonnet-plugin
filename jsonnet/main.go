package main

import (
	"os"

	"time"

	"go.uber.org/zap"
)

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
