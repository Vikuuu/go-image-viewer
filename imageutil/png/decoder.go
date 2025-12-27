package png

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
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

	decompressedData := parseIDAT(pf)
	fmt.Println(len(decompressedData))
	fmt.Println(len(pf.idat))

	return pf.w, pf.h, &image.RGBA{}
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

func parseIDAT(pf *pngFile) []byte {
	r := bytes.NewReader(pf.idat)
	rc, err := zlib.NewReader(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, rc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return b.Bytes()
}
