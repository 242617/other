package fs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/242617/other/tools"
)

func (fs *FS) CreateFile() *Command {
	name := "fs_create_file"
	description := "Create a directory in the file system. Creates parent directories if they don't exist."
	type argsStruct struct {
		Path string `json:"path" description:"File path"`
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
		if _, err := fs.Create(args.Path); err != nil {
			return fmt.Sprintf("Failed to create file %q due to error: %q", args.Path, err.Error())
		}
		return fmt.Sprintf("✅ File %q created successfully", args.Path)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}
