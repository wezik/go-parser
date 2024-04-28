package parser

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"testing"
	"time"
)

type mockReader struct {
	data []byte
	readErr error
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	if len(m.data) == 0 {
		return n, io.EOF
	}

	n = copy(p, m.data)
	m.data = m.data[n:]
	return n, nil
}

// --- Tests ---

func TestReadFile(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": 1619541248}, " +
			"{\"id\": 3, \"state\": \"FINISHED\", \"timestamp\": 1619541249}")
		testedReader := &mockReader{data: testedBytes}
		expectedTimestamps := []LogTimestamp {
			{Id: 1, Timestamp: 1619541248, State: StartFlag},
			{Id: 3, Timestamp: 1619541249, State: FinishFlag}}
		resultCh := make(chan LogTimestamp, 1024)
		
		// When
		err := readFile(testedReader, resultCh, 1024)
		close(resultCh)

		// Then
		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}
		if err != nil || !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Error: %v, Timestamps mismatch. Expected %v, got %v", err, expectedTimestamps, receivedTimestamps)
		}
	})

	t.Run("Split between buffers", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": 1619541248}, " +
			"{\"id\": 3, \"state\": \"FINISHED\", \"timestamp\": 1619541249}")
		testedReader := &mockReader{data: testedBytes}
		expectedTimestamps := []LogTimestamp {
			{Id: 1, Timestamp: 1619541248, State: StartFlag},
			{Id: 3, Timestamp: 1619541249, State: FinishFlag}}
		resultCh := make(chan LogTimestamp, 1024)
		
		// When
		err := readFile(testedReader, resultCh, 70)
		close(resultCh)

		// Then
		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}
		if err != nil || !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Error: %v, Timestamps mismatch. Expected %v, got %v", err, expectedTimestamps, receivedTimestamps)
		}
	})

	t.Run("Small buffer", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": 1619541248}, " +
			"{\"id\": 3, \"state\": \"FINISHED\", \"timestamp\": 1619541249}")
		testedReader := &mockReader{data: testedBytes}
		expectedTimestamps := []LogTimestamp {
			{Id: 1, Timestamp: 1619541248, State: StartFlag},
			{Id: 3, Timestamp: 1619541249, State: FinishFlag}}
		resultCh := make(chan LogTimestamp, 1024)
		
		// When
		err := readFile(testedReader, resultCh, 5)
		close(resultCh)

		// Then
		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}
		if err != nil || !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Error: %v, Timestamps mismatch. Expected %v, got %v", err, expectedTimestamps, receivedTimestamps)
		}
	})

	t.Run("Malformed file", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"p\": 1619541248}, " +
			"{\"id\": 3, \"state\": \"SHED\", \"")
		testedReader := &mockReader{data: testedBytes}
		expectedTimestamps := []LogTimestamp {}
		resultCh := make(chan LogTimestamp, 1024)
		
		// When
		err := readFile(testedReader, resultCh, 5)
		close(resultCh)

		// Then
		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}
		if err != nil || !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Error: %v, Timestamps mismatch. Expected %v, got %v", err, expectedTimestamps, receivedTimestamps)
		}
	})
}

