package components

import (
	"strings"
)

type Component interface {
	String() string
}

// SimpleText
type SimpleText struct {
	text string
}

func (c *SimpleText) SetText(label string) {
	c.text = label
}

func (c SimpleText) String() string {
	return c.text
}

// Input
type Input struct {
	label string
}

func (c *Input) SetLabel(label string) {
	c.label = label
}

func (c Input) String() string {
	return c.label
}

// Util function to get height
func Height(c Component) int {
	s := c.String() + "0" // Add any character to detect empty strings
	seperated := strings.Split(s, "\n")
	return len(seperated)
}
