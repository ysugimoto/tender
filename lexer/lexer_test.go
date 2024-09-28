package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/tiny-template/token"
)

func TestLexer(t *testing.T) {
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
	expects := []token.Token{
		{Type: token.LITERAL, Literal: "This is template spec.\n\n", Line: 1, Position: 1},

		// // single argument loop
		{Type: token.CONTROL_START, Literal: "%{", Line: 3, Position: 1},
		{Type: token.FOR, Literal: "for", Line: 3, Position: 4},
		{Type: token.IDENT, Literal: "v", Line: 3, Position: 8},
		{Type: token.IN, Literal: "in", Line: 3, Position: 10},
		{Type: token.IDENT, Literal: "some_list", Line: 3, Position: 13},
		{Type: token.CONTROL_END, Literal: "~}", Line: 3, Position: 23, RightTrim: true},
		{Type: token.LITERAL, Literal: "\ninside loop, ", Line: 3, Position: 25},
		{Type: token.INTERPORATION, Literal: "v", Line: 4, Position: 16},
		{Type: token.LITERAL, Literal: " is variable interporation.\n", Line: 4, Position: 18},
		{Type: token.CONTROL_START, Literal: "%{", Line: 5, Position: 1},
		{Type: token.ENDFOR, Literal: "endfor", Line: 5, Position: 4},
		{Type: token.CONTROL_END, Literal: "}", Line: 5, Position: 11},

		{Type: token.LITERAL, Literal: "\n\n", Line: 5, Position: 12},

		// // double arguments loop
		{Type: token.CONTROL_START, Literal: "%{~", Line: 7, Position: 1, LeftTrim: true},
		{Type: token.FOR, Literal: "for", Line: 7, Position: 5},
		{Type: token.IDENT, Literal: "i", Line: 7, Position: 9},
		{Type: token.COMMA, Literal: ",", Line: 7, Position: 10},
		{Type: token.IDENT, Literal: "v", Line: 7, Position: 12},
		{Type: token.IN, Literal: "in", Line: 7, Position: 14},
		{Type: token.IDENT, Literal: "some_map", Line: 7, Position: 17},
		{Type: token.CONTROL_END, Literal: "}", Line: 7, Position: 26},
		{Type: token.LITERAL, Literal: "\nAlso can loop for map.\n", Line: 7, Position: 27},
		{Type: token.CONTROL_START, Literal: "%{", Line: 9, Position: 1},
		{Type: token.ENDFOR, Literal: "endfor", Line: 9, Position: 4},
		{Type: token.CONTROL_END, Literal: "}", Line: 9, Position: 11},

		{Type: token.LITERAL, Literal: "\n\n", Line: 9, Position: 12},

		// // // If - elseif - else
		{Type: token.CONTROL_START, Literal: "%{", Line: 11, Position: 1},
		{Type: token.IF, Literal: "if", Line: 11, Position: 4},
		{Type: token.IDENT, Literal: "v", Line: 11, Position: 7},
		{Type: token.EQUAL, Literal: "==", Line: 11, Position: 9},
		{Type: token.STRING, Literal: "v", Line: 11, Position: 12},
		{Type: token.CONTROL_END, Literal: "}", Line: 11, Position: 16},
		{Type: token.LITERAL, Literal: "\nif expression is also supported. Interporation is ", Line: 11, Position: 17},
		{Type: token.INTERPORATION, Literal: "v", Line: 12, Position: 53},
		{Type: token.LITERAL, Literal: ".\n", Line: 12, Position: 55},

		{Type: token.CONTROL_START, Literal: "%{", Line: 13, Position: 1},
		{Type: token.ELSEIF, Literal: "elseif", Line: 13, Position: 4},
		{Type: token.IDENT, Literal: "v", Line: 13, Position: 11},
		{Type: token.EQUAL, Literal: "==", Line: 13, Position: 13},
		{Type: token.STRING, Literal: "w", Line: 13, Position: 16},
		{Type: token.CONTROL_END, Literal: "}", Line: 13, Position: 20},
		{Type: token.LITERAL, Literal: "\nRender when v is \"w\".\n", Line: 13, Position: 21},

		{Type: token.CONTROL_START, Literal: "%{", Line: 15, Position: 1},
		{Type: token.ELSE, Literal: "else", Line: 15, Position: 4},
		{Type: token.CONTROL_END, Literal: "}", Line: 15, Position: 9},
		{Type: token.LITERAL, Literal: "\nelse also.\n", Line: 15, Position: 10},
		{Type: token.CONTROL_START, Literal: "%{", Line: 17, Position: 1},
		{Type: token.ENDIF, Literal: "endif", Line: 17, Position: 4},
		{Type: token.CONTROL_END, Literal: "}", Line: 17, Position: 10},

		{Type: token.LITERAL, Literal: "\n\n", Line: 17, Position: 11},

		// // Complicated if
		{Type: token.CONTROL_START, Literal: "%{", Line: 19, Position: 1},
		{Type: token.IF, Literal: "if", Line: 19, Position: 4},
		{Type: token.LEFT_PAREN, Literal: "(", Line: 19, Position: 7},
		{Type: token.IDENT, Literal: "v", Line: 19, Position: 8},
		{Type: token.EQUAL, Literal: "==", Line: 19, Position: 10},
		{Type: token.STRING, Literal: "v", Line: 19, Position: 13},
		{Type: token.AND, Literal: "&&", Line: 19, Position: 17},
		{Type: token.IDENT, Literal: "w", Line: 19, Position: 20},
		{Type: token.EQUAL, Literal: "==", Line: 19, Position: 22},
		{Type: token.STRING, Literal: "w", Line: 19, Position: 25},
		{Type: token.RIGHT_PAREN, Literal: ")", Line: 19, Position: 28},
		{Type: token.OR, Literal: "||", Line: 19, Position: 30},
		{Type: token.IDENT, Literal: "v", Line: 19, Position: 33},
		{Type: token.NOT_EQUAL, Literal: "!=", Line: 19, Position: 35},
		{Type: token.STRING, Literal: "x", Line: 19, Position: 38},
		{Type: token.CONTROL_END, Literal: "}", Line: 19, Position: 42},
		{Type: token.LITERAL, Literal: "complicated condition", Line: 19, Position: 43},
		{Type: token.CONTROL_START, Literal: "%{", Line: 19, Position: 64},
		{Type: token.ENDIF, Literal: "endif", Line: 19, Position: 66},
		{Type: token.CONTROL_END, Literal: "}", Line: 19, Position: 71},

		{Type: token.LITERAL, Literal: "\n\n%{ should recognize escaped string\n$ is escaped dollar character\n\nThat's all, very simplified!\n", Line: 19, Position: 72},
		{Type: token.EOF, Literal: "", Line: 25, Position: 1},
	}

	l := NewFromString(input)

	for i, tt := range expects {
		tok := l.NextToken()

		if diff := cmp.Diff(tt, tok); diff != "" {
			t.Errorf(`Test[%d] failed, diff=%s`, i, diff)
		}
	}
}

