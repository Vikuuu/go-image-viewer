build:
	go build -o ./bin/image_viewer -v *.go 
run: build
	./bin/image_viewer ./testdata/p3image.ppm
