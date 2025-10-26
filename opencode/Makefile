
clean:
	@docker builder prune

build:
	@docker build \
		--no-cache \
		--tag opencode \
		--file Dockerfile \
		.

restore:
	@go install ../mdprc/cmd/main.go && mv ~/go/bin/main ~/go/bin/mdprc
	@mdprc --source src/ --target ~/.config/opencode/