func TestIllegalTokens(t *testing.T) {
	tests := []struct {
		input   string
		expects []token.Token
	}{
		{
			input: "Illegal % sign must be doubled to escape.",
			expects: []token.Token{
				{Type: token.ILLEGAL, Literal: "Unexpected '%' character found", Line: 1, Position: 9},
			},
		},
		{
			input: "Illegal $ sign must be doubled to escape.",
			expects: []token.Token{
				{Type: token.ILLEGAL, Literal: "Unexpected '$' character found", Line: 1, Position: 9},
			},
		},
	}

	for _, tt := range tests {
		l := NewFromString(tt.input)

		for i, e := range tt.expects {
			tok := l.NextToken()

			if diff := cmp.Diff(e, tok); diff != "" {
				t.Errorf(`Test[%d] failed, diff=%s`, i, diff)
			}
		}
	}
}

func TestBeginningControl(t *testing.T) {
	tests := []struct {
		input   string
		expects []token.Token
	}{
		{
			input: "%{for i in list}%{if i == 0}${i}%{endif}%{endfor}",
			expects: []token.Token{
				{Type: token.CONTROL_START, Literal: "%{", Line: 1, Position: 1},
				{Type: token.FOR, Literal: "for", Line: 1, Position: 3},
				{Type: token.IDENT, Literal: "i", Line: 1, Position: 7},
				{Type: token.IN, Literal: "in", Line: 1, Position: 9},
				{Type: token.IDENT, Literal: "list", Line: 1, Position: 12},
				{Type: token.CONTROL_END, Literal: "}", Line: 1, Position: 16},

				{Type: token.CONTROL_START, Literal: "%{", Line: 1, Position: 17},
				{Type: token.IF, Literal: "if", Line: 1, Position: 19},
				{Type: token.IDENT, Literal: "i", Line: 1, Position: 22},
				{Type: token.EQUAL, Literal: "==", Line: 1, Position: 24},
				{Type: token.INT, Literal: "0", Line: 1, Position: 27},
				{Type: token.CONTROL_END, Literal: "}", Line: 1, Position: 28},

				{Type: token.INTERPORATION, Literal: "i", Line: 1, Position: 31},

				{Type: token.CONTROL_START, Literal: "%{", Line: 1, Position: 33},
				{Type: token.ENDIF, Literal: "endif", Line: 1, Position: 35},
				{Type: token.CONTROL_END, Literal: "}", Line: 1, Position: 40},

				{Type: token.CONTROL_START, Literal: "%{", Line: 1, Position: 41},
				{Type: token.ENDFOR, Literal: "endfor", Line: 1, Position: 43},
				{Type: token.CONTROL_END, Literal: "}", Line: 1, Position: 49},
				{Type: token.EOF, Literal: "", Line: 1, Position: 50},
			},
		},
	}

	for _, tt := range tests {
		l := NewFromString(tt.input)

		for i, e := range tt.expects {
			tok := l.NextToken()

			if diff := cmp.Diff(e, tok); diff != "" {
				t.Errorf(`Test[%d] failed, diff=%s`, i, diff)
			}
		}
	}

}
