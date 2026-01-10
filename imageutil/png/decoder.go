package png

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
)

type ihdr struct {
	w, h              int
	bitDepth          int
	colorType         int
	compressionMethod int
	filterMethod      int
	interlaceMethod   int
}

type pngFile struct {
	*ihdr
	idat []byte
}

type pngChunk struct {
	dataLenByte [4]byte
	chunkType   [4]byte
	data        []byte
	crc         [4]byte
}

func DecodePNG(r *bufio.Reader) (int, int, *image.RGBA) {
	ihdr := &ihdr{}
	pf := &pngFile{ihdr: ihdr}
	verifyPngSig(r)

OUTER:
	for {
		curChunk := &pngChunk{}

		// Read the length of the chunk
		_, err := io.ReadFull(r, curChunk.dataLenByte[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing PNG file: %v\n", err)
			break
		}
		chunkLen := int(binary.BigEndian.Uint32(curChunk.dataLenByte[:]))

		// Read the chunk type
		_, err = io.ReadFull(r, curChunk.chunkType[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing PNG file: %v\n", err)
			break
		}

		curChunk.data = make([]byte, chunkLen)

		// Read the chunk data
		n, err := io.ReadFull(r, curChunk.data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing PNG file: %v\n", err)
			break
		}
		if n < chunkLen {
			break
		}

		// Read the CRC data
		_, err = io.ReadFull(r, curChunk.crc[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing PNG file: %v\n", err)
			break
		}

		switch string(curChunk.chunkType[:]) {
		case "IEND":
			break OUTER
		case "IHDR":
			parseIHDR(ihdr, curChunk.data, curChunk.crc)
		case "IDAT":
			pf.idat = append(pf.idat, curChunk.data...)
		}
	}

	inflateData := inflateIDAT(pf)

	var bpp int
	switch pf.colorType {
	case 6:
		bpp = 4
	default:
		fmt.Fprintln(os.Stderr, "Not implemented color type: ", pf.colorType)
	}

	bytesPerRow := (pf.w * bpp) + 1
	rowData := make([][]byte, pf.h)

	row := 0
	for {
		sc := make([]byte, bytesPerRow)
		_, err := io.ReadFull(inflateData, sc)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Fprintln(os.Stderr, err)
		}
		rowData[row] = make([]byte, 0, (pf.w * 4))
		filterType := sc[0]
		rawPixelData := sc[1:]
		switch filterType {
		case 0:
			rowData[row] = append(rowData[row], rawPixelData...)
		case 1:
			// SUB filter method
			for i := 0; i < len(rawPixelData); i++ {
				subX := int(rawPixelData[i])
				var rawXMinBpp int
				if i-bpp < 0 {
					rawXMinBpp = 0
				} else {
					rawXMinBpp = int(rowData[row][i-bpp])
				}
				pixelData := subX + rawXMinBpp
				rowData[row] = append(rowData[row], byte(pixelData))
			}
		case 2:
			// Up filter method
			if row == 0 {
				rowData[row] = append(rowData[row], rawPixelData...)
			} else {
				for i := 0; i < len(rawPixelData); i++ {
					rawX := int(rawPixelData[i])
					priorX := int(rowData[row-1][i])
					pixelData := rawX + priorX
					rowData[row] = append(rowData[row], byte(pixelData))
				}
			}
		case 3:
			// Average filter method
			for i := 0; i < len(rawPixelData); i++ {
				avgX := int(rawPixelData[i])
				var rawXMinBpp int
				var priorX int
				if i-bpp < 0 {
					rawXMinBpp = 0
				} else {
					rawXMinBpp = int(rowData[row][i-bpp])
				}
				if row == 0 {
					priorX = 0
				} else {
					priorX = int(rowData[row-1][i])
				}
				pixelData := avgX + ((rawXMinBpp + priorX) / 2)
				rowData[row] = append(rowData[row], byte(pixelData))
			}
		case 4:
			// Paeth filter method
			for i := 0; i < len(rawPixelData); i++ {
				rawX := int(rawPixelData[i])
				var rawXMinBpp, priorX, priorXMinBpp int

				if i-bpp < 0 {
					rawXMinBpp = 0
				} else {
					rawXMinBpp = int(rowData[row][i-bpp])
				}
				if row == 0 {
					priorX = 0
				} else {
					priorX = int(rowData[row-1][i])
				}
				if row == 0 || i-bpp < 0 {
					priorXMinBpp = 0
				} else {
					priorXMinBpp = int(rowData[row-1][i-bpp])
				}

				pixelData := rawX + paethPredictor(rawXMinBpp, priorX, priorXMinBpp)

				rowData[row] = append(rowData[row], byte(pixelData))
			}
		default:
			fmt.Fprintln(os.Stderr, "Not implemented > 4")
		}
		row++
	}

	img := image.NewRGBA(image.Rect(0, 0, pf.w, pf.h))
	for h := range rowData {
		for w := 0; w < pf.w; w++ {
			byteIndex := w * 4
			img.SetRGBA(
				w,
				h,
				color.RGBA{
					R: rowData[h][byteIndex],
					G: rowData[h][byteIndex+1],
					B: rowData[h][byteIndex+2],
					A: rowData[h][byteIndex+3],
				},
			)
		}
	}
	return pf.w, pf.h, img
}

func verifyPngSig(reader *bufio.Reader) {
	pngFileSig := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	readPngFileSig := make([]byte, 8)

	n, err := reader.Read(readPngFileSig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	if n != 8 {
		fmt.Fprintln(os.Stderr, "png sig error")
	}

	for i := range len(pngFileSig) {
		if pngFileSig[i] != readPngFileSig[i] {
			panic("Malformed png file signature")
		}
	}
}

func parseIHDR(header *ihdr, data []byte, crc [4]byte) {
	w, h := data[0:4], data[4:8]

	header.w = int(binary.BigEndian.Uint32(w))
	header.h = int(binary.BigEndian.Uint32(h))

	header.bitDepth = int(data[8])
	header.colorType = int(data[9])
	header.compressionMethod = int(data[10])
	header.filterMethod = int(data[11])
	header.interlaceMethod = int(data[12])

	fmt.Printf("%v\n", header)
}

func inflateIDAT(pf *pngFile) io.Reader {
	r := bytes.NewReader(pf.idat)
	rc, err := zlib.NewReader(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer rc.Close()

	var b bytes.Buffer
	_, err = io.Copy(&b, rc)

	return bytes.NewReader(b.Bytes())
}
