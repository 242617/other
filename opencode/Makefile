
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

restore: process
	@rm -rf ~/.config/opencode/ && mkdir -p ~/.config/opencode/
	@cp -R build/ ~/.config/opencode/

up:
	@modd
