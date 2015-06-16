package main

import (
	"os"
	"path/filepath"
	"strings"
)

func isTmpDir(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(settings.outputPath())

	return absolutePath == absoluteTmpPath
}

func isWatchedFile(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(settings.outputPath())

	if strings.HasPrefix(absolutePath, absoluteTmpPath) {
		return false
	}

	name := filepath.Base(path)
	ext := filepath.Ext(path)

	for _, e := range strings.Split(settings["valid_ext"], ",") {
		validExt := strings.TrimSpace(e)
		if validExt == ext || (ext == "" && validExt == name) {
			return true
		}
	}

	return false
}

func createBuildErrorsLog(message string) bool {
	file, err := os.Create(settings.buildErrorsFilePath())
	if err != nil {
		return false
	}

	_, err = file.WriteString(message)
	if err != nil {
		return false
	}

	return true
}

func removeBuildErrorsLog() error {
	err := os.Remove(settings.buildErrorsFilePath())

	return err
}
