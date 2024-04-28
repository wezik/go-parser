package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func Parse(fileName string) error {

	file, err := os.Open(fileName); if err != nil {
		return err
	}
	defer file.Close()

	timestampsCh := make(chan LogTimestamp)
	logsCh := make(chan Log)

	go func() {
		defer close(timestampsCh)
		_ = readFile(file, timestampsCh, 1024)
	}()

	go func() {
		defer close(logsCh)
		collectTimestamps(timestampsCh, logsCh, 8)
	}()

	i := 0
	fmt.Println("Benchmark start")
	start := time.Now()
	fmt.Println("Processing logs...")
	for log := range logsCh {
		i++
		// Placeholder consume for now
		_ = log
		if (i % 16384 == 0) {
			fmt.Print("\rProcessed ", i, " logs")
		}
	}
	fmt.Print("\rFound ", i, " logs flagged over 8 seconds\n")
	fmt.Println("Benchmark end")
	fmt.Println("Elapsed time: ", time.Since(start))

	return nil
}

func readFile(reader io.Reader, ch chan LogTimestamp, bufferSize int) error {
	buffer := make([]byte, bufferSize)
	mergedBuffer := make([]byte, 0)

	var err error

	for {
		clear(buffer)
		_, err = reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}

		mergedBuffer = append(mergedBuffer, buffer...)

		leftover := unmarshalLogTimestamps(mergedBuffer, ch)

		mergedBuffer = append([]byte{}, leftover...)
	}
	return err
}

func collectTimestamps(receiveCh chan LogTimestamp, sendCh chan Log, delayFlag int64) {
	tsMap := make(map[int] LogTimestamp)
	for tsFromCh := range receiveCh {
		if tsFromMap, found := tsMap[tsFromCh.Id]; found {
			delay := tsFromCh.Timestamp - tsFromMap.Timestamp
			if delay >= delayFlag || delay <= (delayFlag * -1) {
				var tsStart, tsFinish time.Time
				if tsFromMap.State == StartFlag {
					tsStart = time.Unix(tsFromMap.Timestamp, 0)
					tsFinish = time.Unix(tsFromCh.Timestamp, 0)
				} else {
					tsStart = time.Unix(tsFromCh.Timestamp, 0)
					tsFinish = time.Unix(tsFromMap.Timestamp, 0)
				}
				sendCh <- Log {
					Id: tsFromCh.Id,
					TimestampStart: tsStart,
					TimestampFinish: tsFinish,
				}
			}
			delete(tsMap, tsFromMap.Id)
		} else {
			tsMap[tsFromCh.Id] = tsFromCh 
		}
	}
}

func unmarshalLogTimestamps(bytes []byte, ch chan LogTimestamp) []byte {
	openingBraceIndex := -1
	closeBraceIndex := -1

	for i, b := range bytes {
		if b == '}' && openingBraceIndex != -1 {
			closeBraceIndex = i
			logTimestamp, err := bytesToLogTimestamp(bytes[openingBraceIndex:i + 1])
			if err != nil {
				continue
			}
			ch <- logTimestamp
		} else if b == '{' {
			openingBraceIndex = i
		}
	}

	if openingBraceIndex >= closeBraceIndex {
		return append([]byte{}, bytes[openingBraceIndex:]...)
	}

	return []byte{}
}

func bytesToLogTimestamp(bytes []byte) (LogTimestamp, error) {
	var logTimestamp LogTimestamp

	err := json.Unmarshal(bytes, &logTimestamp)
	if err != nil {
		return LogTimestamp{}, err
	} else if logTimestamp.Timestamp == 0 || logTimestamp.State == "" {
		return LogTimestamp{}, fmt.Errorf("Invalid log timestamp")
	}

	logTimestamp.State = strings.ToUpper(logTimestamp.State)
	if (logTimestamp.State != StartFlag && logTimestamp.State != FinishFlag) {
		return LogTimestamp{}, fmt.Errorf("Invalid log state")
	}

	return logTimestamp, err
}
