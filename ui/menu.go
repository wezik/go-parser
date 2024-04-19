package ui

import (
	"bufio"
	"com/parser/logParser"
	"com/parser/generator"
	"com/parser/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	menuLine = 
		"== Log parser == Possible actions:\n" +
		"1. Parse file\n" +
		"2. Generate file\n" +
		"Q. Quit\n" +
		"Input: "
	
	askLogCount = "Amount of logs to write: "
	askFileName = "File name: "
)

var scanner *bufio.Scanner

func init() {
	scanner = bufio.NewScanner(os.Stdin)
}

func Start() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go watchProgress(&wg)

	shouldQuit := false

	for !shouldQuit {
		switch readUserInput(menuLine) {
		case "1":
			logParser.Parse()
		case "2":
			err := handleInputGenerate()
			if err != nil {
				continue
			}
		case "q", "quit":
			shouldQuit = true	
		}
	}
}

func readUserInput(prefix string) string {
	fmt.Print(prefix)
	scanner.Scan()
	input := strings.Trim(scanner.Text(), " ")
	return strings.ToLower(input)
}

func handleInputGenerate() error {
	logCountString := readUserInput(askLogCount)
	fileName := readUserInput(askFileName)

	file, err := utils.CreateFile(fileName)
	if err != nil {
		fmt.Println("Error when creating the file")
		return err
	}
	defer file.Close()

	logCount, err := strconv.Atoi(logCountString)
	if err != nil {
		fmt.Println("Error when parsing log count, incorrect format")
		return err
	}

	generator.GenerateToFile(file, logCount)
	
	fmt.Println("Generation success")
	return nil
}
