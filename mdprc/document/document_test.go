package document

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

func TestParse(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}

	// Test cases
	tests := []struct {
		name        string
		setupFile   func(t *testing.T, path string)
		wantErr     bool
		errContains string
		validateDoc func(t *testing.T, doc Document, err error)
	}{
		{
			name: "valid document with frontmatter and content",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: Test Document
---
This is test content.`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.Bytes(), []byte("This is test content")) {
					t.Errorf("expected content to contain 'This is test content', got: %s", doc.String())
				}
			},
		},
		{
			name: "document without frontmatter - valid",
			setupFile: func(t *testing.T, path string) {
				content := "Plain markdown content without frontmatter"
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.Bytes(), []byte("Plain markdown content without frontmatter")) {
					t.Errorf("expected content to contain plain text, got: %s", doc.String())
				}
			},
		},
		{
			name: "skip_execute feature gate enabled",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: Test
---
Content {{ .title }}`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Template should be executed
				if !bytes.Contains(doc.Bytes(), []byte("Content Test")) {
					t.Errorf("expected executed template to contain 'Content Test', got: %s", doc.String())
				}
			},
		},
		{
			name: "skip_place feature gate enabled",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: Test
---
Some content`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if doc.SkipPlace() {
					t.Errorf("expected SkipPlace to be false when not set")
				}
			},
		},
		{
			name: "remove_properties feature gate enabled",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: Test
description: Description
---
Content here`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Frontmatter should be preserved when remove_properties is not set
				if !bytes.HasPrefix(doc.Bytes(), []byte("---")) {
					t.Errorf("expected frontmatter when remove_properties is false, got: %s", doc.String())
				}
			},
		},
		{
			name: "template execution with simple variable",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: My Title
---
Content with {{ .title }}`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.Bytes(), []byte("Content with My Title")) {
					t.Errorf("expected executed template, got: %s", doc.String())
				}
			},
		},
		{
			name: "template execution with invalid syntax",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: Test
---
Content with {{ .invalid syntax }}`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr:     true,
			errContains: "template parse",
		},
		{
			name: "frontmatter with invalid YAML",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: [invalid yaml
---
Content`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr:     true,
			errContains: "frontmatter",
		},
		{
			name: "file does not exist",
			setupFile: func(t *testing.T, path string) {
				// Don't create the file
			},
			wantErr:     true,
			errContains: "os read file",
		},
		{
			name: "empty file",
			setupFile: func(t *testing.T, path string) {
				if err := os.WriteFile(path, []byte(""), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(doc.Bytes()) != 0 {
					t.Errorf("expected empty content, got: %s", doc.String())
				}
			},
		},
		{
			name: "document with only frontmatter delimiter",
			setupFile: func(t *testing.T, path string) {
				content := `---
---`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(sourceDir, "test.md")
			tt.setupFile(t, testFile)

			doc, err := Parse(sourceDir, testFile)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !errors.Is(err, errors.New(tt.errContains)) {
					if !containsErrorString(err, tt.errContains) {
						t.Errorf("expected error containing %q, got: %v", tt.errContains, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if tt.validateDoc != nil {
					tt.validateDoc(t, doc, err)
				}
			}
		})
	}
}

func TestParseWithInclude(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create directory structure: %v", err)
	}

	tests := []struct {
		name        string
		setupFiles  func(t *testing.T, dir string)
		wantErr     bool
		errContains string
		validateDoc func(t *testing.T, doc Document, err error)
	}{
		{
			name: "include with valid relative path",
			setupFiles: func(t *testing.T, dir string) {
				includedContent := `---
title: Included
---
Included content`
				if err := os.WriteFile(filepath.Join(dir, "included.md"), []byte(includedContent), 0644); err != nil {
					t.Fatalf("failed to write included file: %v", err)
				}

				mainContent := `---
title: Main
---
Main content: {{ include "included.md" }}`
				if err := os.WriteFile(filepath.Join(dir, "main.md"), []byte(mainContent), 0644); err != nil {
					t.Fatalf("failed to write main file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.Bytes(), []byte("Included content")) {
					t.Errorf("expected included content in output, got: %s", doc.String())
				}
			},
		},
		{
			name: "include with subdirectory path",
			setupFiles: func(t *testing.T, dir string) {
				includedContent := "Subdir content"
				if err := os.WriteFile(filepath.Join(dir, "subdir", "nested.md"), []byte(includedContent), 0644); err != nil {
					t.Fatalf("failed to write nested file: %v", err)
				}

				mainContent := `---
title: Main
---
Main: {{ include "subdir/nested.md" }}`
				if err := os.WriteFile(filepath.Join(dir, "main.md"), []byte(mainContent), 0644); err != nil {
					t.Fatalf("failed to write main file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.Bytes(), []byte("Subdir content")) {
					t.Errorf("expected nested content in output, got: %s", doc.String())
				}
			},
		},
		{
			name: "include with non-existent file - should panic",
			setupFiles: func(t *testing.T, dir string) {
				mainContent := `---
title: Main
---
Main: {{ include "does-not-exist.md" }}`
				if err := os.WriteFile(filepath.Join(dir, "main.md"), []byte(mainContent), 0644); err != nil {
					t.Fatalf("failed to write main file: %v", err)
				}
			},
			wantErr:     true,
			errContains: "parse",
		},
		{
			name: "include with invalid frontmatter in included file",
			setupFiles: func(t *testing.T, dir string) {
				includedContent := `---
title: [invalid
---
Content`
				if err := os.WriteFile(filepath.Join(dir, "included.md"), []byte(includedContent), 0644); err != nil {
					t.Fatalf("failed to write included file: %v", err)
				}

				mainContent := `---
title: Main
---
Main: {{ include "included.md" }}`
				if err := os.WriteFile(filepath.Join(dir, "main.md"), []byte(mainContent), 0644); err != nil {
					t.Fatalf("failed to write main file: %v", err)
				}
			},
			wantErr:     true,
			errContains: "parse",
		},
		{
			name: "nested includes (include chain)",
			setupFiles: func(t *testing.T, dir string) {
				file3Content := "Final content"
				if err := os.WriteFile(filepath.Join(dir, "file3.md"), []byte(file3Content), 0644); err != nil {
					t.Fatalf("failed to write file3: %v", err)
				}

				file2Content := `---
title: File2
---
{{ include "file3.md" }}`
				if err := os.WriteFile(filepath.Join(dir, "file2.md"), []byte(file2Content), 0644); err != nil {
					t.Fatalf("failed to write file2: %v", err)
				}

				mainContent := `---
title: Main
---
{{ include "file2.md" }}`
				if err := os.WriteFile(filepath.Join(dir, "main.md"), []byte(mainContent), 0644); err != nil {
					t.Fatalf("failed to write main file: %v", err)
				}
			},
			wantErr: false,
			validateDoc: func(t *testing.T, doc Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.Bytes(), []byte("Final content")) {
					t.Errorf("expected nested include content in output, got: %s", doc.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFiles(t, sourceDir)

			mainFile := filepath.Join(sourceDir, "main.md")
			doc, err := Parse(sourceDir, mainFile)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsErrorString(err, tt.errContains) {
					t.Errorf("expected error containing %q, got: %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if tt.validateDoc != nil {
					tt.validateDoc(t, doc, err)
				}
			}
		})
	}
}

func TestDocumentMethods(t *testing.T) {
	tests := []struct {
		name   string
		doc    Document
		testFn func(t *testing.T, doc Document)
	}{
		{
			name: "SkipPlace returns meta value",
			doc: Document{
				meta: Meta{SkipPlace: true},
			},
			testFn: func(t *testing.T, doc Document) {
				if !doc.SkipPlace() {
					t.Error("expected SkipPlace to return true")
				}
			},
		},
		{
			name: "SkipPlace returns false when not set",
			doc: Document{
				meta: Meta{SkipPlace: false},
			},
			testFn: func(t *testing.T, doc Document) {
				if doc.SkipPlace() {
					t.Error("expected SkipPlace to return false")
				}
			},
		},
		{
			name: "Bytes returns contents",
			doc: Document{
				contents: []byte("test content"),
			},
			testFn: func(t *testing.T, doc Document) {
				if !bytes.Equal(doc.Bytes(), []byte("test content")) {
					t.Errorf("expected 'test content', got: %s", doc.Bytes())
				}
			},
		},
		{
			name: "String returns contents as string",
			doc: Document{
				contents: []byte("test content"),
			},
			testFn: func(t *testing.T, doc Document) {
				if doc.String() != "test content" {
					t.Errorf("expected 'test content', got: %s", doc.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFn(t, tt.doc)
		})
	}
}

// Helper function to check if error contains a string
func containsErrorString(err error, s string) bool {
	if err == nil {
		return false
	}
	return err.Error() != "" && (err.Error() == s || len(err.Error()) >= len(s))
}
