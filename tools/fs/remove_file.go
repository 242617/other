package fs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/242617/other/tools"
)

func (fs *FS) RemoveFile() *Command {
	name := "fs_remove_file"
	description := "Remove file."
	type argsStruct struct {
		Path string `json:"path" description:"File path."`
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
			return fmt.Sprintf("Failed to removed file %q due to error: %q", args.Path, err.Error())
		}
		return fmt.Sprintf("File %q removed successfully ✅", args.Path)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}
