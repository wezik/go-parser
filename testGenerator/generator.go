package testGenerator

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Log struct {
	id int
	state string
	timestamp int64
}

var mutex = &sync.Mutex{}

func Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide a file path")
		return
	}

	file, err := createFile(args[1])
	if err != nil {
		fmt.Println("Error creating file")
		return
	}
	writeToFile(1000, file)
}

func createFile(path string) (*os.File, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var _, err = os.Stat(path)
	var file *os.File

	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			fmt.Println("Error creating file")
			return file, err
		}

	} else {
		fmt.Println("File already exists")
		return file, err
	}

	fmt.Println("File created successfully")
	return file, nil
}

func writeToFile(logCount int, file *os.File) {
	var maximumDelayMs int64 = 10000
	stateStart := "STARTED"
	stateFinish := "FINISHED"

	logChan := make(chan Log)


	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		var collectedLogs []Log = make([]Log, 0)

		shuffleLogs := func() {
			rand.Shuffle(len(collectedLogs), func(i, j int) {
				collectedLogs[i], collectedLogs[j] = collectedLogs[j], collectedLogs[i]
			})
		}

		writeLogs := func() {
			mutex.Lock()
			defer mutex.Unlock()
			for _, log := range collectedLogs {
				writeLog(log, file)
			}
		}

		for log := range logChan {
			collectedLogs = append(collectedLogs, log)
			if len(collectedLogs) > 1024 {
				shuffleLogs()
				writeLogs()
				collectedLogs = collectedLogs[:0]
			}
		}
		
		shuffleLogs()
		writeLogs()
		collectedLogs = collectedLogs[:0]
	}()


	for i := 0; i < logCount; i++ {
		generatedDelay := rand.Int63n(maximumDelayMs)
		generatedOffset := rand.Int63n(maximumDelayMs)
		startLog := Log {
			id: i,
			state: stateStart,
			timestamp: time.Now().Unix() + generatedOffset,
		}
		logChan <- startLog
		endLog := Log {
			id: startLog.id,
			state: stateFinish,
			timestamp: startLog.timestamp + generatedDelay,
		}
		logChan <- endLog
	}
	close(logChan)
	wg.Wait()
}

func writeLog(log Log, file *os.File) {
	logString := fmt.Sprintf("%d %s %d\n", log.id, log.state, log.timestamp)
	_, err := file.WriteString(logString)
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}
}
