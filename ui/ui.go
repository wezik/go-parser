package ui

import (
	"bufio"
	"com/parser/eventHandler"
	"fmt"
	"os"
	"strings"
)

const (
	menuLine = "== Log parser == Possible actions\n1. Parse file\n2. Generate file\nQ. Quit\nInput: "
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(menuLine)
		scanner.Scan()
		input := strings.Trim(scanner.Text(), " ")
		switch strings.ToLower(input) {
		case "1":
			eventHandler.SendEvent(eventHandler.EventParse, nil)
		case "2":
			eventHandler.SendEvent(eventHandler.EventGenerate, nil)
		case "q", "quit":
			eventHandler.SendEvent(eventHandler.EventQuit, nil)
			return
		}
	}
}
