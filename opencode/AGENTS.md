# OpenCode Sandbox Agent Guidelines

## Environment Overview
This is a Docker-based sandbox environment for running OpenCode with specialized agents. The configuration includes:
- **Orchestrator**: Coordinates task execution across specialized agents
- **Project Analyst**: Creates/updates PRDs and development plans
- **Go Developer**: Writes production-ready Go code following best practices
- **Test Engineer**: Develops comprehensive unit and integration tests
- **Code Reviewer**: Ensures code quality, security, and performance
- **Technical Writer**: Creates and maintains documentation

## Build Commands
- **Build**: `make build` - Builds OpenCode Docker image
- **Clean**: `make clean` - Cleans Docker builder cache
- **Restore**: `make restore` - Syncs agent configs to local directory

## Agent Workflow

- **Project Analyst** составляет Product Requirements Document или (обновляет его) по существующему проекту и составляет план на доработку (добавление фичи).
- **Orchestrator/Coordinator Agent** делегирует выполнение подзадач разным агентам:
  - **Test Engineer Agent** подготавливает Unit- и интеграционные тесты,
  - **Go Developer Agent** пишет код,
  - **Code Reviewer Agent** проверяет корректность кода и полноту выполнения задания,
  - **Technical Writer Agent** составляет/актуализирует документацию (важно это делать каждый раз после внесения изменений)
- **Orchestrator/Coordinator Agent** проверяет выполнение всех поставленных задач.

## Development Guidelines
- All agents work based on PRD requirements
- Follow Go best practices for code written by Go Developer
- Use Context7 for library documentation and API references
- Maintain comprehensive test coverage (>80% for critical paths)
- Ensure proper error handling and input validation

## Running the Sandbox
```bash
make build
oc  # Use the provided shell function to run OpenCode
```
