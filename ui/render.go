package ui

import (
	"com/parser/ui/components"
	"fmt"
)

func ClearScreen(rows int) {
	fmt.Printf("\r")               // Move cursor to the begining of the line
	if rows > 1 {
		fmt.Printf("\033[%dA", rows - 1) // Move cursor up the height of the component
	}
	fmt.Printf("\033[J")           // Clear from cursor to the end
}

func Render(c components.Component) {
	fmt.Print(c.String())
}

func RenderOnTop(c components.Component, rows int) {
	ClearScreen(rows)
	Render(c)
}

func ReRender(c components.Component) {
	RenderOnTop(c, components.Height(c))
}
