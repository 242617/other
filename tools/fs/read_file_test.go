package fs_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/other/tools/fs"
)

func TestReadFile(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	testContent := "Hello, World!\nThis is a test file.\nWith multiple lines."
	testFile := "testfile.txt"
	err := afero.WriteFile(memFs, testFile, []byte(testContent), 0644)
	require.NoError(t, err, "Cannot write test file")

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "read existing file",
			path:        testFile,
			expectError: false,
		},
		{
			name:        "read non-existent file",
			path:        "nonexistent.txt",
			expectError: true,
		},
		{
			name:        "read directory instead of file",
			path:        "somedir",
			expectError: false, // The tool doesn't validate if it's a directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "read directory instead of file" {
				err := memFs.MkdirAll(tt.path, 0755)
				require.NoError(t, err, "Cannot create test directory")
			}

			readCmd := fs.ReadFile()

			args := map[string]any{
				"path": tt.path,
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := readCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result), "Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result), "Unexpected error: %s", result)

			if tt.name == "read directory instead of file" {
				// Directories return empty content when read as files
				assert.Empty(t, result, "Directory should return empty content")
			} else {
				assert.Equal(t, testContent, result, "Content should match expected value")
			}
		})
	}
}

// read_file: Read binary file contents
// Example: {"path": "binaryfile.bin"}
func TestReadFileBinaryContent(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	binaryContent := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x57, 0x6f, 0x72, 0x6c, 0x64}
	binaryFile := "binary.bin"
	err := afero.WriteFile(memFs, binaryFile, binaryContent, 0644)
	require.NoError(t, err, "Cannot write binary test file")

	readCmd := fs.ReadFile()

	args := map[string]interface{}{
		"path": binaryFile,
	}

	argsJSON, err := json.Marshal(args)
	require.NoError(t, err, "Cannot marshal args")

	ctx := context.Background()
	result := readCmd.Call(ctx, string(argsJSON))

	assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
	assert.Len(t, result, len(binaryContent), "Binary content length should match")
}
