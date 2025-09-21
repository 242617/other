// Package fs provides file system operations for tools.
//
// Available tools with usage examples:
//
// create_directory: Create directories with automatic parent directory creation
//
//	Example: {"path": "mydirectory"}
//
// create_file: Create empty files (does not write content)
//
//	Example: {"path": "emptyfile.txt"}
//
// read_file: Read file contents (supports text and binary files)
//
//	Example: {"path": "myfile.txt"}
//
// remove_directory: Remove directories (works on any filesystem object)
//
//	Example: {"path": "directory_to_remove"}
//
// remove_file: Remove files (works on any filesystem object)
//
//	Example: {"path": "file_to_remove.txt"}
//
// replace_in_file: Make targeted edits using SEARCH/REPLACE blocks
//
//	Example: {"path": "file.txt", "diff": "------- SEARCH\nold content\n=======\nnew content\n+++++++ REPLACE"}
//
// write_to_file: Write content to files (creates or overwrites files)
//
//	Example: {"path": "myfile.txt", "content": "Hello, World!"}
//
// list_files: List files and directories (supports recursive listing)
//
//	Example: {"path": "mydirectory", "recursive": true}
//
// search_files: Perform regex searches across files with context-rich results
//
//	Example: {"path": "src", "pattern": "func.*main", "file_pattern": "*.go"}
package fs

import (
	"context"
	"fmt"

	"github.com/spf13/afero"

	"github.com/242617/other/agent"

	"github.com/242617/other/tools"
)

func New(dir string) *FS {
	return &FS{Fs: afero.NewBasePathFs(afero.NewOsFs(), dir)}
}

type FS struct{ afero.Fs }

type Command struct {
	name string
	info agent.ToolInfo
	call tools.CallFunc
}

func (cmd *Command) Name() string                                { return cmd.name }
func (cmd *Command) Info() agent.ToolInfo                        { return cmd.info }
func (cmd *Command) Call(ctx context.Context, raw string) string { return cmd.call(ctx, raw) }

func empty(name string, err error) *Command {
	return &Command{name: name, call: func(context.Context, string) string {
		return fmt.Sprintf("Cannot create tool info due to error: %q", err)
	}}
}
