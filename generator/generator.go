package generator

import (
	"com/parser/parser"
	"com/parser/ui"
	"com/parser/ui/components"
	"com/parser/utils"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	labelStart = "STARTED"
	labelFinish = "FINISHED"
	maximumTimestampOffsetS int64 = 10
	batchSize = 252144
)

func LogToTimestampStrings(log parser.Log) (string, string) {
	start := fmt.Sprintf(
		"{\"id\": %d, \"state\": \"%s\", \"timestamp\": %d}",
		log.Id, parser.StartFlag, log.TimestampStart.Unix())
	end := fmt.Sprintf(
		"{\"id\": %d, \"state\": \"%s\", \"timestamp\": %d}",
		log.Id, parser.FinishFlag, log.TimestampFinish.Unix())
	return start, end
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
	bytes uint64
}

func GenerateToFile(file *os.File, count int) {
	writeCh := make(chan parser.Log)
	progressCh := make(chan ProgressWrapper)

	timestamp := time.Now().UnixMilli()

	go func() {
		generateLogs(count, writeCh)
		close(writeCh)
	}()

	go func() {
		err := batchLogs(file, writeCh, progressCh)
		close(progressCh)
		if err != nil {
			fmt.Println(err)
			close(writeCh)
		}
	}()
	
	watchProgress(progressCh, count)

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
		createdBar: components.ProgressBarDefault(),
		writtenBar: components.ProgressBarDefault(),
		label: components.SimpleText{},
	}

	progressComponent.label.SetText("Writing timestamps to file...\n2 timestamps for each log")
	progressComponent.createdBar.SetPrefix("Created:")
	progressComponent.createdBar.SetMax(int64(total))
	progressComponent.writtenBar.SetPrefix("Written:")
	progressComponent.writtenBar.SetMax(int64(total))

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
	totalBytes := uint64(0)

	for u := range ch {
		if u.created > 0 {
			createdCount += u.created
			progressComponent.createdBar.SetValue(int64(createdCount))
		}
		if u.written > 0 {
			writtenCount += u.written
			progressComponent.writtenBar.SetValue(int64(createdCount))
		}
		if u.bytes > 0 {
			totalBytes += u.bytes
			progressComponent.writtenBar.SetSuffix("[%d/%d] [" + utils.FormatBytesToString(totalBytes) + "]")
		}
	}
	close(tickerCh)
	ui.ReRender(progressComponent)
}

func batchLogs(file *os.File, ch chan parser.Log, progressCh chan ProgressWrapper) error {
	var builder strings.Builder

	var logsBatch []string

	writeShuffleAndReset := func() (int, error) {
		rand.Shuffle(len(logsBatch), func(i, j int) {
			logsBatch[i], logsBatch[j] = logsBatch[j], logsBatch[i]
		})
		data := strings.Join(logsBatch, ",")
		builder.WriteString(data)
		bytes, err := file.WriteString(builder.String())
		if err != nil {
			fmt.Println("Error writing to file")
			return bytes, err
		}
		builder.Reset()
		builder.WriteString(",")
		return bytes, nil
	}

	for log := range ch {
		start, end := LogToTimestampStrings(log)
		logsBatch = append(logsBatch, start, end)
		if len(logsBatch) >= batchSize {
			bytes, err := writeShuffleAndReset()
			if err != nil {
				return err
			}
			progressCh <- ProgressWrapper{written: len(logsBatch), bytes: uint64(bytes)}
			logsBatch = logsBatch[:0]
		}
		progressCh <- ProgressWrapper{created: 1}
	}

	if len(logsBatch) > 0 {
		bytes, err := writeShuffleAndReset()
		if err != nil {
			return err
		}
		progressCh <- ProgressWrapper{written: len(logsBatch), bytes: uint64(bytes)}
	}
	return nil
}

func generateLogs(count int, ch chan parser.Log) {
	for i := 0; i < count; i++ {
		randomizedDelay := rand.Int63n(maximumTimestampOffsetS)
		randomizedOffset := rand.Int63n(maximumTimestampOffsetS) - maximumTimestampOffsetS / 2
		timestamp := time.Now().Add(time.Duration(randomizedOffset) * time.Second)
		ch <- parser.Log {
			Id: i,
			TimestampStart: timestamp,
			TimestampFinish: timestamp.Add(time.Duration(randomizedDelay) * time.Second),
		}
	}
}
