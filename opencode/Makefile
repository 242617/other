
clean:
	@docker builder prune

build:
	@docker build \
		--tag opencode \
		--file Dockerfile \
		.

restore:
	@rsync -var opencode/* ~/.config/opencode/
