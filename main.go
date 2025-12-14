package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	fileName := os.Args[1]

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// first scan (file data type P3 or P6)
	scanner.Scan()
	_ = scanner.Text()
	// second line for dimensions
	scanner.Scan()
	dimensions := strings.Split(scanner.Text(), " ")
	w, err := strconv.ParseFloat(dimensions[0], 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	h, err := strconv.ParseFloat(dimensions[1], 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}

	a := app.New()
	win := a.NewWindow("Image viewer")

	// rect := canvas.NewRectangle(color.White)
	// win.SetContent(rect)

	fmt.Println("Size: ", w, " ", h)
	win.Resize(fyne.NewSize(float32(w), float32(h)))
	win.ShowAndRun()
}
