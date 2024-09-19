package main

import (
	"os"

	"time"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func main() {
	start := time.Now()

	// Initialize logger
	rawLogger, _ := zap.NewProduction()
	defer rawLogger.Sync() // flushes buffer, if any
	Logger := rawLogger.Sugar()

	jsonnetRootFolder, fileSearchPattern, outputFolder, generateTestData := parseFlags()

	Logger.Infof("jsonnet Root Folder: %s", jsonnetRootFolder)
	Logger.Infof("fileSearchPattern: %s", fileSearchPattern)
	Logger.Infof("Output Folder: %s", outputFolder)
	Logger.Infof("generateTestData: %t", generateTestData)

	processJsonnetFiles(jsonnetRootFolder, fileSearchPattern, outputFolder, generateTestData)

	duration := time.Since(start)
	Logger.Infof("Completed processing %v. Time in nano seconds: %d", os.Args, duration.Nanoseconds())
}
