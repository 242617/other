package fs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/242617/other/tools"
)

func (fs *FS) CreateDirectory() *Command {
	name := "fs_create_directory"
	description := "Create a directory in the file system. Creates parent directories if they don't exist."
	type argsStruct struct {
		Path string `json:"path" description:"Directory path"`
	}

	toolInfo, err := tools.CreateToolInfo(name, description, argsStruct{})
	if err != nil {
		return empty(name, err)
	}

	call := func(ctx context.Context, raw string) string {
		var args argsStruct
		if err := json.Unmarshal([]byte(raw), &args); err != nil {
			return fmt.Sprintf("Cannot unmarshal arguments due to error: %q", err.Error())
		}
		if err := fs.MkdirAll(args.Path, 0755); err != nil {
			return fmt.Sprintf("Failed to create directory %q due to error: %q", args.Path, err.Error())
		}
		return fmt.Sprintf("✅ Directory %q created successfully", args.Path)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}
