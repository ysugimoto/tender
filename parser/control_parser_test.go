package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ysugimoto/tiny-template/ast"
	"github.com/ysugimoto/tiny-template/lexer"
	"github.com/ysugimoto/tiny-template/token"
)

var ignores = []cmp.Option{
	cmpopts.IgnoreFields(token.Token{}, "Line", "Position", "Type", "Literal"),
}

func TestForControl(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expect  []ast.Node
		isError bool
	}{
		{
			name:  "Basic parsing",
			input: "%{for v in list }foo%{endfor}",
			expect: []ast.Node{
				&ast.For{
					Interator: &ast.Ident{Value: "list"},
					Arg1:      &ast.Ident{Value: "v"},
					Block: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndFor{},
				},
			},
		},
		{
			name:  "With trimming",
			input: "%{~for v in list~}foo%{~endfor~}",
			expect: []ast.Node{
				&ast.For{
					Token: token.Token{
						LeftTrim:  true,
						RightTrim: true,
					},
					Interator: &ast.Ident{Value: "list"},
					Arg1:      &ast.Ident{Value: "v"},
					Block: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndFor{
						Token: token.Token{
							LeftTrim:  true,
							RightTrim: true,
						},
					},
				},
			},
		},
		{
			name:  "double arguments",
			input: "%{~for v, w in list~}foo%{~endfor~}",
			expect: []ast.Node{
				&ast.For{
					Token: token.Token{
						LeftTrim:  true,
						RightTrim: true,
					},
					Interator: &ast.Ident{Value: "list"},
					Arg1:      &ast.Ident{Value: "v"},
					Arg2:      &ast.Ident{Value: "w"},
					Block: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndFor{
						Token: token.Token{
							LeftTrim:  true,
							RightTrim: true,
						},
					},
				},
			},
		},
		{
			name: "nested loop, interporation",
			input: `%{ for v in list }
%{~ for w in v ~}${w}%{~ endfor ~}
%{ endfor }`,
			expect: []ast.Node{
				&ast.For{
					Interator: &ast.Ident{Value: "list"},
					Arg1:      &ast.Ident{Value: "v"},
					Block: []ast.Node{
						&ast.Literal{Value: "\n"},
						&ast.For{
							Token: token.Token{
								LeftTrim:  true,
								RightTrim: true,
							},
							Interator: &ast.Ident{Value: "v"},
							Arg1:      &ast.Ident{Value: "w"},
							Block: []ast.Node{
								&ast.Interporation{
									Value: &ast.Ident{Value: "w"},
								},
							},
							End: &ast.EndFor{
								Token: token.Token{
									LeftTrim:  true,
									RightTrim: true,
								},
							},
						},
						&ast.Literal{Value: "\n"},
					},
					End: &ast.EndFor{},
				},
			},
		},
		{
			name:    "Invalid syntax - argument is not specified",
			input:   `%{ for in list }foo%{endfor}`,
			isError: true,
		},
		{
			name:    "Invalid syntax - \"in\" keyword is not specified",
			input:   `%{ for v list }foo%{endfor}`,
			isError: true,
		},
		{
			name:    "Invalid syntax - second argument is not specified",
			input:   `%{ for v, in list }foo%{endfor}`,
			isError: true,
		},
		{
			name:    "Invalid syntax - iterator is not specified",
			input:   `%{ for v in }foo%{endfor}`,
			isError: true,
		},
		{
			name:    "Invalid syntax - endfor control is not specified ",
			input:   `%{ for v in list }foo`,
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := New(lexer.NewFromString(tt.input)).Parse()
			if err != nil {
				if !tt.isError {
					t.Errorf("Unexpected error: %s", err)
					return
				}
				return
			}
			if tt.isError {
				t.Errorf("Expects error but got nil")
				return
			}
			if diff := cmp.Diff(tt.expect, parsed, ignores...); diff != "" {
				t.Errorf("Unmatch parsed result, diff=%s", diff)
			}
		})
	}

}
