package ppm

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

	"github.com/Vikuuu/go-image-viewer/imageutil"
)

func ParsePPM(f *os.File, win *fyne.Window) {
	reader := bufio.NewReader(f)

	// read first line (P3 or P6)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var winW, winH int
	imgCanvas := &canvas.Image{}
	switch line {
	case "P3":
		w, h, img := ParseP3(reader)
		imgCanvas = canvas.NewImageFromImage(img)
		winW, winH = w, h
	case "P6":
		w, h, img := ParseP6(reader)
		imgCanvas = canvas.NewImageFromImage(img)
		winW, winH = w, h
	}

	(*win).SetContent(imgCanvas)
	(*win).Resize(fyne.NewSize(float32(winW), float32(winH)))
}

func ParseP3(reader *bufio.Reader) (w, h int, img *image.RGBA) {
	// read second line
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the Dimensions line of PPM p3 file: %v\n", err)
	}

	dimensions := strings.Split(line, " ")
	tempW, err := strconv.Atoi(strings.TrimSpace(dimensions[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	tempH, err := strconv.Atoi(strings.TrimSpace(dimensions[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// read third line (255)
	maxPixelValStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	maxPixelValStr = strings.TrimSpace(maxPixelValStr)
	maxPixelVal, err := strconv.Atoi(maxPixelValStr)

	factor := uint8(255 / maxPixelVal)

	w = int(tempW)
	h = int(tempH)
	img = image.NewRGBA(image.Rect(0, 0, w, h))

	input, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	colorScanner := imageutil.NewScanner(input)
	for i := range h {
		for j := range w {
			colors := make([]uint8, 0, 3)
			for len(colors) < 3 {
				c := colorScanner.NextNumber()
				colors = append(colors, c)
			}
			img.SetRGBA(j, i, color.RGBA{
				R: colors[0] * factor,
				G: colors[1] * factor,
				B: colors[2] * factor,
				A: 255,
			})
		}
	}

	return w, h, img
}

func ParseP6(reader *bufio.Reader) (w, h int, img *image.RGBA) {
	// read second line
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the Dimensions line of PPM p6 file: %v\n", err)
	}

	dimensions := strings.Split(line, " ")
	tempW, err := strconv.Atoi(strings.TrimSpace(dimensions[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	tempH, err := strconv.Atoi(strings.TrimSpace(dimensions[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// read third line (255)
	maxPixelValStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	maxPixelValStr = strings.TrimSpace(maxPixelValStr)
	maxPixelVal, err := strconv.Atoi(maxPixelValStr)

	factor := uint8(255 / maxPixelVal)

	w = int(tempW)
	h = int(tempH)
	img = image.NewRGBA(image.Rect(0, 0, w, h))
	for i := range h {
		for y := range w {
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
			img.SetRGBA(y, i, color.RGBA{
				R: colors[0] * factor,
				G: colors[1] * factor,
				B: colors[2] * factor,
				A: 255,
			})
		}
	}

	return w, h, img
}
