package generator

import (
	// "com/parser/ui"
	"com/parser/ui/components"
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
	logCh := make(chan Log)
	// progressCh := make(chan ProgressWrapper)
	var wg sync.WaitGroup

	timestamp := time.Now().UnixMilli()
	wg.Add(2)
	// wg.Add(3)
	// go func() {
	// 	defer wg.Done()
	// 	watchProgress(progressCh, count)
	// }()
	go func() {
		defer wg.Done()
		// defer close(progressCh)
		batchLogs(file, logCh/* , progressCh */)
	}()
	generateLogs(count, logCh)
	wg.Wait()

	elapsed := time.Now().UnixMilli() - timestamp

	fmt.Printf("\nGenerated %d logs in %d s %d ms", count, elapsed / 1000, elapsed % 1000)
}

// func watchProgress(ch chan ProgressWrapper, total int) {
// 	progressComponent := ProgressCustomComponent {
// 		createdBar: *components.ProgressBarDefault(),
// 		writtenBar: *components.ProgressBarDefault(),
// 		label: components.SimpleText{},
// 	}
//
// 	progressComponent.label.SetText("Writing logs to file...")
// 	progressComponent.createdBar.SetPrefix("Created logs:  ")
// 	progressComponent.createdBar.SetSuffix(fmt.Sprintf("0/%d", total))
// 	progressComponent.writtenBar.SetPrefix("Writen to file:")
// 	progressComponent.writtenBar.SetSuffix(fmt.Sprintf("0/%d", total))
//
// 	ui.Render(progressComponent)
//
// 	for u := range ch {
// 		if u.created > 0 {
// 			progressComponent.createdBar.SetSuffix(fmt.Sprintf("%d/%d", u.created, total))
// 			progressComponent.createdBar.SetPercentage(float32(u.created) / float32(total) * 100)
// 			ui.ReRender(progressComponent)
// 		}
// 		if u.written > 0 {
// 			progressComponent.writtenBar.SetSuffix(fmt.Sprintf("%d/%d", u.written, total))
// 			progressComponent.writtenBar.SetPercentage(float32(u.written) / float32(total) * 100)
// 			ui.ReRender(progressComponent)
// 		}
// 	}
//
// 	ui.ReRender(progressComponent)
// }

func batchLogs(file *os.File, ch chan Log/* , progressCh chan ProgressWrapper */) {
	var builder strings.Builder
	var logsBatch []string
	// logsWritten := 0
	// logsCollected := 0
	// lastUiUpdate := time.Now()
	// lastWrittenUiUpdate := time.Now()

	for log := range ch {
		logsBatch = append(logsBatch, log.String())
		// logsCollected ++
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

			// logsWritten += len(logsBatch)
			logsBatch = logsBatch[:0]
			// if time.Since(lastWrittenUiUpdate) > time.Millisecond * 200 {
			// 	progressCh <- ProgressWrapper{written: logsWritten, created: logsCollected}
			// }
		}
		// if time.Since(lastUiUpdate) > time.Millisecond * 200 {
		// 	lastUiUpdate = time.Now()
		// 	progressCh <- ProgressWrapper{created: logsCollected}
		// }
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
	// logsWritten += len(logsBatch)
	// progressCh <- ProgressWrapper{written: logsWritten, created: logsCollected}
	logsBatch = logsBatch[:0]
}

func generateLogs(count int, ch chan Log) {
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
