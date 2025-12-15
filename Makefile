build:
	go build -v *.go
run: build
	./main p3image.ppm
