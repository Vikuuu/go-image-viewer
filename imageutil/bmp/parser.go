package bmp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func ParseBMP(f *os.File, win *fyne.Window) {
	reader := bufio.NewReader(f)

	// Reading file header part
	fh := make([]byte, 14)
	n, err := reader.Read(fh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	if n != 14 {
		fmt.Fprintln(os.Stderr, "header not 14 byte, corrupted file")
	}
	// parsing file header
	bfType := string(fh[:2])
	if bfType != "BM" {
		fmt.Fprintln(os.Stderr, "file header mismatch, corrupt file")
	}
	// File size
	_ = string(fh[2:6])
	// Reserved bit (can ignore them or assert zero)
	_ = string(fh[6:8])
	_ = string(fh[8:10])
	// Offset with respect to file header pixel data start
	bfOffBit := binary.LittleEndian.Uint32(fh[10:14])

	// Header size(must be at least 40 bytes)
	biSize := make([]byte, 4)
	_, err = reader.Read(biSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	hSize := binary.LittleEndian.Uint32(biSize)
	if hSize < 40 {
		fmt.Fprintln(os.Stderr, "Header size must at least be 40 bytes")
	}

	// Reading image header part
	ih := make([]byte, hSize)
	n, err = reader.Read(ih)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	if n != 40 {
		fmt.Fprintln(os.Stderr, "header not 14 byte, corrupted file")
	}

	// width of image
	bWidth := binary.LittleEndian.Uint32(ih[0:4])
	fmt.Println(bWidth)
	// height of image
	bHeight := binary.LittleEndian.Uint32(ih[4:8])
	fmt.Println(bHeight)
	// biPlanes must be 1
	biPlanes := binary.LittleEndian.Uint16(ih[8:10])
	if biPlanes != 1 {
		fmt.Fprintln(os.Stderr, "BiPlanes must be equals to 1")
	}
	// Bits per pixel- 1, 4, 8, 16, 24, or 32
	// accept on 24 bits
	biBitCount := binary.LittleEndian.Uint16(ih[10:12])
	if biBitCount != 24 {
		fmt.Fprintln(os.Stderr, "Only bit 24 accepted")
	}
	// Compression type (0 = uncompressed)
	// biCompression := binary.LittleEndian.Uint32(ih[12:16])
	// // Image size, may be zero for uncompressed image
	// biSizeImage := binary.LittleEndian.Uint32(ih[16:20])
	// // Preferred resolution in pixels per meter
	// biXPelsPerMeter := binary.LittleEndian.Uint32(ih[20:24])
	// // Preferred resolution in pixels per meter
	// biYPelsPerMeter := binary.LittleEndian.Uint32(ih[24:28])
	// // Number color map entries that are actually used
	// biClrUsed := binary.LittleEndian.Uint32(ih[28:32])
	// // Number of significant colors
	// biClrImp := binary.LittleEndian.Uint32(ih[32:36])
	// input, err := io.ReadAll(reader)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// }
	// colorScanner := imageutil.NewScanner(input)

	_, err = f.Seek(int64(bfOffBit), 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	reader = bufio.NewReader(f)

	img := image.NewRGBA(image.Rect(0, 0, int(bWidth), int(bHeight)))

	r := (bWidth * 3)
	padding := (4 - (r % 4)) % 4
	w := make([]byte, r+padding)
	for row := 0; row < int(bHeight); row++ {
		_, err := reader.Read(w)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		y := int(bHeight) - row - 1
		x := 0
		for i := 0; i < int(r); i += 3 {
			img.SetRGBA(x, y, color.RGBA{
				R: w[i+2],
				G: w[i+1],
				B: w[i+0],
				A: 255,
			})
			x++
		}
	}

	imgCanvas := canvas.NewImageFromImage(img)

	(*win).SetContent(imgCanvas)
	(*win).Resize(fyne.NewSize(float32(bWidth), float32(bHeight)))
}
