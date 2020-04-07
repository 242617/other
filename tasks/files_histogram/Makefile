
.PHONY: build
build:
	go build \
		-o main \
		cmd/main/main.go

.PHONY: test
test:
	go test --race ./...

.PHONY: run
run: build test
	./main
