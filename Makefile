.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o bin/gwc main.go
