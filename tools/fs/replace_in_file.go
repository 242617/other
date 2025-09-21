package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/242617/other/tools"
)

func (fs *FS) ReplaceInFile() *Command {
	name := "fs_replace_in_file"
	description := "Request to replace sections of content in an existing file using SEARCH/REPLACE blocks that define exact changes to specific parts of the file. This tool should be used when you need to make targeted changes to specific parts of a file."
	type argsStruct struct {
		Path string `json:"path" description:"The path of the file to modify (relative to the current working directory)"`
		Diff string `json:"diff" description:"One or more SEARCH/REPLACE blocks following this exact format:\n------- SEARCH\n[exact content to find]\n=======\n[new content to replace with]\n+++++++ REPLACE"`
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

		contentBytes, err := io.ReadAll(f)
		if err != nil {
			return fmt.Sprintf("Cannot read all from %q due to error: %q", args.Path, err.Error())
		}
		content := string(contentBytes)

		if err := f.Close(); err != nil {
			return fmt.Sprintf("cannot close file %q due to error: %q", args.Path, err.Error())
		}

		modified, err := applyDiff(content, args.Diff)
		if err != nil {
			return fmt.Sprintf("Cannot apply diff to %q due to error: %q", args.Path, err.Error())
		}

		if content == modified {
			return fmt.Sprintf("No changes made to %q - SEARCH patterns not found", args.Path)
		}

		f, err = fs.OpenFile(args.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Sprintf("Cannot open file %q for writing due to error: %q", args.Path, err.Error())
		}
		defer f.Close()

		if _, err := io.WriteString(f, modified); err != nil {
			return fmt.Sprintf("Cannot write modified content to %q due to error: %q", args.Path, err.Error())
		}
		return fmt.Sprintf("✅ Content replaced in %q successfully", args.Path)
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}

// applyDiff parses and applies SEARCH/REPLACE blocks to the content
func applyDiff(content, diff string) (string, error) {
	blocks := parseDiffBlocks(diff)
	if len(blocks) == 0 {
		return content, fmt.Errorf("no valid SEARCH/REPLACE blocks found in diff")
	}

	modifiedContent := content
	for _, block := range blocks {
		if !strings.Contains(modifiedContent, block.search) {
			return content, fmt.Errorf("SEARCH content not found in file: %q", strings.TrimSpace(block.search))
		}

		// Replace only the first occurrence (as specified in requirements)
		modifiedContent = strings.Replace(modifiedContent, block.search, block.replace, 1)
	}

	return modifiedContent, nil
}

type diffBlock struct {
	search, replace string
}

// parseDiffBlocks parses the diff string into individual SEARCH/REPLACE blocks
func parseDiffBlocks(diff string) []diffBlock {
	var blocks []diffBlock
	lines := strings.Split(diff, "\n")

	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "------- SEARCH" {
			var searchLines []string
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) == "=======" {
					var replaceLines []string
					for k := j + 1; k < len(lines); k++ {
						if strings.TrimSpace(lines[k]) == "+++++++ REPLACE" {
							block := diffBlock{
								search:  strings.Join(searchLines, "\n"),
								replace: strings.Join(replaceLines, "\n"),
							}
							blocks = append(blocks, block)
							i = k
							break
						}
						replaceLines = append(replaceLines, lines[k])
					}
					break
				}
				searchLines = append(searchLines, lines[j])
			}
		}
	}

	return blocks
}
