package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/lexer"
	"github.com/ysugimoto/tender/token"
)

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
					Iterator: &ast.Ident{Value: "list"},
					Arg1:     &ast.Ident{Value: "v"},
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
					Iterator: &ast.Ident{Value: "list"},
					Arg1:     &ast.Ident{Value: "v"},
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
					Iterator: &ast.Ident{Value: "list"},
					Arg1:     &ast.Ident{Value: "v"},
					Arg2:     &ast.Ident{Value: "w"},
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
					Iterator: &ast.Ident{Value: "list"},
					Arg1:     &ast.Ident{Value: "v"},
					Block: []ast.Node{
						&ast.Literal{Value: "\n"},
						&ast.For{
							Token: token.Token{
								LeftTrim:  true,
								RightTrim: true,
							},
							Iterator: &ast.Ident{Value: "v"},
							Arg1:     &ast.Ident{Value: "w"},
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
		{
			name:    "Invalid syntax - unexpected control ",
			input:   `%{ for v in list }foo%{endif}`,
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

func TestIfControl(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expect  []ast.Node
		isError bool
	}{
		{
			name:  "Basic if",
			input: `%{ if v == "v" }foo%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "v"},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:  "With trimming",
			input: `%{~ if v == "v" ~}foo%{~endif ~}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{
						LeftTrim:  true,
						RightTrim: true,
					},
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "v"},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndIf{
						Token: token.Token{
							LeftTrim:  true,
							RightTrim: true,
						},
					},
				},
			},
		},
		{
			name:  "if-else",
			input: `%{ if v == "v" }foo%{else}bar%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "v"},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					Alternative: &ast.Else{
						Consequence: []ast.Node{
							&ast.Literal{Value: "bar"},
						},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:  "if-elseif-else",
			input: `%{ if v == "v" }foo%{elseif v != "w"}bar%{else}%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "v"},
					},
					Another: []*ast.ElseIf{
						{
							Condition: &ast.InfixExpression{
								Left:     &ast.Ident{Value: "v"},
								Operator: "!=",
								Right:    &ast.String{Value: "w"},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Value: "bar",
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					Alternative: &ast.Else{
						Consequence: []ast.Node{},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:  "seprated else if",
			input: `%{ if v == "v" }foo%{else if v != "w"}bar%{else}%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "v"},
					},
					Another: []*ast.ElseIf{
						{
							Condition: &ast.InfixExpression{
								Left:     &ast.Ident{Value: "v"},
								Operator: "!=",
								Right:    &ast.String{Value: "w"},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Value: "bar",
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					Alternative: &ast.Else{
						Consequence: []ast.Node{},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:  "if-elseif-elseif-else",
			input: `%{ if v == "v" }foo%{elseif v != "w"}bar%{ elseif v > 0 }baz%{else}%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "v"},
					},
					Another: []*ast.ElseIf{
						{
							Condition: &ast.InfixExpression{
								Left:     &ast.Ident{Value: "v"},
								Operator: "!=",
								Right:    &ast.String{Value: "w"},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Value: "bar",
								},
							},
						},
						{
							Condition: &ast.InfixExpression{
								Left:     &ast.Ident{Value: "v"},
								Operator: ">",
								Right:    &ast.Int{Value: 0},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Value: "baz",
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					Alternative: &ast.Else{
						Consequence: []ast.Node{},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:  "single condition expression",
			input: `%{ if v }foo%{elseif !v}bar%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.Ident{Value: "v"},
					Another: []*ast.ElseIf{
						{
							Condition: &ast.PrefixExpression{
								Operator: "!",
								Right:    &ast.Ident{Value: "v"},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Value: "bar",
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:  "complicated condition",
			input: `%{ if (v == "v" && w == "w") || x > -1 || y < 0 || z >= 1 }foo%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Condition: &ast.InfixExpression{
						Left: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left: &ast.GroupedExpression{
									Right: &ast.InfixExpression{
										Left: &ast.InfixExpression{
											Left:     &ast.Ident{Value: "v"},
											Operator: "==",
											Right:    &ast.String{Value: "v"},
										},
										Operator: "&&",
										Right: &ast.InfixExpression{
											Left:     &ast.Ident{Value: "w"},
											Operator: "==",
											Right:    &ast.String{Value: "w"},
										},
									},
								},
								Operator: "||",
								Right: &ast.InfixExpression{
									Left:     &ast.Ident{Value: "x"},
									Operator: ">",
									Right: &ast.PrefixExpression{
										Operator: "-",
										Right:    &ast.Int{Value: 1},
									},
								},
							},
							Operator: "||",
							Right: &ast.InfixExpression{
								Left:     &ast.Ident{Value: "y"},
								Operator: "<",
								Right:    &ast.Int{Value: 0},
							},
						},
						Operator: "||",
						Right: &ast.InfixExpression{
							Left:     &ast.Ident{Value: "z"},
							Operator: ">=",
							Right:    &ast.Int{Value: 1},
						},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "foo",
						},
					},
					End: &ast.EndIf{},
				},
			},
		},
		{
			name:    "Invalid syntax - invalid expression",
			input:   `%{ if v == }foo%{endif}`,
			isError: true,
		},
		{
			name:    "Invalid syntax - unexpected elseif found after else",
			input:   `%{ if v == "v" }foo%{else}bar%{elseif v == "w"}baz%{endif}`,
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
