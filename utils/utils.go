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
