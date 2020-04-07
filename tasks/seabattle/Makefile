
.PHONY: build
build:
	go build \
		-o seabattle \
		cmd/seabattle/main.go

.PHONY: test
test:
	go test --race ./...

.PHONY: run
run: build test
	./seabattle
