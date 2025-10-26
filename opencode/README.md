# OpenCode Sandbox

A Docker-based development environment for running and developing with OpenCode.

## Quick Start

1. Build the environment:

```sh
make build
```

2. Add function for running OpenCode with local directory mounted and */workspace/* (for `fish`):

```
function oc
	docker run --rm -it \
		--cpus 2 \
		--memory 4g \
		--network host \
		-e DEEPSEEK_API_KEY=$DEEPSEEK_API_KEY \
		-e PERPLEXITY_API_KEY=$PERPLEXITY_API_KEY \
		-v $(pwd):/workspace \
		opencode
end
```

3. Run OpenCode in sandbox:

```sh
oc
```

## TODO

- [ ] Give all permissions for everyone.
- [ ] Save sessions (*.local/share/opencode/storage/session*).
