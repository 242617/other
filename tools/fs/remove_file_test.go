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

func TestRemoveFile(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	tests := []struct {
		name        string
		path        string
		create      bool
		expectError bool
	}{
		{
			name:        "remove existing file",
			path:        "testfile.txt",
			create:      true,
			expectError: false,
		},
		{
			name:        "remove non-existent file",
			path:        "nonexistent.txt",
			create:      false,
			expectError: true,
		},
		{
			name:        "remove directory instead of file",
			path:        "somedir",
			create:      true,
			expectError: false, // The tool removes any filesystem object
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.create {
				if tt.name == "remove directory instead of file" {
					err := memFs.MkdirAll(tt.path, 0755)
					require.NoError(t, err, "Cannot create test directory")
				} else {
					err := afero.WriteFile(memFs, tt.path, []byte("content"), 0644)
					require.NoError(t, err, "Cannot create test file")
				}
			}

			removeCmd := fs.RemoveFile()

			args := map[string]interface{}{
				"path": tt.path,
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := removeCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result), "Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
			assert.Contains(t, result, "success", "Expected success message")

			exists, err := afero.Exists(memFs, tt.path)
			require.NoError(t, err, "Cannot check if file exists")
			assert.False(t, exists, "File should not exist after removal")
		})
	}
}