func TestUnmarshalLogTimestamps(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": 1619541248}, " +
			"{\"id\": 3, \"state\": \"FINISHED\", \"timestamp\": 1619541249}")
		expectedTimestamps := []LogTimestamp {
			{Id: 1, Timestamp: 1619541248, State: StartFlag},
			{Id: 3, Timestamp: 1619541249, State: FinishFlag}}
		resultCh := make(chan LogTimestamp, len(expectedTimestamps))

		// When
		resultLeftover := unmarshalLogTimestamps(testedBytes, resultCh)

		// Then
		close(resultCh)
		if !bytes.Equal(resultLeftover, []byte{}) {
			t.Errorf("Leftover bytes mismatch. Expected %s, got %s", "", resultLeftover)
		}

		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}

		if !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Timestamps mismatch. Expected %v, got %v", expectedTimestamps, receivedTimestamps)
		}
	})

	t.Run("With leftover", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": 1619541248}, " +
			"{\"id\": 3, \"st")
		expectedLeftover := []byte("{\"id\": 3, \"st")
		expectedTimestamps := []LogTimestamp {{Id: 1, Timestamp: 1619541248, State: StartFlag}}
		resultCh := make(chan LogTimestamp, len(expectedTimestamps))

		// When
		resultLeftover := unmarshalLogTimestamps(testedBytes, resultCh)

		// Then
		close(resultCh)
		if !bytes.Equal(resultLeftover, expectedLeftover) {
			t.Errorf("Leftover bytes mismatch. Expected %s, got %s", "", resultLeftover)
		}

		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}

		if !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Timestamps mismatch. Expected %v, got %v", expectedTimestamps, receivedTimestamps)
		}
	})

	t.Run("With funny data", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1,\"sTAte\": \"started\", \"tIMEstamp\": 1619541248}, " +
			"{\"ID\":3,   \"staTE\": \"FiNisheD\",\"timeSTAMP\": 1619541249  } ,")
		expectedTimestamps := []LogTimestamp {
			{Id: 1, Timestamp: 1619541248, State: StartFlag},
			{Id: 3, Timestamp: 1619541249, State: FinishFlag}}
		resultCh := make(chan LogTimestamp, len(expectedTimestamps))

		// When
		resultLeftover := unmarshalLogTimestamps(testedBytes, resultCh)

		// Then
		close(resultCh)
		if !bytes.Equal(resultLeftover, []byte{}) {
			t.Errorf("Leftover bytes mismatch. Expected %s, got %s", "", resultLeftover)
		}

		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}

		if !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Timestamps mismatch. Expected %v, got %v", expectedTimestamps, receivedTimestamps)
		}
	})


	t.Run("Test malformed entry", func(t *testing.T) {
		// Given
		testedBytes := []byte(
			"{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": 1619541248}, " +
			"{\"id\": 3, \"st\": \"FINISHED\", \"timestamp\": 1619541249}")
		expectedTimestamps := []LogTimestamp {{Id: 1, Timestamp: 1619541248, State: StartFlag}}
		resultCh := make(chan LogTimestamp, len(expectedTimestamps))

		// When
		resultLeftover := unmarshalLogTimestamps(testedBytes, resultCh)

		// Then
		close(resultCh)
		if !bytes.Equal(resultLeftover, []byte{}) {
			t.Errorf("Leftover bytes mismatch. Expected %s, got %s", "", resultLeftover)
		}

		receivedTimestamps := make([]LogTimestamp, 0)
		for ts := range resultCh {
			receivedTimestamps = append(receivedTimestamps, ts)
		}

		if !reflect.DeepEqual(receivedTimestamps, expectedTimestamps) {
			t.Errorf("Timestamps mismatch. Expected %v, got %v", expectedTimestamps, receivedTimestamps)
		}
	})
}

func TestBytesToLogTimestamp(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": %d}", timestamp)
		bytes := []byte(str)

		expected := LogTimestamp {
			Id: 1,
			Timestamp: timestamp,
			State: StartFlag,
		}

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
	})

	t.Run("Test malformed id", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": \"1\", \"state\": \"STARTED\", \"timestamp\": %d}", timestamp)
		bytes := []byte(str)

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != (LogTimestamp{}) || err == nil {
			t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
		}
	})

	t.Run("Test malformed state", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": 1, \"ste\": \"STARTED\", \"timestamp\": %d}", timestamp)
		bytes := []byte(str)

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != (LogTimestamp{}) || err == nil {
			t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
		}
	})

	t.Run("Test malformed state 2", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": 1, \"state\": STARTED, \"timestamp\": %d}", timestamp)
		bytes := []byte(str)

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != (LogTimestamp{}) || err == nil {
			t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
		}
	})

	t.Run("Test malformed state 3", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": 1, \"state\": \"START\", \"timestamp\": %d}", timestamp)
		bytes := []byte(str)

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != (LogTimestamp{}) || err == nil {
			t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
		}
	})

	t.Run("Test malformed timestamp", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": 1, \"state\": \"STARTED\", \"timeamp\": %d}", timestamp)
		bytes := []byte(str)

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != (LogTimestamp{}) || err == nil {
			t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
		}
	})

	t.Run("Test malformed timestamp 2", func(t *testing.T) {
		// Given
		timestamp := time.Now().Unix()
		str := fmt.Sprintf("{\"id\": 1, \"state\": \"STARTED\", \"timestamp\": \"%d\"}", timestamp)
		bytes := []byte(str)

		// When
		result, err := bytesToLogTimestamp(bytes)

		// Then
		if result != (LogTimestamp{}) || err == nil {
			t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
		}
	})
}	

