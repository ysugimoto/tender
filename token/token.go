package token

import "fmt"

type TokenType string

type Token struct {
	Type      TokenType
	Literal   string
	LeftTrim  bool
	RightTrim bool

	Line     int
	Position int
	File     string
}

func (t Token) String() string {
	return fmt.Sprintf(
		"type=%s, literal=`%s`, line=%d, position=%d",
		string(t.Type), t.Literal, t.Line, t.Position,
	)
}

var Null = Token{
	Type:    ILLEGAL,
	Literal: "NULL",
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	FLOAT  = "FLOAT"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	LF     = "LF"

	// Template spec
	LITERAL       = "LITERAL"       // arbitrary strings outside control
	INTERPORATION = "INTERPORATION" // "${"
	CONTROL_START = "CONTROL_START" // "%{"
	CONTROL_END   = "CONTROL_END"   // "}"

	// Operators
	EQUAL              = "EQUAL"              // "=="
	NOT_EQUAL          = "NOTEQUAL"           // "!="
	GREATER_THAN       = "GREATER_THAN"       // ">"
	LESS_THAN          = "LESS_THAN"          // "<"
	GREATER_THAN_EQUAL = "GREATER_THAN_EQUAL" // >="
	LESS_THAN_EQUAL    = "LESS_THAN_EQUAL"    // <="
	AND                = "AND"                // "&&"
	OR                 = "OR"                 // "||"

	// Punctuation
	LEFT_PAREN    = "LEFT_PAREN"    // "("
	RIGHT_PAREN   = "RIGHT_PAREN"   // ")"
	LEFT_BRACKET  = "LEFT_BRACKET"  // "["
	RIGHT_BRACKET = "RIGHT_BRACKET" // "]"
	COMMA         = "COMMA"         // ","
	NOT           = "NOT"           // "!"
	TILDA         = "TILDA"         // "~"

	// Keywords
	FOR    = "FOR"    // for
	IN     = "IN"     // in
	ENDFOR = "ENDFOR" // endfor
	IF     = "IF"     // if
	ELSEIF = "ELSEIF" // elseif
	ELSE   = "ELSE"   // else
	ENDIF  = "ENDIF"  // endif
)

var keywords = map[string]TokenType{
	"for":    FOR,
	"in":     IN,
	"endfor": ENDFOR,
	"if":     IF,
	"else":   ELSE,
	"elseif": ELSEIF,
	"endif":  ENDIF,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if v, ok := keywords[ident]; ok {
		return v
	}

	return IDENT
}
