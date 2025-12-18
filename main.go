package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"
)

type state int

const (
	skip state = iota
	integer
)

func main() {
	filePath := os.Args[1]
	fileName := filepath.Base(filePath)

	myApp := app.New()
	win := myApp.NewWindow(fileName)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	defer file.Close()

	parseImage(file, &win)
	win.ShowAndRun()
}
