package menu

import (
	"bufio"
	"com/parser/database"
	"com/parser/generator"
	"com/parser/parser"
	"com/parser/ui"
	"com/parser/ui/components"
	"com/parser/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	menuPrint = "---Log parser---|\n" +
		"1. Parse File   |\n" +
		"2. Generate Test|\n" +
		"3. DB test add  |\n" +
		"4. Db test fetch|\n" +
		"Q. Quit         |\n" +
		"-----------------\n"

	askInput    = "Input: "
	askLogCount = "Amount of logs to write: "
	askFileName = "File name: "

	errorFileCreation = "Error creating the file:"
	errorLogCount     = "Error parsing log count: Use a number"
)

var scanner *bufio.Scanner

func init() {
	scanner = bufio.NewScanner(os.Stdin)
}

func Start() {
	var menuText components.SimpleText
	menuText.SetText(menuPrint)
	for {
		ui.Render(menuText)
		switch strings.ToUpper(readUserInput(askInput)) {
		case "1":
			fileName := readUserInput(askFileName)
			err := parser.Parse(fileName)
			if err != nil {
				fmt.Println(err)
			}
		case "2":
			count, file, err := readGenerateInputs()
			if err != nil {
				continue
			}
			setupGenerate(count, file)
		case "3":
			log := parser.Log{
				Id:              99,
				TimestampStart:  time.Now(),
				TimestampFinish: time.Now(),
			}
			database.InsertLog(log)
		case "4":
			logs, err := database.FetchLogs()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Logs:", len(logs))
			for _, log := range logs {
				fmt.Println(log)
			}
		case "Q", "QUIT":
			return
		}
	}
}

func readUserInput(label string) string {
	var input components.Input
	input.SetLabel(label)
	ui.Render(input)

	scanner.Scan()
	return strings.Trim(scanner.Text(), " ")
}

func readGenerateInputs() (int, string, error) {
	count, err := strconv.Atoi(readUserInput(askLogCount))
	if err != nil {
		fmt.Println(errorLogCount)
		return 0, "", err
	}
	fileName := readUserInput(askFileName)
	return count, fileName, nil
}

func setupGenerate(count int, fileName string) {
	file, err := utils.CreateFile(fileName)
	if err != nil {
		fmt.Println(errorFileCreation, err)
		return
	}
	defer file.Close()

	generator.GenerateToFile(file, count)
}
