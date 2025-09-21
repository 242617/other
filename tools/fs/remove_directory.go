package fs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/242617/other/tools"
)

func (fs *FS) RemoveDirectory() *Command {
	name := "fs_remove_directory"
	description := "Remove a directory from the file system."
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
		if err := fs.Remove(args.Path); err != nil {
			return fmt.Sprintf("Failed to remove directory %q due to error: %q", args.Path, err.Error())
		}
		return fmt.Sprintf("✅ Directory %q removed successfully", args.Path)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}
