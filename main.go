package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2/app"
)

type state int

const (
	skip state = iota
	integer
)

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Image Viewer in Go")

	fileName := os.Args[1]

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	defer file.Close()

	parseImage(file, &win)
	win.ShowAndRun()
}
