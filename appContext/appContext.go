package appContext

import "com/parser/events"

var globalEventHandler *events.EventHandler

func init() {
	 globalEventHandler = events.NewEventHandler()
}

func EventHandler() *events.EventHandler {
	return globalEventHandler
}
