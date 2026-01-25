
process:
	@rm -rf build/ && mkdir -p build/
	@mdprc --source src/ --target build/

clean:
	@docker builder prune

build: process
	@docker build \
		--no-cache \
		--tag opencode \
		--file Dockerfile \
		.

up:
	@modd
