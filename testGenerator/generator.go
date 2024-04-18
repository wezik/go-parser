package testGenerator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"
)

type Log struct {
	Id int
	state string
	timestamp int64
	delay *big.Int 
}

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
	maximumDelay := 10000
	stateStart := "STARTED"
	// stateFinish := "FINISHED"

	startedLogs := make([]Log, 0)

	for i := 0; i < logCount; i++ {
		maximumDelayInt := big.NewInt(int64(maximumDelay / 10))
		random := rand.Reader
		generatedDelay, err := rand.Int(random, maximumDelayInt)
		if err != nil {
			fmt.Println("Error generating random number")
			return
		}
		log := Log{
			Id: i,
			state: stateStart,
			timestamp: time.Now().Unix(),
			delay: generatedDelay,

		}
		startedLogs = append(startedLogs, log)
		writeLog(log, file)
	}
	
}

func writeLog(log Log, file *os.File) {
	logString := fmt.Sprintf("%d %s %d %d\n", log.Id, log.state, log.timestamp, log.delay)
	_, err := file.WriteString(logString)
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}
}
