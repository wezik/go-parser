package utils

import (
	"fmt"
	"os"
)

func CreateFile(fileName string) (*os.File, error) {
	var _, err = os.Stat(fileName)
	var file *os.File

	if os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file")
			return file, err
		}

	} else {
		fmt.Println("File already exists")
		return file, err
	}

	fmt.Println("File created successfully")
	return file, nil
}
