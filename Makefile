build:
	go build -v *.go
run: build
	./main image.ppm
