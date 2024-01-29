all: clean test build

test:
	go test

format:
	go fmt ./...

build:
	go build next-task.go

clean:
	git clean -fdX

install:
	go install .
