package ui

import (
	"bufio"
	"com/parser/logParser"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
menuLine = 
`== Log parser == Possible actions
1. Parse file
2. Generate file
Q. Quit
Input: `
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go watchProgress(&wg)

	shouldQuit := false

	for !shouldQuit {
		switch readUserInput(scanner) {
		case "1":
			logParser.Parse()
		case "2":

		case "q", "quit":
			shouldQuit = true	
		}
	}
}

func readUserInput(scanner *bufio.Scanner) string {
	fmt.Print(menuLine)
	scanner.Scan()
	input := strings.Trim(scanner.Text(), " ")
	return strings.ToLower(input)
}
