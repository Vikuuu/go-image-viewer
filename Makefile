build:
	go build -o ./bin/image_viewer -v *.go 
run: build
	# Testing PPM P3 image
	# ./bin/image_viewer ./testdata/p3image.ppm
	# Testing PPM P6 image
	# ./bin/image_viewer ./testdata/p6image.ppm
	# Testing BMP iamge
	./bin/image_viewer ./testdata/bmp-sample.bmp
