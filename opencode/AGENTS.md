# OpenCode Sandbox Agent Guidelines

Это проект с настройками OpenCode и описаниями ролей для эффетивной разработки ПО.

## Build Commands
- **Clean**: `make clean` - Cleans Docker builder cache
- **Build**: `make build` - Builds OpenCode Docker image
- **Restore**: `make restore` - Syncs agent configs to local directory
- **Run**: `oc` - Use shell function to run OpenCode in sandbox

## Code Style Guidelines
В папке `opencode/agent` располагается описания ролей.
Эти роли позволяют эффективно разрабатывать программное обеспечение с помощью больших языковых моделей. 

## Testing
No testing assumed.

## Agent Workflow
- **Project Analyst**: Creates/updates PRDs and requirements
- **Orchestrator**: Декомпозирует задачу и делегируют выполнение подзадач другим агентам
- **Go/Python/Javascript Developer**: Writes production code following best practices
- **Test Engineer**: Implements comprehensive test suites
- **Code Reviewer**: Ensures quality, security, and performance
- **Technical Writer**: Maintains documentation and API docs

## Tools & Dependencies
- Use Context7 for library documentation and API references
