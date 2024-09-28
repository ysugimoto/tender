package parser

import (
	"fmt"
	"strings"

	"github.com/ysugimoto/tiny-template/token"
)

type ParseError struct {
	Token   token.Token
	Message string
}

func (p *ParseError) Error() string {
	return fmt.Sprintf(
		`Parse Error: %s at line %d, position %d`,
		p.Message,
		p.Token.Line,
		p.Token.Position,
	)
}

func UnexpectedToken(t token.Token, expect ...token.TokenType) *ParseError {
	if len(expect) == 0 {
		return &ParseError{
			Token:   t,
			Message: fmt.Sprintf(`Unexpected Token "%s" found`, t.Type),
		}
	}

	expects := make([]string, len(expect))
	for i, e := range expect {
		expects[i] = `"` + string(e) + `"`
	}

	return &ParseError{
		Token:   t,
		Message: fmt.Sprintf(`Unexpected Token "%s" found, expects %s`, t.Type, strings.Join(expects, " or ")),
	}
}

func TypeConversionError(t token.Token, to string) *ParseError {
	return &ParseError{
		Token:   t,
		Message: fmt.Sprintf(`Failed to convert type to "%s" from "%s"`, to, t.Literal),
	}
}

func UndefinedPrefix(t token.Token) *ParseError {
	return &ParseError{
		Token:   t,
		Message: fmt.Sprintf(`Undefined prefix parser for "%s"`, t.Type),
	}
}
