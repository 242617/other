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

func TestReplaceInFile(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	initialContent := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    var x = 42
    fmt.Println("The answer is:", x)
}`

	testFile := "test.go"
	err := afero.WriteFile(memFs, testFile, []byte(initialContent), 0644)
	require.NoError(t, err, "Cannot create initial test file")

	tests := []struct {
		name        string
		diff        string
		wantContent string
		expectError bool
	}{
		{
			name: "single replacement",
			diff: `------- SEARCH
    fmt.Println("Hello, World!")
=======
    fmt.Println("Hello, Go!")
+++++++ REPLACE`,
			wantContent: `package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
    var x = 42
    fmt.Println("The answer is:", x)
}`,
			expectError: false,
		},
		{
			name: "multiple replacements",
			diff: `------- SEARCH
    fmt.Println("Hello, World!")
=======
    fmt.Println("Hello, Go!")
+++++++ REPLACE
------- SEARCH
    var x = 42
    fmt.Println("The answer is:", x)
=======
    var answer = 42
    fmt.Println("The answer is:", answer)
+++++++ REPLACE`,
			wantContent: `package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
    var answer = 42
    fmt.Println("The answer is:", answer)
}`,
			expectError: false,
		},
		{
			name: "search content not found",
			diff: `------- SEARCH
    nonexistent content
=======
    new content
+++++++ REPLACE`,
			wantContent: initialContent,
			expectError: true,
		},
		{
			name: "malformed diff block",
			diff: `------- SEARCH
    fmt.Println("Hello, World!")
=======
    fmt.Println("Hello, Go!")
+++++++`,
			wantContent: initialContent,
			expectError: true,
		},
		{
			name:        "empty diff",
			diff:        "",
			wantContent: initialContent,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset file content for each test
			err := afero.WriteFile(memFs, testFile, []byte(initialContent), 0644)
			require.NoError(t, err, "Cannot reset file content for test")

			replaceCmd := fs.ReplaceInFile()

			args := map[string]interface{}{
				"path": testFile,
				"diff": tt.diff,
			}

			argsJSON, err := json.Marshal(args)
			require.NoError(t, err, "Cannot marshal args")

			ctx := context.Background()
			result := replaceCmd.Call(ctx, string(argsJSON))

			if tt.expectError {
				assert.True(t, isErrorResult(result), "Expected error but got success: %s", result)
				return
			}

			assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
			assert.Contains(t, result, "success", "Expected success message")

			content, err := afero.ReadFile(memFs, testFile)
			require.NoError(t, err, "Cannot read file content")
			assert.Equal(t, tt.wantContent, string(content), "File content should match expected value")
		})
	}
}

// replace_in_file: Test complex replacement scenarios
func TestReplaceInFileComplex(t *testing.T) {
	memFs := afero.NewMemMapFs()
	fs := &fs.FS{Fs: memFs}

	initialContent := `// Header comment
package main

import "fmt"

func main() {
    // Print greeting
    fmt.Println("Hello, World!")
    
    // Calculate result
    result := 42
    fmt.Println("Result:", result)
}`

	testFile := "complex.go"
	err := afero.WriteFile(memFs, testFile, []byte(initialContent), 0644)
	require.NoError(t, err, "Cannot create complex test file")

	complexDiff := `------- SEARCH
    // Print greeting
    fmt.Println("Hello, World!")
=======
    // Print customized greeting
    fmt.Println("Hello, Go Developer!")
+++++++ REPLACE
------- SEARCH
    // Calculate result
    result := 42
    fmt.Println("Result:", result)
=======
    // Calculate and display result
    answer := 42
    fmt.Println("The answer is:", answer)
+++++++ REPLACE`

	replaceCmd := fs.ReplaceInFile()

	args := map[string]interface{}{
		"path": testFile,
		"diff": complexDiff,
	}

	argsJSON, err := json.Marshal(args)
	require.NoError(t, err, "Cannot marshal args")

	ctx := context.Background()
	result := replaceCmd.Call(ctx, string(argsJSON))

	assert.False(t, isErrorResult(result), "Unexpected error: %s", result)
	assert.Contains(t, result, "success", "Expected success message")

	expectedContent := `// Header comment
package main

import "fmt"

func main() {
    // Print customized greeting
    fmt.Println("Hello, Go Developer!")
    
    // Calculate and display result
    answer := 42
    fmt.Println("The answer is:", answer)
}`

	content, err := afero.ReadFile(memFs, testFile)
	require.NoError(t, err, "Cannot read file content")
	assert.Equal(t, expectedContent, string(content), "File content should match expected value after complex replacement")
}
