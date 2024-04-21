package generator

import (
	"com/parser/ui"
	"com/parser/ui/components"
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
	batchSize = 252144
)

type Log struct {
	id int
	state string
	timestamp int64
}

func (l *Log) String() string {
	return fmt.Sprintf("{id:%d, state:\"%s\", timestamp:%d}", l.id, l.state, l.timestamp)
}

type ProgressCustomComponent struct {
	createdBar components.ProgressBar
	writtenBar components.ProgressBar
	label components.SimpleText
}

func (pw ProgressCustomComponent) String() string {
	return pw.label.String() + "\n" + pw.createdBar.String() + "\n" + pw.writtenBar.String() + "\n"
}

type ProgressWrapper struct {
	created int
	written int
}

func GenerateToFile(file *os.File, count int) {
	writeCh := make(chan Log)
	progressCh := make(chan ProgressWrapper)

	timestamp := time.Now().UnixMilli()

	go func() {
		generateLogs(count, writeCh)
		close(writeCh)
	}()

	go func() {
		batchLogs(file, writeCh, progressCh)
		close(progressCh)
	}()
	
	watchProgress(progressCh, count * 2)

	timeElapsed := time.Since(time.UnixMilli(timestamp)).Milliseconds()
	timeElapsedComponent := components.SimpleText{}
	timeElapsedComponent.SetText(formatElapsedString(timeElapsed))
	ui.Render(timeElapsedComponent)
}

func formatElapsedString(elapsed int64) string {
	var b strings.Builder
	ms := elapsed % 1000
	s := (elapsed / 1000) % 60
	m := (elapsed / 1000) / 60
	b.WriteString("Completed in [")
	if m > 0 {
		b.WriteString(fmt.Sprintf("%dm ", m))
	}
	if s > 0 {
		b.WriteString(fmt.Sprintf("%ds ", s))
	}
	b.WriteString(fmt.Sprintf("%dms]\n", ms))
	return b.String()
}


func watchProgress(ch chan ProgressWrapper, total int) {
	progressComponent := ProgressCustomComponent {
		createdBar: *components.ProgressBarDefault(),
		writtenBar: *components.ProgressBarDefault(),
		label: components.SimpleText{},
	}

	progressComponent.label.SetText("Writing timestamps to file...\n2 timestamps for each log")
	progressComponent.createdBar.SetPrefix("Created:")
	progressComponent.createdBar.SetMax(total)
	progressComponent.writtenBar.SetPrefix("Written:")
	progressComponent.writtenBar.SetMax(total)

	ui.Render(progressComponent)
	
	tickerCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <- ticker.C:
				ui.ReRender(progressComponent)
			case <- tickerCh:
				return
			}

		}
	}()

	createdCount := 0
	writtenCount := 0

	for u := range ch {
		if u.created > 0 {
			createdCount += u.created
			progressComponent.createdBar.SetValue(createdCount)
		}
		if u.written > 0 {
			writtenCount += u.written
			progressComponent.writtenBar.SetValue(writtenCount)
		}
	}
	close(tickerCh)
	ui.ReRender(progressComponent)
}

func batchLogs(file *os.File, ch chan Log, progressCh chan ProgressWrapper) {
	var builder strings.Builder

	var logsBatch []string

	writeShuffleAndReset := func() {
		// rand.Shuffle(len(logsBatch), func(i, j int) {
		// 	logsBatch[i], logsBatch[j] = logsBatch[j], logsBatch[i]
		// })
		data := strings.Join(logsBatch, ",")
		builder.WriteString(data)
		_, err := file.WriteString(builder.String())
		if err != nil {
			fmt.Println("Error writing to file")
			return
		}
		builder.Reset()
		builder.WriteString(",")
	}

	for log := range ch {
		logsBatch = append(logsBatch, log.String())
		if len(logsBatch) >= batchSize {
			writeShuffleAndReset()
			progressCh <- ProgressWrapper{written: len(logsBatch)}
			logsBatch = logsBatch[:0]
		}
		progressCh <- ProgressWrapper{created: 1}
	}

	if len(logsBatch) > 0 {
		writeShuffleAndReset()
		progressCh <- ProgressWrapper{written: len(logsBatch)}
	}
}

func generateLogs(count int, ch chan Log) {
	for i := 0; i < count; i++ {
		randomizedDelay := rand.Int63n(maximumOffsetMs)
		randomizedOffset := rand.Int63n(maximumOffsetMs) - maximumOffsetMs / 2
		startLog := Log {
			id: i,
			state: labelStart,
			timestamp: time.Now().Unix() + randomizedOffset,
		}
		endLog := Log {
			id: startLog.id,
			state: labelFinish,
			timestamp: startLog.timestamp + randomizedDelay,
		}
		ch <- startLog
		ch <- endLog
	}
}
