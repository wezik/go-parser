package eventHandler

import "fmt"

type Event struct {
	eventType EventType
	eventData interface{}
}

type EventType int

const (
	EventParse EventType = iota
	EventGenerate
	EventQuit
)

func (e *Event) SetType(t EventType) {
	e.eventType = t
}

func (e *Event) SetData(d interface{}) {
	e.eventData = d
}

func (e *Event) EventType() EventType {
	return e.eventType
}

func (e *Event) EventData() interface{} {
	return e.eventData
}

func SendEvent(t EventType) {
	SendEventWithData(t, nil)
}

func SendEventWithData(t EventType, d interface{}) {
	e := Event{}
	e.SetType(t)
	e.SetData(d)
	handleEvent(e)
}

func handleEvent(e Event) {
	switch e.EventType() {
	case EventParse:
		// Parse file
		fmt.Println("Parsing file")
	case EventGenerate:
		// Generate file
		fmt.Println("Generating file")
	case EventQuit:
		fmt.Println("Quitting application")
	}
}
