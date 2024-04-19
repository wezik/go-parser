package generator

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
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
	fmt.Println("Writing to file")
	ch := make(chan Log)
	go generateLogsRoutine(logCount, ch)
	batchLogs(file, ch)
}

func batchLogs(file *os.File, ch chan Log) {
	var builder strings.Builder
	var logsBatch []string

	for log := range ch {
		logsBatch = append(logsBatch, log.String())
		if len(logsBatch) >= shuffleSize {
			rand.Shuffle(len(logsBatch), func(i, j int) {
				logsBatch[i], logsBatch[j] = logsBatch[j], logsBatch[i]
			})
			logsStringifed := strings.Join(logsBatch, ",")
			logsBatch = logsBatch[:0]
			builder.WriteString(logsStringifed)
			_, err := file.WriteString(builder.String())
			if err != nil {
				fmt.Println("Error writing to file")
				return
			}
			builder.Reset()
			builder.WriteString(",")
		}
	}
	if len(logsBatch) <= 0 {
		return
	}
	rand.Shuffle(len(logsBatch), func(i, j int) {
		logsBatch[i], logsBatch[j] = logsBatch[j], logsBatch[i]
	})
	logsStringifed := strings.Join(logsBatch, ",")
	logsBatch = logsBatch[:0]
	builder.WriteString(logsStringifed)
	_, err := file.WriteString(builder.String())
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}
}

func generateLogsRoutine(count int, ch chan Log) {
	defer close(ch)
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
