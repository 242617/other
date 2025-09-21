package fs_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/other/tools/fs"
)

// list_files: List files and directories (supports recursive listing)
// Example: {"path": "mydirectory", "recursive": true}
func TestListFiles(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	testDir := "testdir"
	err := memFs.MkdirAll(testDir+"/subdir1", 0755)
	require.NoError(t, err, "Cannot create test directory")

	err = memFs.MkdirAll(testDir+"/subdir2", 0755)
	require.NoError(t, err, "Cannot create test directory")

	err = afero.WriteFile(memFs, testDir+"/file1.txt", []byte("content1"), 0644)
	require.NoError(t, err, "Cannot create test file")
	err = afero.WriteFile(memFs, testDir+"/file2.go", []byte("content2"), 0644)
	require.NoError(t, err, "Cannot create test file")
	err = afero.WriteFile(memFs, testDir+"/subdir1/file3.txt", []byte("content3"), 0644)
	require.NoError(t, err, "Cannot create test file")
	err = afero.WriteFile(memFs, testDir+"/subdir2/file4.go", []byte("content4"), 0644)
	require.NoError(t, err, "Cannot create test file")

	tests := []struct {
		name        string
		path        string
		recursive   bool
		wantLines   int
		expectError bool
	}{
		{
			name:        "top level listing",
			path:        testDir,
			recursive:   false,
			wantLines:   4, // file1.txt, file2.go, subdir1/, subdir2/
			expectError: false,
		},
		{
			name:        "recursive listing",
			path:        testDir,
			recursive:   true,
			wantLines:   6, // all files and directories recursively
			expectError: false,
		},
		{
			name:        "non-existent directory",
			path:        "nonexistent",
			recursive:   false,
			wantLines:   0,
			expectError: true,
		},
		{
			name:        "file instead of directory",
			path:        testDir + "/file1.txt",
			recursive:   false,
			wantLines:   1, // Error message should be returned as a single line
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listCmd := fs.ListFiles()

			args := map[string]interface{}{
				"path":      tt.path,
				"recursive": tt.recursive,
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := listCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result) || strings.Contains(result, "is not a directory"),
					"Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result) || strings.Contains(result, "is not a directory"),
				"Unexpected error: %s", result)

			lines := strings.Split(result, "\n")
			assert.Len(t, lines, tt.wantLines, "Number of lines should match expected count")

			if tt.recursive {
				assert.True(t, containsLine(lines, "subdir1/"), "Should find subdir1/ in recursive listing")
				assert.True(t, containsLine(lines, "subdir1/file3.txt"), "Should find subdir1/file3.txt in recursive listing")
			} else {
				assert.True(t, containsLine(lines, "subdir1/"), "Should find subdir1/ in top-level listing")
				assert.False(t, containsLine(lines, "subdir1/file3.txt"), "Should not find subdir1/file3.txt in top-level listing")
			}
		})
	}
}

func TestListFilesEmptyDirectory(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	testDir := "emptydir"
	err := memFs.MkdirAll(testDir, 0755)
	require.NoError(t, err, "Cannot create empty directory for test")

	listCmd := fs.ListFiles()
	args := map[string]interface{}{
		"path": testDir,
	}

	argsJSON, err := json.Marshal(args)
	require.NoError(t, err, "Cannot marshal args")

	ctx := context.Background()
	result := listCmd.Call(ctx, string(argsJSON))

	assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
	assert.Contains(t, result, "No files found", "Should indicate no files found")
}

func containsLine(lines []string, substr string) bool {
	for _, line := range lines {
		if strings.Contains(line, substr) {
			return true
		}
	}
	return false
}
