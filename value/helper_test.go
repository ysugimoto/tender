package value

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseIdent(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []Field
	}{
		{
			name:  "dot access",
			input: "foo.bar",
			expect: []Field{
				{name: "foo", syntax: none},
				{name: "bar", syntax: dot},
			},
		},
		{
			name:  "array index access",
			input: "foo[0]",
			expect: []Field{
				{name: "foo", syntax: none},
				{name: "0", syntax: sliceBracket},
			},
		},
		{
			name:  "object key access",
			input: `foo["bar"]`,
			expect: []Field{
				{name: "foo", syntax: none},
				{name: "bar", syntax: mapBracket},
			},
		},
		{
			name:  "combinated access",
			input: `foo.bar["baz"][0]`,
			expect: []Field{
				{name: "foo", syntax: none},
				{name: "bar", syntax: dot},
				{name: "baz", syntax: mapBracket},
				{name: "0", syntax: sliceBracket},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first, fields := parseFields(tt.input)
			parsed := append([]Field{first}, fields...)

			if diff := cmp.Diff(tt.expect, parsed, cmpopts.EquateComparable(Field{})); diff != "" {
				t.Errorf("Parsed ident unmatch, diff=%s", diff)
			}
		})
	}
}

func TestIsStruct(t *testing.T) {
	type V struct {
		Name string
	}

	tests := []struct {
		input  any
		expect bool
	}{
		{
			input:  V{Name: "foo"},
			expect: true,
		},
		{
			input:  &V{Name: "foo"},
			expect: true,
		},
	}

	for _, tt := range tests {
		if diff := cmp.Diff(tt.expect, IsStruct(reflect.ValueOf(tt.input))); diff != "" {
			t.Errorf("isStruct unmatch, diff=%s", diff)
		}
	}

}
