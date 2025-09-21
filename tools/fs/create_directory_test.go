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

func TestCreateDirectory(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "create new directory",
			path:        "newdir",
			expectError: false,
		},
		{
			name:        "create nested directory",
			path:        "parent/child",
			expectError: false,
		},
		{
			name:        "create directory that already exists",
			path:        "existingdir",
			expectError: false,
		},
		{
			name:        "create directory with empty path",
			path:        "",
			expectError: false, // The filesystem handles empty paths by creating in current dir
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "create directory that already exists" {
				err := memFs.MkdirAll(tt.path, 0755)
				require.NoError(t, err, "Cannot create existing directory for test")
			}

			createCmd := fs.CreateDirectory()

			args := map[string]interface{}{
				"path": tt.path,
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := createCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result), "Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
			assert.Contains(t, result, "success", "Expected success message")

			exists, err := afero.Exists(memFs, tt.path)
			require.NoError(t, err, "Cannot check if directory exists")
			assert.True(t, exists, "Directory should exist after creation")

			isDir, err := afero.IsDir(memFs, tt.path)
			require.NoError(t, err, "Cannot check if path is directory")
			assert.True(t, isDir, "Path should be a directory")
		})
	}
}
