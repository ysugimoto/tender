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
	cmpopts.IgnoreFields(token.Token{}, "Line", "Position", "Type", "Literal"),
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
			Value: `This is template spec.

`,
		},
		&ast.For{
			Token: token.Token{
				RightTrim: true,
			},
			Iterator: &ast.Ident{Value: "some_list"},
			Arg1:     &ast.Ident{Value: "v"},
			Block: []ast.Node{
				&ast.Literal{
					Value: "\ninside loop, ",
				},
				&ast.Interporation{
					Value: &ast.Ident{Value: "v"},
				},
				&ast.Literal{
					Value: " is variable interporation.\n",
				},
			},
			End: &ast.EndFor{},
		},
		&ast.Literal{
			Value: "\n\n",
		},
		&ast.For{
			Token: token.Token{
				LeftTrim: true,
			},
			Iterator: &ast.Ident{Value: "some_map"},
			Arg1:     &ast.Ident{Value: "i"},
			Arg2:     &ast.Ident{Value: "v"},
			Block: []ast.Node{
				&ast.Literal{
					Value: "\nAlso can loop for map.\n",
				},
			},
			End: &ast.EndFor{},
		},
		&ast.Literal{
			Value: "\n\n",
		},
		&ast.If{
			Condition: &ast.InfixExpression{
				Left:     &ast.Ident{Value: "v"},
				Operator: "==",
				Right:    &ast.String{Value: "v"},
			},
			Another: []*ast.ElseIf{
				&ast.ElseIf{
					Condition: &ast.InfixExpression{
						Left:     &ast.Ident{Value: "v"},
						Operator: "==",
						Right:    &ast.String{Value: "w"},
					},
					Consequence: []ast.Node{
						&ast.Literal{
							Value: "\nRender when v is \"w\".\n",
						},
					},
				},
			},
			Consequence: []ast.Node{
				&ast.Literal{
					Value: "\nif expression is also supported. Interporation is ",
				},
				&ast.Interporation{
					Value: &ast.Ident{Value: "v"},
				},
				&ast.Literal{
					Value: ".\n",
				},
			},
			Alternative: &ast.Else{
				Consequence: []ast.Node{
					&ast.Literal{
						Value: "\nelse also.\n",
					},
				},
			},
			End: &ast.EndIf{},
		},
		&ast.Literal{
			Value: "\n\n",
		},
		&ast.If{
			Condition: &ast.InfixExpression{
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
					Left:     &ast.Ident{Value: "v"},
					Operator: "!=",
					Right:    &ast.String{Value: "x"},
				},
			},
			Another: []*ast.ElseIf{},
			Consequence: []ast.Node{
				&ast.Literal{
					Value: "complicated condition",
				},
			},
			End: &ast.EndIf{},
		},
		&ast.Literal{
			Value: `

%{ should recognize escaped string
$ is escaped dollar character

That's all, very simplified!
`,
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