func TestCollectTimestamps(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		// Given
		timestampStart := time.Time.Add(time.Now(), time.Second * 17).Unix()
		timestampFinish := time.Time.Add(time.Now(), time.Second * 90).Unix()

		expectedLogs := []Log {
			{Id: 3, TimestampStart: time.Unix(timestampStart, 0), TimestampFinish: time.Unix(timestampFinish, 0)},
			{Id: 1, TimestampStart: time.Unix(timestampStart, 0), TimestampFinish: time.Unix(timestampFinish, 0)},
		}

		testStamps := []LogTimestamp {
			{Id: 3, Timestamp: timestampStart, State: StartFlag},
			{Id: 1, Timestamp: timestampStart, State: StartFlag},
			{Id: 3, Timestamp: timestampFinish, State: FinishFlag},
			{Id: 1, Timestamp: timestampFinish, State: FinishFlag},
		}

		tsCh := make(chan LogTimestamp)
		logCh := make(chan Log, len(testStamps))

		var wg sync.WaitGroup

		// When
		wg.Add(1)
		go func() {
			defer wg.Done()
			collectTimestamps(tsCh, logCh, 0)
		}()

		for _, ts := range testStamps {
			tsCh <- ts
		}

		close(tsCh)

		// Then
		wg.Wait()

		close(logCh)

		receivedLogs := make([]Log, 0)
		for log := range logCh {
			receivedLogs = append(receivedLogs, log)
		}

		if !reflect.DeepEqual(receivedLogs, expectedLogs) {
			t.Errorf("Logs mismatch. Expected %v, got %v", expectedLogs, receivedLogs)
		}
	})

	t.Run("Missing matches", func(t *testing.T) {
		// Given
		timestampStart := time.Time.Add(time.Now(), time.Second * 17).Unix()
		timestampFinish := time.Time.Add(time.Now(), time.Second * 90).Unix()

		expectedLogs := []Log {
			{Id: 1, TimestampStart: time.Unix(timestampStart, 0), TimestampFinish: time.Unix(timestampFinish, 0)},
		}

		testStamps := []LogTimestamp {
			{Id: 3, Timestamp: timestampStart, State: StartFlag},
			{Id: 1, Timestamp: timestampStart, State: StartFlag},
			{Id: 7, Timestamp: timestampFinish, State: FinishFlag},
			{Id: 1, Timestamp: timestampFinish, State: FinishFlag},
		}

		tsCh := make(chan LogTimestamp)
		logCh := make(chan Log, len(testStamps))

		var wg sync.WaitGroup

		// When
		wg.Add(1)
		go func() {
			defer wg.Done()
			collectTimestamps(tsCh, logCh, 0)
		}()

		for _, ts := range testStamps {
			tsCh <- ts
		}

		close(tsCh)

		// Then
		wg.Wait()

		close(logCh)

		receivedLogs := make([]Log, 0)
		for log := range logCh {
			receivedLogs = append(receivedLogs, log)
		}

		if !reflect.DeepEqual(receivedLogs, expectedLogs) {
			t.Errorf("Logs mismatch. Expected %v, got %v", expectedLogs, receivedLogs)
		}
	})

	t.Run("No entries", func(t *testing.T) {
		// Given
		tsCh := make(chan LogTimestamp)
		logCh := make(chan Log)

		var wg sync.WaitGroup

		// When
		wg.Add(1)
		go func() {
			defer wg.Done()
			collectTimestamps(tsCh, logCh, 0)
		}()

		close(tsCh)

		// Then
		wg.Wait()

		close(logCh)

		receivedLogs := make([]Log, 0)
		for log := range logCh {
			receivedLogs = append(receivedLogs, log)
		}

		if !reflect.DeepEqual(receivedLogs, []Log{}) {
			t.Errorf("Logs mismatch. Expected %v, got %v", []Log{}, receivedLogs)
		}
	})

	t.Run("With flag time", func(t *testing.T) {
		// Given
		tsStartShort := time.Time.Add(time.Now(), time.Second * 17).Unix()
		tsFinishShort := time.Time.Add(time.Now(), time.Second * 22).Unix()
		tsStartLong := time.Time.Add(time.Now(), time.Second * 17).Unix()
		tsFinishLong := time.Time.Add(time.Now(), time.Second * 23).Unix()

		expectedLogs := []Log {
			{Id: 1, TimestampStart: time.Unix(tsStartLong, 0), TimestampFinish: time.Unix(tsFinishLong, 0)},
		}

		testStamps := []LogTimestamp {
			{Id: 3, Timestamp: tsStartShort, State: StartFlag},
			{Id: 1, Timestamp: tsFinishLong, State: FinishFlag},
			{Id: 3, Timestamp: tsFinishShort, State: FinishFlag},
			{Id: 1, Timestamp: tsStartLong, State: StartFlag},
		}

		tsCh := make(chan LogTimestamp)
		logCh := make(chan Log, len(testStamps))

		var wg sync.WaitGroup

		// When
		wg.Add(1)
		go func() {
			defer wg.Done()
			collectTimestamps(tsCh, logCh, 6)
		}()

		for _, ts := range testStamps {
			tsCh <- ts
		}
		close(tsCh)

		// Then
		wg.Wait()

		close(logCh)

		receivedLogs := make([]Log, 0)
		for log := range logCh {
			receivedLogs = append(receivedLogs, log)
		}

		if !reflect.DeepEqual(receivedLogs, expectedLogs) {
			t.Errorf("Logs mismatch. Expected %v, got %v", expectedLogs, receivedLogs)
		}
	})
}
