
clean:
	@docker builder prune

build:
	@docker build \
		--no-cache \
		--tag opencode \
		--file Dockerfile \
		.

restore:
	@rsync -var opencode/* ~/.config/opencode/
