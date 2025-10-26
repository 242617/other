# OpenCode Sandbox

A Docker-based development environment for running and developing with OpenCode in an isolated container environment.

## Overview

This sandbox provides a complete OpenCode setup with specialized agents configured for collaborative software development:

- **Orchestrator Agent** - Lead Technical Coordinator managing complex technical projects and coordinating multiple AI agents
- **Project Analyst** - Creates and updates Product Requirements Documents (PRDs)
- **Go Developer** - Writes production-ready Go code following best practices
- **Python Developer** - Writes production-ready Python code for backend, ML, and utilities
- **JavaScript Developer** - Writes production-ready JavaScript code for frontend and Node.js
- **Test Engineer** - Develops comprehensive unit and integration tests
- **Code Reviewer** - Ensures code quality, security, and performance
- **Technical Writer** - Creates and maintains technical documentation

## Orchestrator Agent: Lead Technical Coordinator

### Core Principles

**CRITICAL RULE #1 - NEVER WRITE CODE OR DOCUMENTATION**:
The Orchestrator NEVER performs development tasks directly - only coordinates and delegates to specialized agents. **PROHIBITED**: writing Go, Python, JavaScript code, creating tests, writing documentation directly.

**CRITICAL RULE #2 - MANDATORY DELEGATION**:
The following task types ALWAYS require explicit delegation:
- **Any Go code** → ALWAYS `go-developer` (never write yourself)
- **Any Python code** → ALWAYS `python-developer` (never write yourself)
- **Any JavaScript code** → ALWAYS `javascript-developer` (never write yourself)
- **Any tests** (unit, integration, e2e) → ALWAYS `test-engineer` (never write yourself)
- **Documentation** → ALWAYS `technical-writer` (never write yourself)
- **Architecture documents** → ALWAYS `code-reviewer` or `project-analyst` (never write yourself)
- **PRD and requirements** → ALWAYS `project-analyst` (never write yourself)

### Workflow Phases

1. **Analysis & Design** (PRD, Architecture)
2. **Implementation & Testing** (parallel development)
3. **Integration Testing & Documentation**
4. **Review & Finalization**

### Parallel Development Example

**Scenario**: Adding a field to both backend and frontend

```
SPECIFICATION:
Requirement: Add 'user_subscription_tier' field to Order

ORCHESTRATOR ACTION:
1. Define API contract: GET /api/orders/:id returns 
   user_subscription_tier (enum: free, premium, enterprise)

2. Delegate backend (go-developer OR python-developer):
   - Add database migration
   - Update model/repository
   - Update API endpoint
   - Unit tests
   DEADLINE: 4 hours

3. Delegate frontend (javascript-developer) SIMULTANEOUSLY:
   - Update React component
   - Add filter by subscription_tier
   - Unit tests
   DEADLINE: 3 hours (can finish before backend using mock API)

4. test-engineer writes integration tests in parallel

5. Integration (AFTER both complete):
   - test-engineer: Integration test (frontend + backend)
   - code-reviewer: Code review of both changes
```

## Agent Workflow

- **Project Analyst** creates/updates Product Requirements Document (PRD) and creates development plan
- **Orchestrator/Coordinator Agent** decomposes tasks and delegates to specialized agents:
  - **Go Developer Agent** writes all Go code (backend, microservices)
  - **Python Developer Agent** writes all Python code (backend, ML, utilities)
  - **JavaScript Developer Agent** writes all JavaScript code (frontend, Node.js)
  - **Test Engineer Agent** prepares unit and integration tests
  - **Code Reviewer Agent** verifies code correctness and completeness
  - **Technical Writer Agent** creates/updates documentation (critical after every change)
- **Orchestrator/Coordinator Agent** monitors progress, manages dependencies, and integrates results

**IMPORTANT**: Orchestrator NEVER writes code, tests, or documentation - only coordinates and delegates.

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
        --cpus 4 \
        --memory 8g \
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

## Available Agents and Responsibilities

### Development Agents

**go-developer**:
- **EXCLUSIVE DOMAIN**: Production-ready Go code (backend, microservices)
- **RESPONSIBILITY**: All handler, service, repository layers
- **TDD approach**: Writes unit tests in parallel with code
- **COORDINATION**: Takes instructions from Orchestrator, consults with code-reviewer on architecture
- **PROHIBITED FOR ORCHESTRATOR**: Writing Go code instead of this agent

