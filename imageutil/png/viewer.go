package png

import (
	"bufio"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func ViewPNGImage(f *os.File, win *fyne.Window) {
	reader := bufio.NewReader(f)
	imgCanvas := &canvas.Image{}

	winW, winH, img := DecodePNG(reader)
	imgCanvas = canvas.NewImageFromImage(img)

	(*win).SetContent(imgCanvas)
	(*win).Resize(fyne.NewSize(float32(winW), float32(winH)))
}
