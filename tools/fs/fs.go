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
