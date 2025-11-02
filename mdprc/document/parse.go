package document

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
	"github.com/pkg/errors"
)

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

	// Remove mdprc properties
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
		changed, err := remove(res, key)
		if errors.Is(err, ErrKeyNotFound) {
			continue
		}
		if err != nil {
			return errors.Wrap(err, "unexpected behaviour: remove")
		}
		res = changed
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
