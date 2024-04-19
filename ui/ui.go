package ui

import (
	"bufio"
	"com/parser/appContext"
	"com/parser/events"
	"fmt"
	"os"
	"strings"
)

const (
	menuLine = "== Log parser == Possible actions\n1. Parse file\n2. Generate file\nQ. Quit\nInput: "
)

func Start() {
	eventHandler := appContext.EventHandler()
	scanner := bufio.NewScanner(os.Stdin)
	
	parseHandler := eventHandler.Subscribe(events.EventParse)
	generateHandler := eventHandler.Subscribe(events.EventGenerate)
	go func() {
		for {
			select {
			case event := <-parseHandler:
				fmt.Println("Parsing file", event)
			case event := <-generateHandler:
				fmt.Println("Generating file", event)
			}
		}
	}()

	for {
		fmt.Print(menuLine)
		scanner.Scan()
		input := strings.Trim(scanner.Text(), " ")
		switch strings.ToLower(input) {
		case "1":
			eventHandler.Publish(events.EventParse, nil)
		case "2":
			eventHandler.Publish(events.EventGenerate, nil)
		case "q", "quit":
			eventHandler.Publish(events.EventQuit, nil)
			return
		}
	}
}
