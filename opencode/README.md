# OpenCode Sandbox

A Docker-based development environment for running and developing with OpenCode in an isolated container environment.

## Overview

This sandbox provides a complete OpenCode setup with specialized agents configured for collaborative software development:

- **Orchestrator Agent** - Coordinates task execution across specialized agents
- **Project Analyst** - Creates and updates Product Requirements Documents (PRDs)
- **Go Developer** - Writes production-ready Go code following best practices
- **Test Engineer** - Develops comprehensive unit and integration tests
- **Code Reviewer** - Ensures code quality, security, and performance
- **Technical Writer** - Creates and maintains technical documentation

## Agent Workflow

- **Project Analyst** составляет Product Requirements Document или (обновляет его) по существующему проекту и составляет план на доработку (добавление фичи).
- **Orchestrator/Coordinator Agent** делегирует выполнение подзадач разным агентам:
  - **Test Engineer Agent** подготавливает Unit- и интеграционные тесты,
  - **Go Developer Agent** пишет код,
  - **Code Reviewer Agent** проверяет корректность кода и полноту выполнения задания,
  - **Technical Writer Agent** составляет/актуализирует документацию (важно это делать каждый раз после внесения изменений)
- **Orchestrator/Coordinator Agent** проверяет выполнение всех поставленных задач.

## Quick Start

### 1. Build the environment

```bash
make build
```

### 2. Set up shell function for easy access

For **fish shell** (as provided):
```fish
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

For **bash/zsh**:
```bash
oc() {
    docker run --rm -it \
        --cpus 2 \
        --memory 4g \
        --network host \
        -e DEEPSEEK_API_KEY=$DEEPSEEK_API_KEY \
        -e PERPLEXITY_API_KEY=$PERPLEXITY_API_KEY \
        -v $(pwd):/workspace \
        opencode
}
```

### 3. Run OpenCode in sandbox

```bash
oc
```

## Configuration

- **Model**: Qwen3 30B Coder (primary), with DeepSeek fallback
- **Tools**: Full access to Context7 for library documentation
- **Permissions**: Controlled access based on agent roles
- **Environment**: Isolated Docker container with mounted workspace

## Available Commands

- `make clean` - Clean Docker builder cache
- `make build` - Build the OpenCode Docker image
- `make restore` - Sync local OpenCode agent configurations

## Environment Variables

- `DEEPSEEK_API_KEY` - API key for DeepSeek models
- `PERPLEXITY_API_KEY` - API key for Perplexity search

## TODO

- [ ] Configure universal permissions for all agents
- [ ] Implement session persistence (*.local/share/opencode/storage/session*)
- [ ] Add support for additional model providers
- [ ] Improve error handling and recovery mechanisms
