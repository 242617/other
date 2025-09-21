package fs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/242617/other/tools"
)

func (fs *FS) SearchFiles() *Command {
	name := "fs_search_files"
	description := "Perform regex searches across files in directories with context-rich results. Searches through all subdirectories from the starting path. The search is case-sensitive and matches the exact regex pattern. Returns full paths to all matching items with context around matches."
	type argsStruct struct {
		Path        string `json:"path" description:"The path of the directory to search in (relative to the current working directory)"`
		Pattern     string `json:"pattern" description:"The regular expression pattern to search for"`
		FilePattern string `json:"file_pattern,omitempty" description:"Glob pattern to filter files (e.g., '*.go' for Go files, '*.txt' for text files). If not provided, it will search all files."`
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

		info, err := fs.Stat(args.Path)
		if err != nil {
			return fmt.Sprintf("Cannot access directory %q due to error: %q", args.Path, err.Error())
		}
		if !info.IsDir() {
			return fmt.Sprintf("Path %q is not a directory", args.Path)
		}

		regex, err := regexp.Compile(args.Pattern)
		if err != nil {
			return fmt.Sprintf("Invalid regex pattern %q due to error: %q", args.Pattern, err.Error())
		}

		results, err := searchFilesRecursive(fs, args.Path, regex, args.FilePattern)
		if err != nil {
			return fmt.Sprintf("Cannot search files in %q due to error: %q", args.Path, err.Error())
		}

		if len(results) == 0 {
			return fmt.Sprintf("No matches found for pattern %q in %q", args.Pattern, args.Path)
		}

		return strings.Join(results, "\n\n")
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}

func searchFilesRecursive(fs *FS, path string, regex *regexp.Regexp, filePattern string) ([]string, error) {
	var results []string

	err := recursiveSearchWalk(fs, path, path, regex, filePattern, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func recursiveSearchWalk(fs *FS, basePath, currentPath string, regex *regexp.Regexp, filePattern string, results *[]string) error {
	file, err := fs.Open(currentPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for _, info := range fileInfos {
		fullPath := filepath.Join(currentPath, info.Name())

		if info.IsDir() {
			err := recursiveSearchWalk(fs, basePath, fullPath, regex, filePattern, results)
			if err != nil {
				return err
			}
		} else {
			if filePattern != "" {
				matched, err := filepath.Match(filePattern, info.Name())
				if err != nil || !matched {
					continue
				}
			}

			fileMatches, err := searchInFile(fs, fullPath, regex, basePath)
			if err != nil {
				return err
			}

			if len(fileMatches) > 0 {
				*results = append(*results, fileMatches...)
			}
		}
	}

	return nil
}

func searchInFile(fs *FS, filePath string, regex *regexp.Regexp, basePath string) ([]string, error) {
	file, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		relPath = filePath
	}

	var matches []string
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	var contextLines []string

	for scanner.Scan() {
		line := scanner.Text()
		if regex.MatchString(line) {
			context := getContext(contextLines, line)
			match := fmt.Sprintf("File: %s\nLine: %d\n%s", relPath, lineNumber, context)
			matches = append(matches, match)
		}

		contextLines = append(contextLines, line)
		if len(contextLines) > 5 {
			contextLines = contextLines[1:]
		}

		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func getContext(contextLines []string, currentLine string) string {
	var context []string
	if len(contextLines) > 0 {
		context = append(context, "Context:")
		for i, line := range contextLines {
			context = append(context, fmt.Sprintf("  %d: %s", len(contextLines)-i, line))
		}
		context = append(context, fmt.Sprintf("  -> %s", currentLine))
	} else {
		context = append(context, fmt.Sprintf("Match: %s", currentLine))
	}
	return strings.Join(context, "\n")
}
