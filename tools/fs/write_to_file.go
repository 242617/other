package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/242617/other/tools"
)

func (fs *FS) WriteToFile() *Command {
	name := "fs_write_to_file"
	description := "Request to write content to a file at the specified path. If the file exists, it will be overwritten with the provided content. If the file doesn't exist, it will be created. This tool will automatically create any directories needed to write the file."
	type argsStruct struct {
		Path    string `json:"path" description:"The path of the file to write to (relative to the current working directory)"`
		Content string `json:"content" description:"The content to write to the file. ALWAYS provide the COMPLETE intended content of the file, without any truncation or omissions. You MUST include ALL parts of the file, even if they haven't been modified."`
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

		f, err := fs.OpenFile(args.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Sprintf("Cannot open file %q due to error: %q", args.Path, err.Error())
		}
		defer f.Close()

		if _, err := io.WriteString(f, args.Content); err != nil {
			return fmt.Sprintf("Cannot write string into %q due to error: %q", args.Path, err.Error())
		}
		return fmt.Sprintf("✅ Content written into %q successfully", args.Path)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}
