package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/lexer"
	"github.com/ysugimoto/tender/token"
)

var ignores = []cmp.Option{
	cmpopts.IgnoreFields(token.Token{}, "Line", "Position", "Type"),
}

func TestParser(t *testing.T) {
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
	nodes, err := New(lexer.NewFromString(input)).Parse()
	if err != nil {
		t.Errorf("Unexpected error, %s", err)
		return
	}

	expect := []ast.Node{
		&ast.Literal{
			Token: token.Token{
				Literal: `This is template spec.

`,
			},
		},
		&ast.For{
			Token: token.Token{
				Literal:   "for",
				RightTrim: true,
			},
			Iterator: &ast.Ident{
				Token: token.Token{
					Literal: "some_list",
				},
				Value: "some_list",
			},
			Arg1: &ast.Ident{
				Token: token.Token{
					Literal: "v",
				},
				Value: "v",
			},
			Block: []ast.Node{
				&ast.Literal{
					Token: token.Token{Literal: "\ninside loop, "},
				},
				&ast.Interporation{
					Token: token.Token{
						Literal: "v",
					},
					Value: &ast.Ident{
						Token: token.Token{
							Literal: "v",
						},
						Value: "v",
					},
				},
				&ast.Literal{
					Token: token.Token{Literal: " is variable interporation.\n"},
				},
			},
			End: &ast.EndFor{
				Token: token.Token{
					Literal: "endfor",
				},
			},
		},
		&ast.Literal{
			Token: token.Token{Literal: "\n\n"},
		},
		&ast.For{
			Token: token.Token{
				Literal:  "for",
				LeftTrim: true,
			},
			Iterator: &ast.Ident{
				Token: token.Token{
					Literal: "some_map",
				},
				Value: "some_map",
			},
			Arg1: &ast.Ident{
				Token: token.Token{
					Literal: "i",
				},
				Value: "i",
			},
			Arg2: &ast.Ident{
				Token: token.Token{
					Literal: "v",
				},
				Value: "v",
			},
			Block: []ast.Node{
				&ast.Literal{
					Token: token.Token{Literal: "\nAlso can loop for map.\n"},
				},
			},
			End: &ast.EndFor{
				Token: token.Token{
					Literal: "endfor",
				},
			},
		},
		&ast.Literal{
			Token: token.Token{Literal: "\n\n"},
		},
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
						Token: token.Token{Literal: "=="},
						Left: &ast.Ident{
							Token: token.Token{Literal: "v"},
							Value: "v",
						},
						Operator: "==",
						Right: &ast.String{
							Token: token.Token{Literal: "w"},
							Value: "w",
						},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Token: token.Token{Literal: "\nRender when v is \"w\".\n"},
						},
					},
				},
			},
			Consequence: []ast.Node{
				&ast.Literal{
					Token: token.Token{Literal: "\nif expression is also supported. Interporation is "},
				},
				&ast.Interporation{
					Token: token.Token{Literal: "v"},
					Value: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
				},
				&ast.Literal{
					Token: token.Token{Literal: ".\n"},
				},
			},
			Alternative: &ast.Else{
				Token: token.Token{Literal: "else"},
				Consequence: []ast.Node{
					&ast.Literal{
						Token: token.Token{Literal: "\nelse also.\n"},
					},
				},
			},
			End: &ast.EndIf{
				Token: token.Token{Literal: "endif"},
			},
		},
		&ast.Literal{
			Token: token.Token{Literal: "\n\n"},
		},
		&ast.If{
			Token: token.Token{Literal: "if"},
			Condition: &ast.InfixExpression{
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
					Token: token.Token{Literal: "!="},
					Left: &ast.Ident{
						Token: token.Token{Literal: "v"},
						Value: "v",
					},
					Operator: "!=",
					Right: &ast.String{
						Token: token.Token{Literal: "x"},
						Value: "x",
					},
				},
			},
			Another: []*ast.ElseIf{},
			Consequence: []ast.Node{
				&ast.Literal{
					Token: token.Token{Literal: "complicated condition"},
				},
			},
			End: &ast.EndIf{
				Token: token.Token{Literal: "endif"},
			},
		},
		&ast.Literal{
			Token: token.Token{Literal: `

%{ should recognize escaped string
$ is escaped dollar character

That's all, very simplified!
`},
		},
	}

	if diff := cmp.Diff(expect, nodes, ignores...); diff != "" {
		t.Errorf("Parsed result mismatch, diff=%s", diff)
	}
}

func BenchmarkPar(b *testing.B) {
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
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		New(lexer.NewFromString(input)).Parse()
	}
}
