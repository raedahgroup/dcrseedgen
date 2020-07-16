package helper

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"time"
)

type Message struct {
	Message string
	Variant string
}

const (
	exportDir            = "./exports"
	exportFilenamePrefix = "dcrseedgen_"
)

func CreateDataDirectory() error {
	_, err := os.Stat(exportDir)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		os.Mkdir(exportDir, os.ModePerm)
		return nil
	}
	return err
}

func CreateCSV(data [][]string) (string, error) {
	filename := filepath.Join(exportDir, exportFilenamePrefix+time.Now().Format("2006-01-02_15:04:05")+".csv")
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
	return fp, nil
}
