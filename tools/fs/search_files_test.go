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

func TestSearchFiles(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	testDir := "testdir"
	err := memFs.MkdirAll(testDir+"/subdir", 0755)
	require.NoError(t, err, "Cannot create test directory")

	content1 := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    var testVar = 42
    fmt.Println("The answer is:", testVar)
}`

	content2 := `// This is a test file
const testConst = "hello"
var anotherVar = testConst + " world"

func testFunction() {
    fmt.Println(testConst)
}`

	err = afero.WriteFile(memFs, testDir+"/main.go", []byte(content1), 0644)
	require.NoError(t, err, "Cannot create test file")
	err = afero.WriteFile(memFs, testDir+"/subdir/utils.go", []byte(content2), 0644)
	require.NoError(t, err, "Cannot create test file")
	err = afero.WriteFile(memFs, testDir+"/readme.txt", []byte("This is a text file with test content"), 0644)
	require.NoError(t, err, "Cannot create test file")

	tests := []struct {
		name        string
		path        string
		pattern     string
		filePattern string
		wantLines   int
		expectError bool
	}{
		{
			name:        "basic regex search",
			path:        testDir,
			pattern:     `fmt\.Println`,
			wantLines:   2,
			expectError: false,
		},
		{
			name:        "variable search",
			path:        testDir,
			pattern:     `testVar`,
			wantLines:   1,
			expectError: false,
		},
		{
			name:        "file pattern filtering",
			path:        testDir,
			pattern:     `test`,
			filePattern: "*.go",
			wantLines:   3,
			expectError: false,
		},
		{
			name:        "no matches found",
			path:        testDir,
			pattern:     `nonexistentpattern`,
			wantLines:   1,
			expectError: false,
		},
		{
			name:        "invalid regex pattern",
			path:        testDir,
			pattern:     `[invalid`,
			wantLines:   1,
			expectError: true, // The tool should error on invalid regex patterns
		},
		{
			name:        "non-existent directory",
			path:        "nonexistent",
			pattern:     `test`,
			wantLines:   1,
			expectError: true,
		},
		{
			name:        "file instead of directory",
			path:        testDir + "/main.go",
			pattern:     `test`,
			wantLines:   1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchCmd := fs.SearchFiles()

			args := map[string]interface{}{
				"path":    tt.path,
				"pattern": tt.pattern,
			}
			if tt.filePattern != "" {
				args["file_pattern"] = tt.filePattern
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := searchCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result) || strings.Contains(result, "is not a directory"),
					"Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result) || strings.Contains(result, "is not a directory"),
				"Unexpected error: %s", result)

			lines := strings.Split(result, "\n")
			assert.GreaterOrEqual(t, len(lines), tt.wantLines, "Should have at least expected number of lines")

			if strings.Contains(result, "No matches found") && tt.wantLines > 1 {
				assert.Fail(t, "Expected matches but got 'No matches found'")
			}
		})
	}
}

// search_files: Test that search results contain proper context information
func TestSearchFilesContextFormat(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	testContent := `package main

func main() {
    // This is a test function
    fmt.Println("Hello, World!")
}`

	testFile := "context_test.go"
	err := afero.WriteFile(memFs, testFile, []byte(testContent), 0644)
	require.NoError(t, err, "Cannot create test file for context format test")

	searchCmd := fs.SearchFiles()

	args := map[string]interface{}{
		"path":    ".",
		"pattern": `fmt\.Println`,
	}

	argsJSON, err := json.Marshal(args)
	require.NoError(t, err, "Cannot marshal args")

	ctx := context.Background()
	result := searchCmd.Call(ctx, string(argsJSON))

	assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
	assert.Contains(t, result, "File: context_test.go", "Should contain file path")
	assert.Contains(t, result, "Line:", "Should contain line number")
	assert.Contains(t, result, "Context:", "Should contain context information")
	assert.Contains(t, result, `fmt.Println("Hello, World!")`, "Should contain matching line")
}
