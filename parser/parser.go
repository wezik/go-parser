package parser

import (
	"com/parser/ui"
	"com/parser/ui/components"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type ProgressWrapper struct {
	bytes int64
}

type ProgressCustomComponent struct {
	bytesBar components.ProgressBar
	label components.SimpleText
}

func (pw ProgressCustomComponent) String() string {
	return pw.label.String() + "\n" + pw.bytesBar.String() + "\n"
}

func ReadAndFlag(file *os.File) {
	timestamp := time.Now().UnixMilli()
	buffer := make([]byte, 1024)
	logs := make([]string, 0)
	builder := strings.Builder{}
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	progressCh := make(chan ProgressWrapper)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		watchProgress(progressCh, fileStat.Size())
	}()
	for {
		bytes, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		startIndex := 0
		endIndex := 0
		for i, b := range buffer {
			if b == '}' {
				endIndex = i
				builder.WriteString(string(buffer[startIndex:endIndex + 1]))
				logs = append(logs, builder.String())
				builder.Reset()
			} else if b == '{' {
				startIndex = i
			}
		}
		// Handle logs split between buffers
		if startIndex >= endIndex {
			builder.WriteString(string(buffer[startIndex:]))
		}
		// Buffer can be filled with previous read, so we need to clear it
		clear(buffer)
		progressCh <- ProgressWrapper{bytes: int64(bytes)}
	}
	close(progressCh)
	wg.Wait()
	fmt.Println("Timestamps read:", len(logs))
	fmt.Println("ReadAndFlag took", time.Now().UnixMilli()-timestamp, "ms")
}

func watchProgress(progressCh chan ProgressWrapper, fileSize int64) {
	component := ProgressCustomComponent{
		bytesBar: components.ProgressBarDefault(),
		label: components.SimpleText{},
	}
	component.label.SetText("Reading file...")	
	component.bytesBar.SetMax(fileSize)

	ui.Render(component)

	tickerStopCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		defer wg.Done()
		for {
			select {	
			case <- ticker.C:
				ui.ReRender(component)
			case <- tickerStopCh:
				return
			}
		}
	}()

	totalBytes := int64(0)
	for progress := range progressCh {
		totalBytes += progress.bytes
		component.bytesBar.SetValue(totalBytes)
	}
	close(tickerStopCh)
	wg.Wait()
	ui.ReRender(component)
}
