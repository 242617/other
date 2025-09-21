package tools

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

var (
	dockerCommandName = "docker"
	dockerCommandBase = []string{"exec", "-it", "sandbox", "sh", "-c"}
)

func CommandRun(ctx context.Context, commands ...string) error {
	commands = append(dockerCommandBase, strings.Join(commands, " "))
	cmd := exec.CommandContext(ctx, dockerCommandName, commands...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
