package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/lexer"
	"github.com/ysugimoto/tender/token"
)

func ref(s string) *string {
	return &s
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
					Token: token.Token{Literal: "for"},
					Iterator: &ast.Ident{
						Token: token.Token{Literal: "list"},
						Value: "list",
					},
					Arg1: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
					Block: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndFor{
						Token: token.Token{Literal: "endfor"},
					},
				},
			},
		},
		{
			name:  "With trimming",
			input: "%{~for v in list~}foo%{~endfor~}",
			expect: []ast.Node{
				&ast.For{
					Token: token.Token{
						Literal:   "for",
						LeftTrim:  true,
						RightTrim: true,
					},
					Iterator: &ast.Ident{
						Token: token.Token{Literal: "list"},
						Value: "list",
					},
					Arg1: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
					Block: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndFor{
						Token: token.Token{
							Literal:   "endfor",
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
						Literal:   "for",
						LeftTrim:  true,
						RightTrim: true,
					},
					Iterator: &ast.Ident{
						Token: token.Token{Literal: "list"},
						Value: "list",
					},
					Arg1: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
					Arg2: &ast.Ident{
						Token: token.Token{Literal: "w"},
						Value: "w",
					},
					Block: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndFor{
						Token: token.Token{
							Literal:   "endfor",
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
					Token: token.Token{
						Literal: "for",
					},
					Iterator: &ast.Ident{
						Token: token.Token{Literal: "list"},
						Value: "list",
					},
					Arg1: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
					Block: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "\n"},
						},
						&ast.For{
							Token: token.Token{
								Literal:   "for",
								LeftTrim:  true,
								RightTrim: true,
							},
							Iterator: &ast.Ident{
								Token: token.Token{Literal: "v"},
								Value: "v",
							},
							Arg1: &ast.Ident{
								Token: token.Token{Literal: "w"},
								Value: "w",
							},
							Block: []ast.Node{
								&ast.Interporation{
									Token: token.Token{Literal: "w"},
									Value: &ast.Ident{
										Token: token.Token{Literal: "w"},
										Value: "w",
									},
								},
							},
							End: &ast.EndFor{
								Token: token.Token{
									Literal:   "endfor",
									LeftTrim:  true,
									RightTrim: true,
								},
							},
						},
						&ast.Literal{
							Token: token.Token{Literal: "\n"},
						},
					},
					End: &ast.EndFor{
						Token: token.Token{
							Literal: "endfor",
						},
					},
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
					Token: token.Token{Literal: "if"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
				},
			},
		},
		{
			name:  "With trimming",
			input: `%{~ if v == "v" ~}foo%{~endif ~}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{
						Literal:   "if",
						LeftTrim:  true,
						RightTrim: true,
					},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndIf{
						Token: token.Token{
							Literal:   "endif",
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
					Token: token.Token{Literal: "if"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					Alternative: &ast.Else{
						Token: token.Token{Literal: "else"},
						Consequence: []ast.Node{
							&ast.Literal{
								Token: token.Token{Literal: "bar"},
							},
						},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
				},
			},
		},
		{
			name:  "if-elseif-else",
			input: `%{ if v == "v" }foo%{elseif v != "w"}bar%{else}%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{Literal: "if"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
					},
					Another: []*ast.ElseIf{
						{
							Token: token.Token{Literal: "elseif"},
							Condition: &ast.InfixExpression{
								Token: token.Token{Literal: "!="},
								Left: &ast.Ident{
									Token: token.Token{Literal: "v"},
									Value: "v",
								},
								Operator: "!=",
								Right: &ast.String{
									Token: token.Token{Literal: "w"},
									Value: "w",
								},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Token: token.Token{Literal: "bar"},
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					Alternative: &ast.Else{
						Token:       token.Token{Literal: "else"},
						Consequence: []ast.Node{},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
				},
			},
		},
		{
			name:  "separated else if",
			input: `%{ if v == "v" }foo%{else if v != "w"}bar%{else}%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{Literal: "if"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
					},
					Another: []*ast.ElseIf{
						{
							Token: token.Token{Literal: "else if"},
							Condition: &ast.InfixExpression{
								Token: token.Token{Literal: "!="},
								Left: &ast.Ident{
									Token: token.Token{Literal: "v"},
									Value: "v",
								},
								Operator: "!=",
								Right: &ast.String{
									Value: "w",
									Token: token.Token{Literal: "w"},
								},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Token: token.Token{Literal: "bar"},
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					Alternative: &ast.Else{
						Token:       token.Token{Literal: "else"},
						Consequence: []ast.Node{},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
				},
			},
		},
		{
			name:  "if-elseif-elseif-else",
			input: `%{ if v == "v" }foo%{elseif v != "w"}bar%{ elseif v > 0 }baz%{else}%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{Literal: "if"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
					},
					Another: []*ast.ElseIf{
						{
							Token: token.Token{Literal: "elseif"},
							Condition: &ast.InfixExpression{
								Token: token.Token{Literal: "!="},
								Left: &ast.Ident{
									Token: token.Token{Literal: "v"},
									Value: "v",
								},
								Operator: "!=",
								Right: &ast.String{
									Token: token.Token{Literal: "w"},
									Value: "w",
								},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Token: token.Token{Literal: "bar"},
								},
							},
						},
						{
							Token: token.Token{Literal: "elseif"},
							Condition: &ast.InfixExpression{
								Token: token.Token{Literal: ">"},
								Left: &ast.Ident{
									Token: token.Token{Literal: "v"},
									Value: "v",
								},
								Operator: ">",
								Right: &ast.Int{
									Token: token.Token{Literal: "0"},
									Value: 0,
								},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Token: token.Token{Literal: "baz"},
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					Alternative: &ast.Else{
						Token:       token.Token{Literal: "else"},
						Consequence: []ast.Node{},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
				},
			},
		},
		{
			name:  "single condition expression",
			input: `%{ if v }foo%{elseif !v}bar%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{Literal: "if"},
					Condition: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
					Another: []*ast.ElseIf{
						{
							Token: token.Token{Literal: "elseif"},
							Condition: &ast.PrefixExpression{
								Token:    token.Token{Literal: "!"},
								Operator: "!",
								Right: &ast.Ident{
									Token: token.Token{Literal: "v"},
									Value: "v",
								},
							},
							Consequence: []ast.Node{
								&ast.Literal{
									Token: token.Token{Literal: "bar"},
								},
							},
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
				},
			},
		},
		{
			name:  "complicated condition",
			input: `%{ if (v == "v" && w == "w") || x > -1 || y < 0 || z >= 1 }foo%{endif}`,
			expect: []ast.Node{
				&ast.If{
					Token: token.Token{Literal: "if"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Literal: "||"},
						Left: &ast.InfixExpression{
							Token: token.Token{Literal: "||"},
							Left: &ast.InfixExpression{
								Token: token.Token{Literal: "||"},
								Left: &ast.GroupedExpression{
									Token: token.Token{Literal: "("},
									Right: &ast.InfixExpression{
										Token: token.Token{Literal: "&&"},
										Left: &ast.InfixExpression{
											Token: token.Token{Literal: "=="},
											Left: &ast.Ident{
												Token: token.Token{Literal: "v"},
												Value: "v",
											},
											Operator: "==",
											Right: &ast.String{
												Token: token.Token{Literal: "v"},
												Value: "v",
											},
										},
										Operator: "&&",
										Right: &ast.InfixExpression{
											Token: token.Token{Literal: "=="},
											Left: &ast.Ident{
												Token: token.Token{Literal: "w"},
												Value: "w",
											},
											Operator: "==",
											Right: &ast.String{
												Token: token.Token{Literal: "w"},
												Value: "w",
											},
										},
									},
								},
								Operator: "||",
								Right: &ast.InfixExpression{
									Token: token.Token{Literal: ">"},
									Left: &ast.Ident{
										Token: token.Token{Literal: "x"},
										Value: "x",
									},
									Operator: ">",
									Right: &ast.PrefixExpression{
										Token:    token.Token{Literal: "-"},
										Operator: "-",
										Right: &ast.Int{
											Token: token.Token{Literal: "1"},
											Value: 1,
										},
									},
								},
							},
							Operator: "||",
							Right: &ast.InfixExpression{
								Token: token.Token{Literal: "<"},
								Left: &ast.Ident{
									Token: token.Token{Literal: "y"},
									Value: "y",
								},
								Operator: "<",
								Right: &ast.Int{
									Token: token.Token{Literal: "0"},
									Value: 0,
								},
							},
						},
						Operator: "||",
						Right: &ast.InfixExpression{
							Token: token.Token{Literal: ">="},
							Left: &ast.Ident{
								Token: token.Token{Literal: "z"},
								Value: "z",
							},
							Operator: ">=",
							Right: &ast.Int{
								Token: token.Token{Literal: "1"},
								Value: 1,
							},
						},
					},
					Another: []*ast.ElseIf{},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "foo"},
						},
					},
					End: &ast.EndIf{
						Token: token.Token{Literal: "endif"},
					},
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
