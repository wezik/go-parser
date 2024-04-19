package events 

type EventType int

const (
	EventParse EventType = iota
	EventGenerate
	EventQuit
)

func (e *Event) EventType() EventType {
	return e.eventType
}

func (e *Event) EventData() interface{} {
	return e.eventData
}

func (e *Event) SetType(t EventType) {
	e.eventType = t
}

func (e *Event) SetData(d interface{}) {
	e.eventData = d
}

type Event struct {
	eventType EventType
	eventData interface{}
}

type EventHandler struct {
	subscribers map[EventType][]chan interface{}
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		subscribers: make(map[EventType][]chan interface{}),
	}
}

func (eh *EventHandler) Subscribe(t EventType) chan interface{} {
	ch := make(chan interface{})
	eh.subscribers[t] = append(eh.subscribers[t], ch)
	return ch
}

func (eh *EventHandler) Publish(t EventType, d interface{}) {
	event := Event{}
	event.SetType(t)
	event.SetData(d)
	for _, ch := range eh.subscribers[t] {
		ch <- event // Might be worth to make it into a goroutine
	}
}
