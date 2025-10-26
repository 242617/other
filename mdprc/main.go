package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func init() { log.SetFlags(log.Lshortfile) }
func main() {
	source, target := flag.String("source", "", "Source path (directory)"), flag.String("target", "", "Target path (directory)")
	flag.Parse()

	if *source == "" {
		log.Fatal(errors.New("empty source"))
	}
	if *target == "" {
		log.Fatal(errors.New("empty target"))
	}

	src, trg := filepath.Clean(*source), filepath.Clean(*target)

	for _, path := range []string{src, trg} {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(errors.Wrap(err, "os open"))
		}
		defer f.Close()
		if info, err := f.Stat(); err != nil {
			log.Fatal(errors.Wrap(err, "f stat"))
		} else if !info.IsDir() {
			log.Fatal(fmt.Errorf("path %q is not directory, as needed", path))
		}
	}

	if err := filepath.WalkDir(src, walk(src, trg)); err != nil {
		log.Fatal(errors.Wrap(err, "filepath walk dir"))
	}
}

func walk(source, target string) fs.WalkDirFunc {
	return func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath, err := filepath.Rel(source, path)
		if err != nil {
			return errors.Wrap(err, "filepath rel")
		}
		targetPath = filepath.Join(target, targetPath)

		info, err := dirEntry.Info()
		if err != nil {
			return errors.Wrap(err, "dir entry info")
		}

		if dirEntry.Type().IsDir() {
			if err := os.MkdirAll(targetPath, info.Mode().Perm()); err != nil {
				return errors.Wrap(err, "mkdir all")
			}
			return nil
		}

		if err := place(path, source, targetPath, info.Mode().Perm(), filepath.Ext(dirEntry.Name()) == ".md"); err != nil {
			return errors.Wrap(err, "place")
		}

		return nil
	}
}

func place(from, source, to string, perm fs.FileMode, process bool) error {
	if !process {
		b, err := os.ReadFile(from)
		if err != nil {
			return errors.Wrap(err, "os read file")
		}
		if err := os.WriteFile(to, b, perm); err != nil {
			return errors.Wrap(err, "os write file")
		}
		return nil
	}

	document, err := Parse(source, from)
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	if document.SkipPlace() {
		return nil
	}

	if err := os.WriteFile(to, document.Bytes(), perm); err != nil {
		return errors.Wrap(err, "os write file")
	}
	return nil
}

/// Document

func Parse(root, path string) (Document, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Document{}, errors.Wrap(err, "os read file")
	}

	document := Document{
		contents: b,
		tpl: template.New(path).Funcs(
			template.FuncMap{
				"include": func(arg reflect.Value, _ ...reflect.Value) reflect.Value {
					regarding := filepath.Dir(filepath.Join(path))
					include := filepath.Join(regarding, arg.String())

					document, err := Parse(root, include)
					if err != nil {
						panic(errors.Wrap(err, "parse"))
					}
					return reflect.ValueOf(document.String())
				},
			},
		),
	}

	var meta Meta
	if _, err := frontmatter.MustParse(bytes.NewReader(document.contents), &meta); errors.Is(err, frontmatter.ErrNotFound) {
		return document, nil
	} else if err != nil {
		return Document{}, errors.Wrap(err, "frontmatter must parse")
	}
	document.meta = meta

	if err := document.apply(); err != nil {
		return Document{}, errors.Wrap(err, "document apply")
	}

	if err := document.execute(); err != nil {
		return Document{}, errors.Wrap(err, "document execute")
	}

	return document, nil
}

func (d *Document) apply() error {
	rest, err := frontmatter.Parse(bytes.NewReader(d.contents), &d.properties)
	if err != nil {
		return errors.Wrap(err, "frontmatter parse")
	}

	if d.meta.RemoveProperties {
		d.contents = rest
		return nil
	}

	props, ok := bytes.CutSuffix(d.contents, rest)
	if !ok {
		return errors.New("unexpected behaviour: bytes cut suffix")
	}
	unwrapped := unwrap(string(props))

	res := unwrapped
	for _, key := range []string{
		`.[mdprc:skip_execute]`,
		`.[mdprc:skip_place]`,
		`.[mdprc:remove_properties]`,
	} {
		res, err = remove(res, key)
		if err != nil && !errors.Is(err, ErrKeyNotFound) {
			return errors.Wrap(err, "unexpected behaviour: remove")
		}
	}
	res = strings.TrimSpace(res)
	if res == "{}" {
		d.contents = rest
		return nil
	}
	wrapped := wrap(res)

	d.contents = append([]byte(wrapped), rest...)

	return nil
}

func (d *Document) execute() error {
	if d.meta.SkipExecute {
		return nil
	}

	tpl, err := d.tpl.Parse(string(d.contents))
	if err != nil {
		return errors.Wrap(err, "template parse")
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, d.properties)
	if err != nil {
		return errors.Wrap(err, "tpl execute")
	}

	d.contents = buf.Bytes()

	return nil
}

type (
	Document struct {
		meta       Meta
		properties map[string]any
		contents   []byte
		tpl        *template.Template
	}
	Meta struct {
		SkipExecute      bool `yaml:"[mdprc:skip_execute]"`
		SkipPlace        bool `yaml:"[mdprc:skip_place]"`
		RemoveProperties bool `yaml:"[mdprc:remove_properties]"`
	}
)

func (d Document) SkipPlace() bool { return d.meta.SkipPlace }

func (d Document) Bytes() []byte  { return d.contents }
func (d Document) String() string { return string(d.Bytes()) }

/// Utils

const (
	trailing = "---\n"
	leading  = "\n---\n"
)

func unwrap(str string) string {
	str = strings.TrimPrefix(str, trailing)
	str = strings.TrimSuffix(str, leading)
	return str
}
func wrap(str string) string {
	return trailing + str + leading
}

func remove(doc string, path string) (string, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &root); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	trimmedPath := strings.TrimPrefix(path, ".")
	if trimmedPath == "" {
		return "", fmt.Errorf("empty path")
	}

	keys := strings.Split(trimmedPath, ".")
	if err := removeNestedKey(&root, keys); err != nil {
		return "", err
	}

	out, err := yaml.Marshal(&root)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(out), nil
}

var ErrKeyNotFound = errors.New("key not found")

func removeNestedKey(node *yaml.Node, keys []string) error {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return fmt.Errorf("empty document")
		}
		return removeNestedKey(node.Content[0], keys)
	}

	// if node.Kind != yaml.MappingNode {
	// 	return fmt.Errorf("expected mapping node")
	// }

	if len(keys) == 1 {
		// Последний ключ - удаляем его
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == keys[0] {
				node.Content = append(node.Content[:i], node.Content[i+2:]...)
				return nil
			}
		}
		return errors.Wrap(ErrKeyNotFound, fmt.Sprintf("key %q not found", keys[0]))
	}

	// Ищем следующий уровень вложенности
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == keys[0] {
			return removeNestedKey(node.Content[i+1], keys[1:])
		}
	}
	return errors.Wrap(ErrKeyNotFound, fmt.Sprintf("key %q not found", keys[0]))
}
