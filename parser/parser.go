package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func ReadAndFlag(reader io.Reader) error {
	timestamps := make([]LogTimestamp, 0)

	buffer := make([]byte, 1024)
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

		leftover := findTimestamps(mergedBuffer, &timestamps)

		clear(mergedBuffer)
		mergedBuffer = append(mergedBuffer, leftover...)
	}
	return err
}

func findTimestamps(bytes []byte, timestamps *[]LogTimestamp) []byte {
	openingBraceIndex := -1
	closeBraceIndex := -1
	var leftover []byte

	for i, b := range bytes {
		if b == '}' && openingBraceIndex != -1 {
			closeBraceIndex = i
			logTimestamp, err := unmarshalTimestamp(bytes[openingBraceIndex:i + 1])
			if err != nil {
				continue
			}
			*timestamps = append(*timestamps, logTimestamp)
		} else if b == '{' {
			openingBraceIndex = i
		}
	}

	if openingBraceIndex >= closeBraceIndex {
		leftover = append(leftover, bytes[openingBraceIndex:]...)
	}

	return leftover
}

func unmarshalTimestamp(bytes []byte) (LogTimestamp, error) {
	var logTimestamp LogTimestamp

	err := json.Unmarshal(bytes, &logTimestamp)
	if err != nil {
		return LogTimestamp{}, err
	} else if logTimestamp.Id == 0 || logTimestamp.Timestamp == 0 || logTimestamp.State == "" {
		return LogTimestamp{}, fmt.Errorf("Invalid log timestamp")
	}

	logTimestamp.State = strings.ToUpper(logTimestamp.State)
	return logTimestamp, err
}
