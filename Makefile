build:
	- go build -o vigo360 -ldflags "-X main.version=$(shell git rev-parse --short HEAD)" .

run:
	- @export $(shell cat .env | grep -v '^#' | xargs)
	- go run -ldflags "-X main.version=$(shell git rev-parse --short HEAD)" .