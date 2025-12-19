package main

import (
	"os"
)

func main() {
	filePath := os.Args[1]

	parseImage(filePath)
}
