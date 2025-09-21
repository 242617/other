package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/242617/other/tools"
)

func (fs *FS) ListFiles() *Command {
	name := "fs_list_files"
	description := "Request to list files and directories within the specified directory. If recursive is true, it will list all files and directories recursively. If recursive is false or not provided, it will only list the top-level contents. Do not use this tool to confirm the existence of files you may have created, as the user will let you know if the files were created successfully or not."
	type argsStruct struct {
		Path      string `json:"path" description:"The path of the directory to list contents for (relative to the current working directory)"`
		Recursive bool   `json:"recursive,omitempty" description:"Whether to list files recursively. Use true for recursive listing, false or omit for top-level only."`
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

		var entries []string
		if args.Recursive {
			entries, err = listFilesRecursive(fs, args.Path)
		} else {
			entries, err = listFilesTopLevel(fs, args.Path)
		}
		if err != nil {
			return fmt.Sprintf("Cannot list files in %q due to error: %q", args.Path, err.Error())
		}

		if len(entries) == 0 {
			return fmt.Sprintf("No files found in %q", args.Path)
		}

		return strings.Join(entries, "\n")
	}

	return &Command{
		name: name,
		info: toolInfo,
		call: call,
	}
}

func listFilesTopLevel(fs *FS, path string) ([]string, error) {
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var entries []string
	for _, info := range fileInfos {
		if info.IsDir() {
			entries = append(entries, info.Name()+"/")
		} else {
			entries = append(entries, info.Name())
		}
	}

	sort.Strings(entries)
	return entries, nil
}

func listFilesRecursive(fs *FS, path string) ([]string, error) {
	var entries []string

	err := recursiveWalk(fs, path, path, &entries)
	if err != nil {
		return nil, err
	}

	sort.Strings(entries)
	return entries, nil
}

func recursiveWalk(fs *FS, basePath, currentPath string, entries *[]string) error {
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
		relPath, err := filepath.Rel(basePath, fullPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Add trailing slash for directories
			*entries = append(*entries, relPath+"/")
			err := recursiveWalk(fs, basePath, fullPath, entries)
			if err != nil {
				return err
			}
		} else {
			// Files without any prefix
			*entries = append(*entries, relPath)
		}
	}

	return nil
}
