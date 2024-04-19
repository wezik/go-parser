package ui

import (
	"bufio"
	"com/parser/logParser"
	"com/parser/testGenerator"
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

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go watchProgress(&wg)

	shouldQuit := false

	for !shouldQuit {
		switch readUserInput(menuLine, scanner) {
		case "1":
			logParser.Parse()
		case "2":
			logCountString := readUserInput(askLogCount, scanner)
			fileName := readUserInput(askFileName, scanner)
			logCount, err := strconv.Atoi(logCountString)
			if err != nil {
				continue
			}
			testGenerator.Run(logCount, fileName)
		case "q", "quit":
			shouldQuit = true	
		}
	}
}

func readUserInput(prefix string, scanner *bufio.Scanner) string {
	fmt.Print(prefix)
	scanner.Scan()
	input := strings.Trim(scanner.Text(), " ")
	return strings.ToLower(input)
}
