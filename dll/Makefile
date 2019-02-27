dll:
	del lib.dll
	go build -o lib.dll -buildmode=c-shared lib.go

run:
	go build -o main.exe main.go
	./main.exe

test: dll run