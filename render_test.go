package tender

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/lexer"
)

func TestInterporation(t *testing.T) {
	os.Setenv("FOO_BAR", "baz")

	tests := []struct {
		name    string
		input   string
		value   any
		expect  string
		isError bool
	}{
		{name: "env", input: "FOO_BAR", value: int(1), expect: "baz"},
		{name: "no env", input: "FOO_BAR_BAZ", value: int(1), isError: true},
		{name: "int", input: "v", value: int(1), expect: "1"},
		{name: "int8", input: "v", value: int8(1), expect: "1"},
		{name: "int16", input: "v", value: int16(1), expect: "1"},
		{name: "int32", input: "v", value: int32(1), expect: "1"},
		{name: "int64", input: "v", value: int64(1), expect: "1"},
		{name: "uint", input: "v", value: uint(1), expect: "1"},
		{name: "uint8", input: "v", value: uint8(1), expect: "1"},
		{name: "uint16", input: "v", value: uint16(1), expect: "1"},
		{name: "uint32", input: "v", value: uint32(1), expect: "1"},
		{name: "uint64", input: "v", value: uint64(1), expect: "1"},
		{name: "float32", input: "v", value: float32(1.1), expect: "1.1"},
		{name: "float64", input: "v", value: float64(1.1), expect: "1.1"},
		{name: "bool(true)", input: "v", value: true, expect: "true"},
		{name: "bool(false)", input: "v", value: false, expect: "false"},
		{name: "string", input: "v", value: "foo", expect: "foo"},
		{name: "slice", input: "v", value: []int{1, 2, 3}, expect: "[1, 2, 3]"},
		{name: "map", input: "v", value: map[string]int{"key": 1}, expect: "{key: 1}"},
		{name: "struct", input: "v", value: &ast.Literal{Value: "foo"}, expect: "{Value: foo}"},
		{name: "struct(unexported fields)", input: "v", value: &lexer.Lexer{}, expect: "{}"},
		{
			name:  "combinated",
			input: "v[0].foo.bar[0].Value",
			value: []map[string]any{
				{
					"foo": map[string][]*ast.Literal{
						"bar": {
							{Value: "lorem"},
							{Value: "ipsum"},
							{Value: "dolor"},
						},
					},
				},
			},
			expect: "lorem",
		},
		{
			name:  "combinated(unexported field)",
			input: "v[0].foo.bar[0].line",
			value: []map[string]any{
				{
					"foo": map[string][]*lexer.Lexer{
						"bar": {
							{},
							{},
							{},
						},
					},
				},
			},
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := "Value is ${" + tt.input + "}"
			vars := Variables{"v": tt.value}
			rendered, err := NewFromString(tmpl).With(vars).Render()
			if tt.isError {
				if err == nil {
					t.Errorf("Expects error, but got-nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected render error\n %+v", err)
				return
			}
			if diff := cmp.Diff("Value is "+tt.expect, rendered); diff != "" {
				t.Errorf("Rendered string mismatch, diff=%s", diff)
				return
			}

		})
	}
}

func TestForSliceLoop(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		value   any
		expect  string
		isError bool
	}{
		{
			name: "single variable",
			input: `
%{ for i in some_list }
Index is ${i}
%{ endfor }`,
			value: []string{"a", "b", "c"},
			expect: `

Index is 0

Index is 1

Index is 2
`,
		},
		{
			name: "double variable",
			input: `
%{ for i, v in some_list }
Index is ${i}, value is ${v}
%{ endfor }`,
			value: []string{"a", "b", "c"},
			expect: `

Index is 0, value is a

Index is 1, value is b

Index is 2, value is c
`,
		},
		{
			name: "trimming render",
			input: `
%{~ for i, v in some_list ~}
Index is ${i}, value is ${v}
%{~ endfor ~}`,
			value:  []string{"a", "b", "c"},
			expect: `Index is 0, value is aIndex is 1, value is bIndex is 2, value is c`,
		},
		{
			name: "not iterable",
			input: `
%{~ for i, v in some_list ~}
Index is ${i}, value is ${v}
%{~ endfor ~}`,
			value:   "foo",
			isError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := Variables{"some_list": tt.value}
			rendered, err := NewFromString(tt.input).With(vars).Render()
			if tt.isError {
				if err == nil {
					t.Errorf("Expects error, but got-nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected render error\n %+v", err)
				return
			}
			if diff := cmp.Diff(tt.expect, rendered); diff != "" {
				t.Errorf("Rendered string mismatch, diff=%s", diff)
				return
			}

		})
	}
}

func TestForMapLoop(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		value   any
		expect  string
		isError bool
	}{
		{
			name: "single variable",
			input: `
%{ for k in some_map }
Key is ${k}
%{ endfor }`,
			value: map[string]any{
				"foo":   "bar",
				"hoge":  "huga",
				"lorem": "ipsum",
			},
			expect: `

Key is foo

Key is hoge

Key is lorem
`,
		},
		{
			name: "double variable",
			input: `
%{ for k, v in some_map }
Key is ${k}, value is ${v}
%{ endfor }`,
			value: map[string]any{
				"foo":   "bar",
				"hoge":  "huga",
				"lorem": "ipsum",
			},
			expect: `

Key is foo, value is bar

Key is hoge, value is huga

Key is lorem, value is ipsum
`,
		},
		{
			name: "trimming render",
			input: `
%{~ for i, v in some_map ~}
Key is ${i}, value is ${v}
%{~ endfor ~}`,
			value: map[string]any{
				"foo":   "bar",
				"hoge":  "huga",
				"lorem": "ipsum",
			},
			expect: `Key is foo, value is barKey is hoge, value is hugaKey is lorem, value is ipsum`,
		},
		{
			name: "not iterable",
			input: `
%{~ for i, v in some_map ~}
Key is ${i}, value is ${v}
%{~ endfor ~}`,
			value:   "foo",
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := Variables{"some_map": tt.value}
			rendered, err := NewFromString(tt.input).With(vars).Render()
			if tt.isError {
				if err == nil {
					t.Errorf("Expects error, but got-nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected render error\n %+v", err)
				return
			}
			if diff := cmp.Diff(tt.expect, rendered); diff != "" {
				t.Errorf("Rendered string mismatch, diff=%s", diff)
				return
			}

		})
	}
}

