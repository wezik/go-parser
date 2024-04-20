package components

import (
	"com/parser/utils"
	"fmt"
	"strings"
)

type ProgressBar struct {
	percentage float32
	prefix string
	suffix string
	width int
}

func (pb *ProgressBar) Percentage() float32 {
	return pb.percentage
}

func (pb *ProgressBar) Prefix() string {
	return pb.prefix
}

func (pb *ProgressBar) Suffix() string {
	return pb.suffix
}

func (pb *ProgressBar) Width() int {
	return pb.width
}

func (pb *ProgressBar) SetPercentage(percentage float32) {
	pb.percentage = percentage
}

func (pb *ProgressBar) SetPrefix(prefix string) {
	pb.prefix = prefix
}

func (pb *ProgressBar) SetSuffix(suffix string) {
	pb.suffix = suffix
}

func (pb *ProgressBar) SetWidth(width int) {
	pb.width = width
}

const (
	fullBlock = "█"
	// boldBlock = "▓"
	// regularBlock = "▒"
	thinBlock = "░"
)

func (pb ProgressBar) String() string {
	var b strings.Builder

	prefixString := pb.Prefix() + " "
	percentageString := fmt.Sprintf(" %.2f%% ", pb.Percentage())
	barWidth := pb.Width()
	suffixString := pb.Suffix()

	expectedLength := len(prefixString) + barWidth + len(percentageString) + len(suffixString)
	terminalWidth := utils.TerminalWidth()

	if expectedLength > terminalWidth {
		lowerBarWidth := terminalWidth - expectedLength + barWidth
		if lowerBarWidth < 10 {
			barWidth = 10
			prefixString = prefixString + "\n"
		}
	}

	filledBlocks := float32(pb.Width()) * pb.Percentage() / 100
	b.WriteString(prefixString)

	for i := 0; i < barWidth; i++ {
		if i < int(filledBlocks) {
			b.WriteString(fullBlock)
		} else {
			b.WriteString(thinBlock)
		}
	}

	b.WriteString(percentageString)
	b.WriteString(pb.Suffix())
	return b.String()
}

func ProgressBarDefault() *ProgressBar {
	return &ProgressBar {
		percentage: 0.0,
		prefix: "Progress: ",
		suffix: "",
		width: 20,
	}
}
