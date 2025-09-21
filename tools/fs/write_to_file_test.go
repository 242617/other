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

func TestWriteToFile(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	tests := []struct {
		name        string
		path        string
		content     string
		create      bool
		expectError bool
	}{
		{
			name:        "write to new file",
			path:        "newfile.txt",
			content:     "Hello, World!",
			create:      false,
			expectError: false,
		},
		{
			name:        "overwrite existing file",
			path:        "existing.txt",
			content:     "New content",
			create:      true,
			expectError: false,
		},
		{
			name:        "write to nested path",
			path:        "subdir/file.txt",
			content:     "Nested content",
			create:      false,
			expectError: false,
		},
		{
			name:        "write empty content",
			path:        "empty.txt",
			content:     "",
			create:      false,
			expectError: false,
		},
		{
			name:        "write to empty path",
			path:        "",
			content:     "content",
			create:      false,
			expectError: false, // The filesystem handles empty paths
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.create {
				err := afero.WriteFile(memFs, tt.path, []byte("Old content"), 0644)
				require.NoError(t, err, "Cannot create existing file for test")
			}

			writeCmd := fs.WriteToFile()

			args := map[string]interface{}{
				"path":    tt.path,
				"content": tt.content,
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := writeCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result), "Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
			assert.Contains(t, result, "success", "Expected success message")

			exists, err := afero.Exists(memFs, tt.path)
			require.NoError(t, err, "Cannot check if file exists")
			assert.True(t, exists, "File should exist after writing")

			content, err := afero.ReadFile(memFs, tt.path)
			require.NoError(t, err, "Cannot read file content")
			assert.Equal(t, tt.content, string(content), "File content should match expected value")
		})
	}
}
