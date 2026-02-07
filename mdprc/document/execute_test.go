package document

import (
	"bytes"
	"testing"
	"text/template"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Document
		wantErr  bool
		validate func(t *testing.T, doc *Document, err error)
	}{
		{
			name: "execute simple template",
			doc: &Document{
				meta: Meta{SkipExecute: false},
				properties: map[string]any{
					"title": "My Title",
					"count": 42,
				},
				contents: []byte("Title: {{ .title }}, Count: {{ .count }}"),
				tpl:      nil, // Will be set in test
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.contents, []byte("Title: My Title")) {
					t.Errorf("expected template execution, got: %s", string(doc.contents))
				}
				if !bytes.Contains(doc.contents, []byte("Count: 42")) {
					t.Errorf("expected template execution, got: %s", string(doc.contents))
				}
			},
		},
		{
			name: "skip_execute enabled",
			doc: &Document{
				meta: Meta{SkipExecute: true},
				properties: map[string]any{
					"title": "My Title",
				},
				contents: []byte("Title: {{ .title }}"),
				tpl:      nil,
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Template should NOT be executed
				if !bytes.Contains(doc.contents, []byte("{{ .title }}")) {
					t.Errorf("expected unexecuted template, got: %s", string(doc.contents))
				}
			},
		},
		{
			name: "invalid template syntax",
			doc: &Document{
				meta:       Meta{SkipExecute: false},
				properties: map[string]any{},
				contents:   []byte("Invalid {{ .unclosed"),
				tpl:        nil,
			},
			wantErr: true,
		},
		{
			name: "template with unknown variable",
			doc: &Document{
				meta: Meta{SkipExecute: false},
				properties: map[string]any{
					"title": "My Title",
				},
				contents: []byte("Value: {{ .unknown }}"),
				tpl:      nil,
			},
			wantErr: false, // Go templates don't error on missing variables, they just output <no value>
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// The template should execute, but with <no value>
				if !bytes.Contains(doc.contents, []byte("<no value>")) {
					t.Errorf("expected <no value> for unknown variable, got: %s", string(doc.contents))
				}
			},
		},
		{
			name: "empty template",
			doc: &Document{
				meta:       Meta{SkipExecute: false},
				properties: map[string]any{},
				contents:   []byte(""),
				tpl:        nil,
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(doc.contents) != 0 {
					t.Errorf("expected empty content, got: %s", string(doc.contents))
				}
			},
		},
		{
			name: "template with multiple variables",
			doc: &Document{
				meta: Meta{SkipExecute: false},
				properties: map[string]any{
					"first":  "John",
					"last":   "Doe",
					"age":    30,
					"active": true,
				},
				contents: []byte("Name: {{ .first }} {{ .last }}, Age: {{ .age }}, Active: {{ .active }}"),
				tpl:      nil,
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				expected := "Name: John Doe, Age: 30, Active: true"
				if !bytes.Contains(doc.contents, []byte(expected)) {
					t.Errorf("expected '%s', got: %s", expected, string(doc.contents))
				}
			},
		},
		{
			name: "template with nested properties",
			doc: &Document{
				meta: Meta{SkipExecute: false},
				properties: map[string]any{
					"user": map[string]any{
						"name":  "Alice",
						"email": "alice@example.com",
					},
				},
				contents: []byte("User: {{ .user.name }}, Email: {{ .user.email }}"),
				tpl:      nil,
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				expected := "User: Alice, Email: alice@example.com"
				if !bytes.Contains(doc.contents, []byte(expected)) {
					t.Errorf("expected '%s', got: %s", expected, string(doc.contents))
				}
			},
		},
		{
			name: "template with conditional",
			doc: &Document{
				meta: Meta{SkipExecute: false},
				properties: map[string]any{
					"show": true,
				},
				contents: []byte("{{ if .show }}Visible{{ end }}"),
				tpl:      nil,
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Contains(doc.contents, []byte("Visible")) {
					t.Errorf("expected 'Visible', got: %s", string(doc.contents))
				}
			},
		},
		{
			name: "template with range loop",
			doc: &Document{
				meta: Meta{SkipExecute: false},
				properties: map[string]any{
					"items": []string{"one", "two", "three"},
				},
				contents: []byte("{{ range .items }}{{ . }} {{ end }}"),
				tpl:      nil,
			},
			wantErr: false,
			validate: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				expected := "one two three"
				if !bytes.Contains(doc.contents, []byte(expected)) {
					t.Errorf("expected '%s', got: %s", expected, string(doc.contents))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create template if not set
			if tt.doc.tpl == nil {
				tt.doc.tpl = newTemplate("test")
			}

			err := tt.doc.execute()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, tt.doc, err)
				}
			}
		})
	}
}

