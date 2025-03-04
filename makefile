.DEFAULT_GOAL := all

all: build

build: *.go
	@echo "Build program"
	go build -ldflags="-w -s" -gcflags=all="-l -B" -o a.out

install:
	go install .

clean:
	@echo "Cleaning binaries"
	rm a.out

run:
	@echo "Run dummy"
	go run .

