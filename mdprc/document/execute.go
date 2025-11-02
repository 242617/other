package document

import (
	"bytes"

	"github.com/pkg/errors"
)

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
