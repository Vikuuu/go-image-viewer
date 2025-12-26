package png

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"os"
)

type pngChunk struct {
	dataLenByte [4]byte
	chunkType   [4]byte
	data        []byte
	crc         [4]byte
}

func DecodePNG(r *bufio.Reader) (int, int, *image.RGBA) {
	verifyPngSig(r)

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

		if string(curChunk.chunkType[:]) == "IEND" {
			_, err = io.ReadFull(r, curChunk.crc[:])
			fmt.Println(chunkLen)
			fmt.Println(curChunk.chunkType[:])
			fmt.Println(curChunk.data)
			fmt.Println(curChunk.crc[:])
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

		fmt.Println(chunkLen)
		fmt.Println(curChunk.chunkType[:])
		fmt.Println(curChunk.data)
		fmt.Println(curChunk.crc[:])
	}

	return 0, 0, &image.RGBA{}
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
