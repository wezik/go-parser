package generator

import (
	"com/parser/appContext"
	"com/parser/events"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	labelStart = "STARTED"
	labelFinish = "FINISHED"
	maximumOffsetMs int64 = 10000
	shuffleSize = 252144
)

type Log struct {
	id int
	state string
	timestamp int64
}

func (l *Log) String() string {
	return fmt.Sprintf("{id:%d, state:\"%s\", timestamp:%d}", l.id, l.state, l.timestamp)
}

func GenerateToFile(file *os.File, logCount int) {
	eh := appContext.EventHandler()
	ch := make(chan Log)
	var wg sync.WaitGroup

	eh.Publish(events.EventProgressStart, "Generating logs to file")

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer close(ch)
		generateLogs(logCount, ch)
	}()
	go func() {
		defer wg.Done()
		batchLogs(file, logCount, ch, eh)
	}()
	wg.Wait()

	eh.Publish(events.EventProgressComplete, "Finished generating logs to file")
}

func batchLogs(file *os.File, count int, ch chan Log, eh *events.EventHandler) {
	var builder strings.Builder
	var logsBatch []string
	logsWritten := 0

	for log := range ch {
		logsBatch = append(logsBatch, log.String())
		if len(logsBatch) >= shuffleSize {
			rand.Shuffle(len(logsBatch), func(i, j int) {
				logsBatch[i], logsBatch[j] = logsBatch[j], logsBatch[i]
			})
			logsStringifed := strings.Join(logsBatch, ",")
			builder.WriteString(logsStringifed)
			_, err := file.WriteString(builder.String())
			if err != nil {
				fmt.Println("Error writing to file")
				return
			}
			builder.Reset()
			builder.WriteString(",")

			logsWritten += len(logsBatch)
			logsBatch = logsBatch[:0]

			eh.Publish(events.EventProgressDraw, formatProgressMessage(logsWritten, count))
		}
	}
	if len(logsBatch) <= 0 {
		return
	}
	rand.Shuffle(len(logsBatch), func(i, j int) {
		logsBatch[i], logsBatch[j] = logsBatch[j], logsBatch[i]
	})
	logsStringifed := strings.Join(logsBatch, ",")
	builder.WriteString(logsStringifed)
	_, err := file.WriteString(builder.String())
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}
	logsWritten += len(logsBatch)
	logsBatch = logsBatch[:0]

	eh.Publish(events.EventProgressDraw, formatProgressMessage(logsWritten, count))
}

func formatProgressMessage(logsWritten, count int) string {
	percentage := float32(logsWritten) / float32(count) * 100
	return fmt.Sprintf("\rLogs written: %d/%d %.2f%%", logsWritten, count, percentage)
}

func generateLogs(count int, ch chan Log) {
	for i := 0; i < count / 2; i++ {
		generatedDelay := rand.Int63n(maximumOffsetMs)
		generatedOffset := rand.Int63n(maximumOffsetMs)
		startLog := Log {
			id: i,
			state: labelStart,
			timestamp: time.Now().Unix() + generatedOffset,
		}
		ch <- startLog
		endLog := Log {
			id: startLog.id,
			state: labelFinish,
			timestamp: startLog.timestamp + generatedDelay,
		}
		ch <- endLog
	}
}
