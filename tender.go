package tender

import (
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/ysugimoto/tender/lexer"
	"github.com/ysugimoto/tender/parser"
	"github.com/ysugimoto/tender/value"
)

// Shorthand type for map[string]any
type Variables map[string]any

// Template struct represents template renderer.
// This struct holds some state that global variables which is provided via caller,
// and assigning local variables.
type Template struct {
	reader io.Reader
	global value.Value
	locals []value.Value
}

// Create template pointer from io.Reader stream
func New(r io.Reader) *Template {
	return &Template{
		reader: r,
		global: value.Value{},
		locals: []value.Value{},
	}
}

// Create template pointer from template string
func NewFromString(tmpl string) *Template {
	return New(strings.NewReader(tmpl))
}

// Assign variables to this template rendering
func (t *Template) With(variables Variables) *Template {
	for k, v := range variables {
		t.global[k] = reflect.ValueOf(v)
	}
	return t
}

// Render the template string with prvided variables.
// This method may return erorr as second return value,
// you can handle the error if your template has syntax, typing problem
func (t *Template) Render() (string, error) {
	nodes, err := parser.New(lexer.New(t.reader)).Parse()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return t.render(nodes)
}
