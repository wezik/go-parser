package events 

type EventType int

const (
	EventProgressStart EventType = iota
	EventProgressDraw
	EventProgressComplete
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
	subscribers map[EventType][]chan Event
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		subscribers: make(map[EventType][]chan Event),
	}
}

func (eh *EventHandler) Subscribe(t EventType) chan Event {
	ch := make(chan Event)
	eh.subscribers[t] = append(eh.subscribers[t], ch)
	return ch
}

func (eh *EventHandler) Unsubscribe(t EventType, ch chan Event) {
	if subscribers, ok := eh.subscribers[t]; ok {
		for i, subscriber := range subscribers {
			if subscriber == ch {
				eh.subscribers[t] = append(subscribers[:i], subscribers[i+1:]...)
				close(ch)
				break
			}
		}
	}
}

func (eh *EventHandler) Publish(t EventType, d interface{}) {
	event := Event{}
	event.SetType(t)
	event.SetData(d)
	for _, ch := range eh.subscribers[t] {
		ch <- event // Might be worth to make it into a goroutine
	}
}
