package helper

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

type Message struct {
	Message string
	Variant string
}

var exportDir string

func CreateDataDirectory() error {
	exportDir = filepath.Join(".", "exports")
	return os.MkdirAll(exportDir, os.ModePerm)
}

func CreateCSV(filename string, data [][]string) (string, error) {
	filename = filepath.Join(exportDir, filename)
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return "", err
		}
	}
	fp, _ := filepath.Abs(filename)
	return fp + ".csv", nil
}
