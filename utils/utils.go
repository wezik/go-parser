package utils

import (
	"fmt"
	"os"
)

func TerminalWidth() int {
	return 200 // Todo implement
}

func CreateFile(fileName string) (*os.File, error) {
	var _, err = os.Stat(fileName)
	var file *os.File

	if os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			return file, err
		}
	} else {
		return file, fmt.Errorf("File with that name already exists")
	}
	return file, nil
}

func OpenFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return file, err
	}
	return file, nil
}

func FormatBytesToString(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024 * 1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes) / 1024)
	} else if bytes < 1024 * 1024 * 1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes) / 1024 / 1024)
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes) / 1024 / 1024 / 1024)
	}
}
