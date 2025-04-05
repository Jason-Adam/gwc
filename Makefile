.PHONY: run
run:
	go run cmd/gwc/main.go

.PHONY: build
build:
	go build -o bin/gwc cmd/gwc/main.go
