# mdprc - Markdown Processor

A Go-based document processing tool that combines Markdown frontmatter with Go template capabilities for dynamic document generation.

## Overview

mdprc (Markdown Processor) is a command-line tool that processes Markdown files with YAML frontmatter and executes Go templates, providing a flexible system for generating dynamic documentation, reports, and other text-based content.

## Usage

### Basic Command

```bash
./mdprc -source ./src -target ./build
```

### Command Line Options

- `-source`: Source directory containing markdown files (required)
- `-target`: Target directory for processed output (required)

## mdprc Properties

mdprc supports special properties in YAML frontmatter that control processing behavior:

### Available Properties

| Property                    | Type    | Description                             |
|-----------------------------|---------|-----------------------------------------|
| `[mdprc:skip_execute]`      | boolean | Skip template execution for this file   |
| `[mdprc:skip_place]`        | boolean | Skip file placement in target directory |
| `[mdprc:remove_properties]` | boolean | Remove all mdprc properties from output |

### Property Usage Examples

```yaml
---
description: "Example Document"
"[mdprc:skip_execute]": false
"[mdprc:skip_place]": false
"[mdprc:remove_properties]": true
---
```

## Template Functions

### Built-in Functions

#### `include`

Include content from other markdown files:

```go
{{ include "path/to/file.md" }}
```

**Example:**
```markdown
{{ include "../documents/task.md" }}
```

This function:
- Resolves paths relative to the current file
- Processes the included file (including frontmatter and templates)
- Returns the processed content as a string

## TODO

- [ ] Resolve recursive includes
