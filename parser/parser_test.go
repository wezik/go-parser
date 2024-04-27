package parser

import (
	"fmt"
	// "io"
	"os"
	"testing"
	"time"
)

// --- Mocks ---

// type mockReader struct {
// 	data []byte
// 	readErr error
// }
//
// func (m *mockReader) Read(p []byte) (n int, err error) {
// 	if m.readErr != nil {
// 		return 0, m.readErr
// 	}
//
// 	finalLength := 70 // Capped at 70 bytes so we can test leftovers
// 	if len(m.data) < finalLength {
// 		finalLength = len(m.data)
// 	}
//
// 	if finalLength == 0 {
// 		return n, io.EOF
// 	}
//
// 	n = copy(p, m.data[:finalLength])
// 	fmt.Println("Read:", string(p))
// 	m.data = m.data[n:]
// 	return n, nil
// }

const (
	mockData = `{"Id": 1, "STATE": "STARTED", "tiMesTamP": 1619541248} , {"iD": 3, "STatE": "FINISHED", "tiMesTamP": 1619541249}`
)

var mockDataBytes []byte
var mockDataTimestamps []LogTimestamp

func TestMain(m *testing.M) {
	fmt.Println("Running setup before tests")

	mockDataBytes = []byte(mockData)
	mockDataTimestamps = make([]LogTimestamp, 0)
	mockDataTimestamps = append(mockDataTimestamps, LogTimestamp{Id: 1, Timestamp: 1619541248, State: StartFlag})
	mockDataTimestamps = append(mockDataTimestamps, LogTimestamp{Id: 3, Timestamp: 1619541249, State: FinishFlag})

	exitCode := m.Run()

	fmt.Println("Running teardown after tests")

	os.Exit(exitCode)
}

// --- FindTimestamps ---

func TestFindTimestamps(t *testing.T) {
	// Given

	// When
	resultTimestamps := make([]LogTimestamp, 0)
	resultLeftover := findTimestamps(mockDataBytes, &resultTimestamps)

	// Then
	if len(resultTimestamps) < 1 || len(resultTimestamps) != len(mockDataTimestamps) {
		fmt.Println(mockDataBytes)
		t.Errorf("Timestamps length Expected %v, got %v", len(mockDataTimestamps), len(resultTimestamps))
	}
	if string(resultLeftover) != "" {
		t.Errorf("Leftover Expected %v, got %v", "nothing", resultLeftover)
	}
	if (len(resultTimestamps) == len(mockDataTimestamps)) {
		for i, expected := range mockDataTimestamps {
			if resultTimestamps[i] != expected {
				t.Errorf("Timestamps Expected %v, got %v", expected, resultTimestamps[i])
			}
		}
	}
}

func TestFindTimestampsWithLeftover(t *testing.T) {
	// Given
	bytes := []byte(`{"Id": 1, "STATE": "STARTED", "tiMesTamP": 1619541248} , {"iD": 3, "STatE": "`)
	resultTimestamps := make([]LogTimestamp, 0)
	expectedTimestamps := make([]LogTimestamp, 0)

	expectedTimestamps = append(expectedTimestamps, LogTimestamp{Id: 1, Timestamp: 1619541248, State: StartFlag})
	expectedLeftover := []byte(`{"iD": 3, "STatE": "`)

	// When
	resultLeftover := findTimestamps(bytes, &resultTimestamps)

	// Then
	if len(resultTimestamps) < 1 || resultTimestamps[0] != expectedTimestamps[0] || len(resultTimestamps) > 1 {
		t.Errorf("Found timestamps Expected %v, got %v", expectedTimestamps, resultTimestamps)
	}
	if string(resultLeftover) != string(expectedLeftover) {
		t.Errorf("Leftover Expected %v, got %v", expectedLeftover, resultLeftover)
	}
}

// --- UnmarshalTimestamp ---

func TestUnmarshalTimestamp(t *testing.T) {
	// Given
	id := 3
	timestamp := time.Now().Unix()
	state := "sTaRtEd"
	stringified := fmt.Sprintf("{\"Id\":  %d, \"STATE\":\"%s\",\"TimeSTaMp\": %d}", id, state, timestamp)

	expected := LogTimestamp {
		Id: id,
		Timestamp: timestamp,
		State: StartFlag,
	}

	bytes := []byte(stringified)

	// When
	result, err := unmarshalTimestamp(bytes)

	// Then
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestUnmarshalTimestampError(t *testing.T) {
	// Given
	stringified := "{\"d\":  3, \"StE\":\"STARTED\", \"tiMesTamP\": 100}"
	bytes := []byte(stringified)

	// When
	result, err := unmarshalTimestamp(bytes)

	// Then
	if result != (LogTimestamp{}) || err == nil {
		t.Errorf("Expected empty LogTimestamp and non-nil error, got %v and %v", result, err)
	}
}
