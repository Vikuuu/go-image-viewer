package png

import (
	"bufio"
	"fmt"
	"image"
	"os"
)

func DecodePNG(reader *bufio.Reader) (int, int, *image.RGBA) {
	verifyPngSig(reader)

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
