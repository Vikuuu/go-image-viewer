package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func parseImage(f *os.File, win *fyne.Window) {
	reader := bufio.NewReader(f)

	// read first line (P3 or P6)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var winW, winH int
	imgCanvas := &canvas.Image{}
	switch line {
	case "P3":
		w, h, pixel := parseP3(reader)
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := range h {
			for x := range w {
				img.SetRGBA(x, y, pixel[y][x])
			}
		}
		imgCanvas = canvas.NewImageFromImage(img)
		winW, winH = w, h
	case "P6":
		w, h, pixel := parseP6(reader)
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := range h {
			for x := range w {
				img.SetRGBA(x, y, pixel[y][x])
			}
		}
		imgCanvas = canvas.NewImageFromImage(img)
		winW, winH = w, h
	}

	(*win).SetContent(imgCanvas)
	(*win).Resize(fyne.NewSize(float32(winW), float32(winH)))
}

func parseP3(reader *bufio.Reader) (w, h int, pixels [][]color.RGBA) {
	// read second line
	line, _ := reader.ReadString('\n')
	dimensions := strings.Split(line, " ")
	tempW, _ := strconv.ParseInt(dimensions[0], 10, 32)
	tempH, _ := strconv.ParseInt(strings.TrimSpace(dimensions[1]), 10, 32)
	// read third line (255)
	line, _ = reader.ReadString('\n')

	w = int(tempW)
	h = int(tempH)
	pixels = make([][]color.RGBA, h)

	input, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	colorScanner := NewScanner(input)
	for i := range len(pixels) {
		pixels[i] = make([]color.RGBA, w)
		for j := range len(pixels[i]) {
			colors := make([]uint8, 0, 3)
			for len(colors) < 3 {
				c := colorScanner.NextNumber()
				colors = append(colors, c)
			}
			if len(colors) < 3 {
				break
			}
			pixels[i][j] = color.RGBA{
				R: colors[0],
				G: colors[1],
				B: colors[2],
				A: 255,
			}
		}
	}

	return w, h, pixels
}

func parseP6(reader *bufio.Reader) (w, h int, pixels [][]color.RGBA) {
	// read second line
	line, _ := reader.ReadString('\n')
	dimensions := strings.Split(line, " ")
	tempW, _ := strconv.ParseInt(dimensions[0], 10, 32)
	tempH, _ := strconv.ParseInt(strings.TrimSpace(dimensions[1]), 10, 32)
	// read third line (255)
	line, _ = reader.ReadString('\n')

	w = int(tempW)
	h = int(tempH)
	pixels = make([][]color.RGBA, h)
	for i := range len(pixels) {
		pixels[i] = make([]color.RGBA, w)
		for y := range len(pixels[i]) {
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
			pixels[i][y] = color.RGBA{
				R: colors[0],
				G: colors[1],
				B: colors[2],
				A: 255,
			}
		}

	}

	return w, h, pixels
}
