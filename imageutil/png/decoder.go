package png

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
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
			_, err = io.ReadFull(r, curChunk.crc[:])
			break OUTER
		case "IHDR":
			parseIHDR(ihdr, curChunk.data, curChunk.crc)
		case "IDAT":
			pf.idat = append(pf.idat, curChunk.data...)
		}
	}

	inflateData := inflateIDAT(pf)

	var bytesPerPixel int
	switch pf.colorType {
	case 6:
		bytesPerPixel = 4
	default:
		fmt.Fprintln(os.Stderr, "Not implemented color type: ", pf.colorType)
	}

	bytesPerRow := (pf.w * bytesPerPixel) + 1
	rowData := make([]byte, pf.h*(pf.w*4))

	for {
		sc := make([]byte, bytesPerRow)
		_, err := io.ReadFull(inflateData, sc)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Fprintln(os.Stderr, err)
		}
		switch sc[0] {
		case 0:
			rowData = append(rowData, sc[1:]...)
		case 1:
			// SUB filter method
			for i := range sc[1:] {
				subX := int(sc[i])
				var rawXMinBpp int
				if i-bytesPerPixel < 0 {
					rawXMinBpp = 0
				} else {
					rawXMinBpp = int(sc[i-bytesPerPixel])
				}
				pixelData := subX + rawXMinBpp
				rowData = append(rowData, byte(pixelData))
			}
		case 2:
		case 3:
		case 4:
			fallthrough
		default:
			fmt.Fprintln(os.Stderr, "Not implemented")
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, pf.w, pf.h))
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
	if rc == nil {
		fmt.Fprintln(os.Stderr, "rc is nil")
	}

	defer rc.Close()
	var b bytes.Buffer
	_, err = io.Copy(&b, rc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if b.Len() != (pf.h * (pf.w*4 + 1)) {
		fmt.Fprintln(os.Stderr, "inflated len not equal")
	}

	return &b
}
