package logParser

import (
	"com/parser/appContext"
	"com/parser/events"
	"fmt"
	"time"
)

func Parse() {
	eh := appContext.EventHandler()
	eh.Publish(events.EventProgressStart, "Tracking progress on parsing")
	for i := 0.0; i < 100; i += 1.24 {
		eh.Publish(events.EventProgressDraw, fmt.Sprintf("\rLoading: %.2f%%", i))
		time.Sleep(50 * time.Millisecond)
	}
	eh.Publish(events.EventProgressDraw, fmt.Sprintf("\rLoading: 100.00%%"))
	eh.Publish(events.EventProgressComplete, "Tracking complete")
}
