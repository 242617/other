package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/242617/other/tools"
)

func (fs *FS) ReadFile() *Command {
	name := "fs_read_file"
	description := "Request to read the contents of a file at the specified path. Use this when you need to examine the contents of an existing file you do not know the contents of, for example to analyze code, review text files, or extract information from configuration files. Automatically extracts raw text from PDF and DOCX files. May not be suitable for other types of binary files, as it returns the raw content as a string. Do NOT use this tool to list the contents of a directory. Only use this tool on files."
	type argsStruct struct {
		Path string `json:"path" description:"The path of the file to read (relative to the current working directory)"`
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

		f, err := fs.OpenFile(args.Path, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Sprintf("Cannot open file %q due to error: %q", args.Path, err.Error())
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			return fmt.Sprintf("Cannot read all from %q due to error: %q", args.Path, err.Error())
		}
		return string(b)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}
