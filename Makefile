build:
	- go build -o vigo360 -ldflags "-X main.version=$(shell git rev-parse --short HEAD)" .

run:
	- go run -ldflags "-X main.version=$(shell git rev-parse --short HEAD)" .