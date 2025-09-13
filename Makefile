
test:
	@go test \
		--race \
		-v \
		./...

install\:rename:
	@go install rename/rename.go
	# @go install github.com/242617/other/rename/rename.go
