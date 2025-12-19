package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"

	"github.com/Vikuuu/go-image-viewer/imageutil/bmp"
	"github.com/Vikuuu/go-image-viewer/imageutil/ppm"
)

func parseImage(fp string) {
	fileName := filepath.Base(fp)
	ext := filepath.Ext(fileName)

	myApp := app.New()
	win := myApp.NewWindow(fileName)
	file, err := os.OpenFile(fp, os.O_RDONLY, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	defer file.Close()
	switch ext {
	case ".ppm":
		ppm.ParsePPM(file, &win)
	case ".bmp":
		bmp.ParseBMP(file, &win)
	}
	win.ShowAndRun()
}
