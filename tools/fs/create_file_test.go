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

func TestCreateFile(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "create new file",
			path:        "newfile.txt",
			expectError: false,
		},
		{
			name:        "create file in nested directory",
			path:        "subdir/file.txt",
			expectError: false,
		},
		{
			name:        "overwrite existing file",
			path:        "existing.txt",
			expectError: false,
		},
		{
			name:        "create file with empty path",
			path:        "",
			expectError: false, // The filesystem handles empty paths
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "overwrite existing file" {
				require.NoError(t, afero.WriteFile(memFs, tt.path, []byte("Old content"), 0644), "afero write file")
			}

			createCmd := fs.CreateFile()

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
			require.NoError(t, err, "Cannot check if file exists")
			assert.True(t, exists, "File should exist after creation")

			// create_file only creates empty files, so content should be empty
			content, err := afero.ReadFile(memFs, tt.path)
			require.NoError(t, err, "Cannot read file content")
			assert.Empty(t, content, "File should be empty after creation")
		})
	}
}
