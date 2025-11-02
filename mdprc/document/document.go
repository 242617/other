package document

import "text/template"

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