func TestExecutePanicScenarios(t *testing.T) {
	tests := []struct {
		name        string
		doc         *Document
		expectPanic bool
	}{
		{
			name: "nil template",
			doc: &Document{
				meta:       Meta{SkipExecute: false},
				properties: map[string]any{},
				contents:   []byte("test"),
				tpl:        nil, // This will cause nil pointer dereference
			},
			expectPanic: true,
		},
		{
			name:        "nil document",
			doc:         nil,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectPanic {
				t.Skip("not a panic scenario")
			}

			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic, but did not occur")
				}
			}()

			if tt.doc != nil && tt.doc.tpl == nil {
				tt.doc.tpl = nil // Explicitly set to nil
			}

			tt.doc.execute()
		})
	}
}

func TestExecuteFeatureGateCombinations(t *testing.T) {
	tests := []struct {
		name          string
		skipExecute   bool
		skipPlace     bool
		removeProps   bool
		properties    map[string]any
		contents      []byte
		expectedIn    []byte
		expectedNotIn []byte
	}{
		{
			name:        "all gates disabled - normal execution",
			skipExecute: false,
			skipPlace:   false,
			removeProps: false,
			properties: map[string]any{
				"title": "Test",
			},
			contents:      []byte("{{ .title }}"),
			expectedIn:    []byte("Test"),
			expectedNotIn: []byte("{{ .title }}"),
		},
		{
			name:        "only skip_execute enabled",
			skipExecute: true,
			skipPlace:   false,
			removeProps: false,
			properties: map[string]any{
				"title": "Test",
			},
			contents:      []byte("{{ .title }}"),
			expectedIn:    []byte("{{ .title }}"),
			expectedNotIn: []byte("Test"),
		},
		{
			name:        "only skip_place enabled",
			skipExecute: false,
			skipPlace:   true,
			removeProps: false,
			properties: map[string]any{
				"title": "Test",
			},
			contents:      []byte("{{ .title }}"),
			expectedIn:    []byte("Test"),
			expectedNotIn: []byte("{{ .title }}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{
				meta: Meta{
					SkipExecute:      tt.skipExecute,
					SkipPlace:        tt.skipPlace,
					RemoveProperties: tt.removeProps,
				},
				properties: tt.properties,
				contents:   tt.contents,
				tpl:        newTemplate("test"),
			}

			err := doc.execute()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(tt.expectedIn) > 0 && !bytes.Contains(doc.contents, tt.expectedIn) {
				t.Errorf("expected content to contain %q, got: %s", string(tt.expectedIn), string(doc.contents))
			}

			if len(tt.expectedNotIn) > 0 && bytes.Contains(doc.contents, tt.expectedNotIn) {
				t.Errorf("expected content not to contain %q, got: %s", string(tt.expectedNotIn), string(doc.contents))
			}
		})
	}
}

// Helper function to create a template
func newTemplate(name string) *template.Template {
	// Import text/template
	return template.New(name)
}

// Test execute with special characters and escaping
func TestExecuteSpecialCharacters(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]any
		contents   []byte
		validate   func(t *testing.T, result []byte)
	}{
		{
			name: "HTML escaping - text templates don't auto-escape",
			properties: map[string]any{
				"html": "<script>alert('test')</script>",
			},
			contents: []byte("{{ .html }}"),
			validate: func(t *testing.T, result []byte) {
				// text/template does NOT auto-escape (unlike html/template)
				if !bytes.Contains(result, []byte("<script>")) {
					t.Errorf("expected HTML to NOT be escaped in text/template, got: %s", string(result))
				}
			},
		},
		{
			name: "newlines in content",
			properties: map[string]any{
				"lines": "line1\nline2\nline3",
			},
			contents: []byte("{{ .lines }}"),
			validate: func(t *testing.T, result []byte) {
				if !bytes.Contains(result, []byte("line1")) {
					t.Errorf("expected line1, got: %s", string(result))
				}
			},
		},
		{
			name: "special characters in values",
			properties: map[string]any{
				"special": "value with @#$%^&*()",
			},
			contents: []byte("{{ .special }}"),
			validate: func(t *testing.T, result []byte) {
				if !bytes.Contains(result, []byte("@#$%^&*()")) {
					t.Errorf("expected special chars, got: %s", string(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{
				meta:       Meta{SkipExecute: false},
				properties: tt.properties,
				contents:   tt.contents,
				tpl:        newTemplate("test"),
			}

			err := doc.execute()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			tt.validate(t, doc.contents)
		})
	}
}
