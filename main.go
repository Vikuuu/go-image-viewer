package main

import (
	"bufio"
	"errors"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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

	reader := bufio.NewReader(file)

	// read first line (P3 or P6)
	line, err := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var w, h int64
	pixel := []color.RGBA{}
	switch line {
	case "P3":
		w, h, pixel = parseP3(reader)
	case "P6":
		w, h, pixel = parseP6(reader)
	}

	imgW, imgH := int(w), int(h)
	raster := canvas.NewRasterWithPixels(func(x, y, _, _ int) color.Color {
		if x >= imgW || y >= imgH {
			return color.Black
		}
		return pixel[y*imgW+x]
	})

	fmt.Println(w, h)

	win.SetContent(raster)
	win.Resize(fyne.NewSize(float32(w), float32(h)))
	win.ShowAndRun()
}

func parseP3(reader *bufio.Reader) (w, h int64, pixel []color.RGBA) {
	fmt.Println("Calling P3 parse")
	// read second line
	line, _ := reader.ReadString('\n')
	dimensions := strings.Split(line, " ")
	w, _ = strconv.ParseInt(dimensions[0], 10, 32)
	h, _ = strconv.ParseInt(strings.TrimSpace(dimensions[1]), 10, 32)
	// read third line (255)
	line, _ = reader.ReadString('\n')

	pixel = make([]color.RGBA, w*h)
	idx := 0

	input, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	colorScanner := NewScanner(input)
	for idx < len(pixel) {
		colors := make([]uint8, 0, 3)
		for len(colors) < 3 {
			c := colorScanner.NextNumber()
			colors = append(colors, c)
		}
		if len(colors) < 3 {
			break
		}
		pixel[idx] = color.RGBA{
			R: colors[0],
			G: colors[1],
			B: colors[2],
			A: 255,
		}
		idx++
	}

	return w, h, pixel
}

func parseP6(reader *bufio.Reader) (w, h int64, pixel []color.RGBA) {
	// read second line
	fmt.Println("Calling P6 parse")
	line, _ := reader.ReadString('\n')
	dimensions := strings.Split(line, " ")
	w, _ = strconv.ParseInt(dimensions[0], 10, 32)
	h, _ = strconv.ParseInt(strings.TrimSpace(dimensions[1]), 10, 32)
	// read third line (255)
	line, _ = reader.ReadString('\n')

	pixel = make([]color.RGBA, w*h)
	idx := 0

	for idx < len(pixel) {
		colors := make([]byte, 0, 3)
		for len(colors) < 3 {
			c, err := reader.ReadByte()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			colors = append(colors, c)
		}

		if len(colors) < 3 {
			break
		}

		pixel[idx] = color.RGBA{R: colors[0], G: colors[1], B: colors[2], A: 255}
		idx++
	}

	return w, h, pixel
}
