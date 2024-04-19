package ui

import (
	"com/parser/appContext"
	"com/parser/events"
	"fmt"
	"sync"
)

func watchProgress(wg *sync.WaitGroup) {
	defer wg.Done()

	eh := appContext.EventHandler()

	startEventCh := eh.Subscribe(events.EventProgressStart)
	defer eh.Unsubscribe(events.EventProgressStart, startEventCh)

	for startEvent := range startEventCh {
		wg.Add(1)
		go listenToProgress(wg, eh, startEvent)
	}
}

func listenToProgress(wg *sync.WaitGroup, eh *events.EventHandler, startEvent events.Event) {
	defer wg.Done()
	
	fmt.Println(startEvent.EventData())

	eventCh := eh.Subscribe(events.EventProgressDraw)
	defer eh.Unsubscribe(events.EventProgressDraw, eventCh)
	
	completeCh := eh.Subscribe(events.EventProgressComplete)
	defer eh.Unsubscribe(events.EventProgressComplete, completeCh)
	
	for {
		select {
		case event := <- eventCh:
			fmt.Print(event.EventData())
		case event := <- completeCh:
			fmt.Printf("\n%v\n", event.EventData())
			return
		}
	}
}
