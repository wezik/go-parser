package components

import (
	"com/parser/utils"
	"fmt"
	"strings"
)

type ProgressBar struct {
	value int64
	max int64
	prefix string
	suffix string
	width int
}

func (pb *ProgressBar) Max() int64 {
	return pb.max
}

func (pb *ProgressBar) Prefix() string {
	return pb.prefix
}

func (pb *ProgressBar) Suffix() string {
	return pb.suffix
}

func (pb *ProgressBar) Value() int64 {
	return pb.value
}

func (pb *ProgressBar) Width() int {
	return pb.width
}

func (pb *ProgressBar) SetMax(max int64) {
	pb.max = max
}

func (pb *ProgressBar) SetPrefix(prefix string) {
	pb.prefix = prefix
}

func (pb *ProgressBar) SetSuffix(suffix string) {
	pb.suffix = suffix
}

func (pb *ProgressBar) SetValue(value int64) {
	pb.value = value
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
	percentage := float32(pb.Value()) / float32(pb.Max()) * 100
	percentageString := fmt.Sprintf(" %.2f%% ", percentage)
	barWidth := pb.Width()
	suffixString := fmt.Sprintf(pb.Suffix(), pb.Value(), pb.Max())

	expectedLength := len(prefixString) + barWidth + len(percentageString) + len(suffixString)
	terminalWidth := utils.TerminalWidth()

	if expectedLength > terminalWidth {
		lowerBarWidth := terminalWidth - expectedLength + barWidth
		if lowerBarWidth < 10 {
			barWidth = 10
			prefixString = prefixString + "\n"
		}
	}

	filledBlocks := float32(pb.Width()) * percentage / 100
	b.WriteString(prefixString)

	for i := 0; i < barWidth; i++ {
		if i < int(filledBlocks) {
			b.WriteString(fullBlock)
		} else {
			b.WriteString(thinBlock)
		}
	}

	b.WriteString(percentageString)
	b.WriteString(suffixString)
	return b.String()
}

func ProgressBarDefault() ProgressBar {
	return ProgressBar {
		value: 0,
		max: 100,
		prefix: "Progress:",
		suffix: "[%d/%d]",
		width: 20,
	}
}
