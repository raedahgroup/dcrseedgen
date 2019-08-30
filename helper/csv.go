package helper

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"time"
)

const (
	EXPORT_FOLDER_NAME = "exports"
)

var (
	exportFolder string
)

func CreateExportFolder() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	exportFolder = filepath.Join(wd, EXPORT_FOLDER_NAME)
	if _, err := os.Stat(exportFolder); os.IsNotExist(err) {
		err := os.MkdirAll(exportFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateCSV(data [][]string) (string, error) {
	filename := filepath.Join(exportFolder, time.Now().Format("2006-01-02 15:04:05")+".csv")

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	csvData := make([][]string, len(data)+1)
	csvData[0] = []string{"Address", "Private Key"}

	for index, item := range data {
		csvData[index+1] = item
	}

	for _, value := range csvData {
		err := writer.Write(value)
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}