func TestIfControl(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		value   Variables
		expect  string
		isError bool
	}{
		{
			name: "basic if (true)",
			input: `
%{ if v == "v" }
v is "v"
%{ endif }`,
			value: Variables{
				"v": "v",
			},
			expect: `

v is "v"
`,
		},
		{
			name: "trimming basic if (true)",
			input: `
%{~ if v == "v" ~}
v is "v"
%{~ endif ~}`,
			value: Variables{
				"v": "v",
			},
			expect: `v is "v"`,
		},
		{
			name: "basic if (false)",
			input: `
%{ if v == "w" }
v is "w"
%{ endif }`,
			value: Variables{
				"v": "v",
			},
			expect: `
`,
		},
		{
			name: "trimming basic if (false)",
			input: `
%{~ if v == "w" ~}
v is "w"
%{~ endif ~}`,
			value: Variables{
				"v": "v",
			},
			expect: "",
		},
		{
			name: "basic if-elseif (false)",
			input: `
%{ if v == "w" }
v is "w"
%{ else if v == "v" }
v is "v"
%{ endif }`,
			value: Variables{
				"v": "v",
			},
			expect: `

v is "v"
`,
		},
		{
			name: "trimming basic if-elseif (true)",
			input: `
%{~ if v == "w" ~}
v is "w"
%{~ else if v == "v" ~}
v is "v"
%{~ endif ~}`,
			value: Variables{
				"v": "v",
			},
			expect: `v is "v"`,
		},
		{
			name: "basic if-elseif (true)",
			input: `
%{ if v == "w" }
v is "w"
%{ else if v == "v" }
v is "v"
%{ endif }`,
			value: Variables{
				"v": "w",
			},
			expect: `

v is "w"
`,
		},
		{
			name: "trimming basic if-elseif (true)",
			input: `
%{~ if v == "w" ~}
v is "w"
%{~ else if v == "v" ~}
v is "v"
%{~ endif ~}`,
			value: Variables{
				"v": "w",
			},
			expect: `v is "w"`,
		},
		{
			name: "basic if-elseif-elseif (third)",
			input: `
%{ if v == "w" }
v is "w"
%{ else if v == "v" }
v is "v"
%{ else if v == "z" }
v is "z"
%{ endif }`,
			value: Variables{
				"v": "z",
			},
			expect: `

v is "z"
`,
		},
		{
			name: "trimming basic if-elseif-elseif (third)",
			input: `
%{~ if v == "w" ~}
v is "w"
%{~ else if v == "v" ~}
v is "v"
%{~ else if v == "z" ~}
v is "z"
%{~ endif ~}`,
			value: Variables{
				"v": "z",
			},
			expect: `v is "z"`,
		},
		{
			name: "basic if-elseif-else",
			input: `
%{ if v == "w" }
v is "w"
%{ else if v == "v" }
v is "v"
%{ else }
v is other
%{ endif }`,
			value: Variables{
				"v": "foo",
			},
			expect: `

v is other
`,
		},
		{
			name: "trimming basic if-elseif-else",
			input: `
%{~ if v == "w" ~}
v is "w"
%{~ else if v == "v" ~}
v is "v"
%{~ else ~}
v is other
%{~ endif ~}`,
			value: Variables{
				"v": "foo",
			},
			expect: `v is other`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered, err := NewFromString(tt.input).With(tt.value).Render()
			if tt.isError {
				if err == nil {
					t.Errorf("Expects error, but got-nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected render error\n %+v", err)
				return
			}
			if diff := cmp.Diff(tt.expect, rendered); diff != "" {
				t.Errorf("Rendered string mismatch, diff=%s", diff)
				return
			}

		})
	}
}

func TestRender(t *testing.T) {
	input := `This is template spec.

%{ for v in some_list ~}
inside loop, ${v} is variable interporation.
%{ endfor }

%{~ for i, v in some_map }
Also can loop for map.
%{ endfor }

%{ if v == "v" }
if expression is also supported. Interporation is ${v}.
%{ elseif v == "w" }
Render when v is "w".
%{ else }
else also.
%{ endif }

%{ if (v == "v" && w == "w") || v != "x" }complicated condition%{endif}

%%{ should recognize escaped string
$$ is escaped dollar character

That's all, very simplified!
`

	vars := Variables{
		"some_list": []string{"foo"},
		"some_map":  map[string]any{},
		"v":         "v",
		"w":         "w",
	}

	expect := `This is template spec.

inside loop, 0 is variable interporation.



if expression is also supported. Interporation is v.


complicated condition

%{ should recognize escaped string
$ is escaped dollar character

That's all, very simplified!
`

	output, err := NewFromString(input).With(vars).Render()
	if err != nil {
		t.Errorf("Unexpected render error\n %+v", err)
		return
	}
	if diff := cmp.Diff(expect, output); diff != "" {
		t.Errorf("Rendered string mismatch, diff=%s", diff)
		return
	}
}
