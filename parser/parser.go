package parser

import (
	"fmt"
	"io"
	"strings"
	"time"
)



func ReadAndFlag(reader io.Reader) {
	timestamp := time.Now().UnixMilli()
	buffer := make([]byte, 1024)
	logs := make([]string, 0)
	builder := strings.Builder{}
	for {
		_, err := reader.Read(buffer)
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
	}
	fmt.Println("Timestamps read:", len(logs))
	fmt.Println("ReadAndFlag took", time.Now().UnixMilli()-timestamp, "ms")
}