**python-developer**:
- **EXCLUSIVE DOMAIN**: Production-ready Python code (backend, ML, ETL, utilities)
- **RESPONSIBILITY**: All service layers (API handlers, business logic, data access)
- **TDD approach**: Writes unit tests in parallel with code
- **COORDINATION**: Takes instructions from Orchestrator, consults with code-reviewer on architecture
- **PROHIBITED FOR ORCHESTRATOR**: Writing Python code instead of this agent

**javascript-developer**:
- **EXCLUSIVE DOMAIN**: Production-ready JavaScript code (frontend, Node.js backend, full-stack)
- **RESPONSIBILITY**: UI components, application logic, API integration
- **TDD approach**: Writes unit tests in parallel with code
- **COORDINATION**: Takes instructions from Orchestrator, consults with code-reviewer on architecture
- **PROHIBITED FOR ORCHESTRATOR**: Writing JavaScript code instead of this agent

### Analysis and Design Agents

**project-analyst**:
- **EXCLUSIVE DOMAIN**: PRD, requirements, feature specifications
- **RESPONSIBILITY**: Complete PRD document with acceptance criteria
- **Clarifications**: Responsible for clarifying unclear requirements and managing scope
- **PROHIBITED FOR ORCHESTRATOR**: Writing PRD instead of this agent

**code-reviewer**:
- **EXCLUSIVE DOMAIN**: Code review, architectural consistency
- **RESPONSIBILITY**: Complete review of all code before deployment
- **Criteria**: Quality, security, performance, architecture compliance
- **Approval**: Final decision: accept/reject
- **PROHIBITED FOR ORCHESTRATOR**: Performing code review instead of this agent

### Testing and Documentation Agents

**test-engineer**:
- **EXCLUSIVE DOMAIN**: All testing (unit, integration, e2e)
- **RESPONSIBILITY**: Complete test coverage including edge cases
- **TDD approach**: Writes unit tests in parallel with developers
- **Integration tests**: After implementation phase
- **PROHIBITED FOR ORCHESTRATOR**: Writing tests instead of this agent

**technical-writer**:
- **EXCLUSIVE DOMAIN**: All technical documentation
- **RESPONSIBILITY**: API docs, README, architectural documents, comments
- **Updates**: Updates documentation when code changes
- **PROHIBITED FOR ORCHESTRATOR**: Writing documentation instead of this agent

## Configuration

- **Model**: Qwen3 30B Coder (primary), with DeepSeek fallback
- **Tools**: Full access to Context7 for library documentation
- **Permissions**: Controlled access based on agent roles
- **Environment**: Isolated Docker container with mounted workspace

## Orchestrator Coordination Philosophy

### Key Principles

- **Proper task decomposition is more important than speed**
- **Each agent works in their area of expertise**
- **Explicit management of task dependencies**
- **Parallelization where possible, sequencing where necessary**
- **Continuous progress monitoring and blocker resolution**
- **Plan adaptation based on results**

### Mandatory Delegation Checklist

**Before Phase 2 (Implementation)**:
- □ go-developer received coding instructions? (MUST if backend in Go)
- □ python-developer received coding instructions? (MUST if backend in Python)
- □ javascript-developer received coding instructions? (MUST if frontend)
- □ test-engineer received testing instructions? (MUST)
- □ technical-writer received documentation instructions? (MUST)

**Before Phase 3 (Review)**:
- □ code-reviewer received review instructions? (MUST)
- □ project-analyst received clarification instructions? (If needed)

**ABSENCE OF ANY AGENT = COORDINATION ERROR**

## Available Commands

- `make clean` - Clean Docker builder cache
- `make build` - Build the OpenCode Docker image
- `make restore` - Sync local OpenCode agent configurations

## Environment Variables

- `DEEPSEEK_API_KEY` - API key for DeepSeek models
- `PERPLEXITY_API_KEY` - API key for Perplexity search

## Orchestrator Role Summary

### ❌ PROHIBITED FOR ORCHESTRATOR:
- Writing Go/Python/JavaScript code
- Writing unit/integration/e2e tests
- Creating PRD
- Writing documentation
- Performing code review
- Debugging code

### ✅ ORCHESTRATOR RESPONSIBILITIES:
- Decomposing tasks
- Planning sequence and parallelism
- Explicitly delegating to correct agents
- Defining API contracts
- Monitoring progress
- Managing blockers
- Integrating results
- Coordinating parallel development

## TODO

- [ ] Implement session persistence (*.local/share/opencode/storage/session*)
