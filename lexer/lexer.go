package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/ysugimoto/tender/token"
)

type State int

const (
	Default State = iota
	ControlStart
	ControlStartTrim
	Control
	ControlEnd
	ControlEndTrim
	Interporation
)

type Lexer struct {
	r      *bufio.Reader
	char   rune
	line   int
	index  int
	buffer *bytes.Buffer
	lines  []string
	// file   string
	isEOF  bool
	states []State
}

func New(r io.Reader) *Lexer {
	l := &Lexer{
		r:      bufio.NewReader(r),
		line:   1,
		buffer: new(bytes.Buffer),
		states: []State{Default},
	}
	l.readChar()
	return l
}

func NewFromString(input string) *Lexer {
	return New(strings.NewReader(input))
}

func (l *Lexer) pushState(s State) {
	l.states = append(l.states, s)
}

func (l *Lexer) replaceState(s State) {
	l.states[len(l.states)-1] = s
}

func (l *Lexer) popState() {
	if len(l.states) == 0 {
		return
	}
	l.states = l.states[0 : len(l.states)-1]
}

func (l *Lexer) currentState() State {
	if len(l.states) == 0 {
		return Default
	}
	return l.states[len(l.states)-1]
}

func (l *Lexer) readChar() {
	rn, _, err := l.r.ReadRune()
	if err != nil {
		l.char = 0x00
		l.index += 1
		return
	}
	if l.char == 0x0A { // LF
		l.NewLine()
	}
	l.index += 1
	l.char = rn
	l.buffer.WriteRune(rn)
}

func (l *Lexer) peekChar() rune {
	b, err := l.r.Peek(1)
	if err != nil {
		return 0x00
	}
	return rune(b[0])
}

func (l *Lexer) NewLine() {
	l.lines = append(l.lines, strings.TrimRight(l.buffer.String(), "\n"))
	l.buffer = new(bytes.Buffer)
	l.index = 0
	l.line++
}

func (l *Lexer) NextToken() token.Token {
	// Hook states should return without forward reading character
	switch l.currentState() {
	case ControlStart:
		l.replaceState(Control)
		return newToken(token.CONTROL_START, "%{", l.line, l.index-2)
	case ControlStartTrim:
		l.replaceState(Control)
		t := newToken(token.CONTROL_START, "%{~", l.line, l.index-3)
		t.LeftTrim = true
		return t
	case ControlEnd:
		l.popState()
		return newToken(token.CONTROL_END, "}", l.line, l.index-1)
	case ControlEndTrim:
		l.popState()
		t := newToken(token.CONTROL_END, "~}", l.line, l.index-2)
		t.RightTrim = true
		return t
	}

	// Following state must forward reading
	defer l.readChar()

	switch l.currentState() {
	case Control:
		return l.nextControlToken()
	case Interporation:
		t := l.nextInterporationToken()
		t.Position -= 2
		return t
	default:
		return l.nextToken()
	}
}

func (l *Lexer) nextToken() token.Token {
	var stack []rune

	// Store start line and index
	index, line := l.index, l.line

	for {
		switch l.char {
		case '%':
			switch l.peekChar() {
			case '%': // escaped percent sign
				stack = append(stack, l.char)
				l.readChar()
				goto CONT
			case '{':
				l.readChar()
				if l.peekChar() == '~' { // trim control
					l.readChar()
					if len(stack) == 0 {
						l.pushState(Control)
						t := newToken(token.CONTROL_START, "%{~", line, index)
						t.LeftTrim = true
						return t
					}
					l.pushState(ControlStartTrim)
				} else {
					if len(stack) == 0 {
						l.pushState(Control)
						return newToken(token.CONTROL_START, "%{", line, index)
					}
					l.pushState(ControlStart)
				}
				return newToken(token.LITERAL, string(stack), line, index)
			default:
				return newToken(token.ILLEGAL, "Unexpected '%' character found", l.line, l.index)
			}
		case '$':
			switch l.peekChar() {
			case '$': // escaped dollar sign
				stack = append(stack, l.char)
				l.readChar()
				goto CONT
			case '{':
				l.readChar()
				if len(stack) == 0 {
					l.readChar()
					t := l.nextInterporationToken()
					t.Position -= 2
					return t
				}
				l.pushState(Interporation)
				return newToken(token.LITERAL, string(stack), line, index)
			default:
				return newToken(token.ILLEGAL, "Unexpected '$' character found", l.line, l.index)
			}
		case 0x00: // EOF
			if !l.isEOF {
				l.NewLine()
				l.isEOF = true
			}
			if len(stack) > 0 {
				return newToken(token.LITERAL, string(stack), line, index)
			}
			return newToken(token.EOF, "", line, index)
		default:
			stack = append(stack, l.char)
		}
	CONT:
		l.readChar()
	}
}

func (l *Lexer) nextControlToken() token.Token {
	l.skipWhitespace()

	index, line := l.index, l.line
	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			return newToken(token.EQUAL, "==", line, index)
		}
		return newToken(token.ILLEGAL, "=", l.line, l.index)
	case '-':
		return newToken(token.MINUS, "-", line, index)
	case '}': // end control
		l.popState()
		return newToken(token.CONTROL_END, "}", l.line, l.index)
	case '(':
		return newToken(token.LEFT_PAREN, "(", line, index)
	case ')':
		return newToken(token.RIGHT_PAREN, ")", line, index)
	case ',':
		return newToken(token.COMMA, ",", line, index)
	case '"':
		return newToken(token.STRING, l.readString(), line, index)
	case '|':
		if l.peekChar() == '|' { // "||"
			l.readChar()
			return newToken(token.OR, "||", line, index)
		}
		return newToken(token.ILLEGAL, "|", l.line, l.index)
	case '&':
		if l.peekChar() == '&' { // "&&"
			l.readChar()
			return newToken(token.AND, "&&", line, index)
		}
		return newToken(token.ILLEGAL, "&", l.line, l.index)
	case '>':
		switch l.peekChar() {
		case '=': // ">="
			l.readChar()
			return newToken(token.GREATER_THAN_EQUAL, ">=", line, index)
		default:
			return newToken(token.GREATER_THAN, ">", line, index)
		}
	case '<':
		switch l.peekChar() {
		case '=': // "<="
			l.readChar()
			return newToken(token.LESS_THAN_EQUAL, "<=", line, index)
		default:
			return newToken(token.LESS_THAN, "<", line, index)
		}
	case '!':
		switch l.peekChar() {
		case '=': // "!="
			l.readChar()
			return newToken(token.NOT_EQUAL, "!=", line, index)
		default:
			return newToken(token.NOT, "!", line, index)
		}
	case '~':
		if l.peekChar() == '}' {
			l.readChar()
			l.popState()
			t := newToken(token.CONTROL_END, "~}", line, index)
			t.RightTrim = true
			return t
		}
		return newToken(token.ILLEGAL, string(l.char), line, index)
	case 0x0A: // LF
		return newToken(token.LF, "\n", line, index)
	case 0x00:
		return newToken(token.ILLEGAL, "", line, index)
	default:
		switch {
		case isLetter(l.char):
			literal, ok := l.readLiteral()
			if !ok {
				return newToken(token.ILLEGAL, "", line, index)
			}
			return newToken(token.LookupIdent(literal), literal, line, index)
		case isDigit(l.char):
			num := l.readNumber()

			// If literal contains ".", token should be FLOAT
			if strings.Count(num, ".") == 1 {
				return newToken(token.FLOAT, num, line, index)
			}
			return newToken(token.INT, num, line, index)
		default:
			return newToken(token.ILLEGAL, string(l.char), line, index)
		}
	}
}

func (l *Lexer) nextInterporationToken() token.Token {
	index, line := l.index, l.line

	var literal string
	for l.char != '}' {
		switch l.char {
		case 0x00:
			return newToken(token.ILLEGAL, "", line, index)
		default:
			// Interporation only accepts identifier
			// TODO: Maybe we implement filter process in the future
			switch {
			case isLetter(l.char):
				lt, ok := l.readLiteral()
				if !ok {
					return newToken(token.ILLEGAL, string(l.char), line, index)
				}
				literal = lt
			default:
				return newToken(token.ILLEGAL, string(l.char), line, index)
			}
		}
		l.readChar()
	}
	l.popState()
	return newToken(token.INTERPORATION, literal, line, index)
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readLiteral() (string, bool) {
	literal := l.readIdentifier()

	// Read more neighbor digit, dot, underscore, left bracket
	for {
		peek := l.peekChar()
		switch {
		case isLetter(peek):
			l.readChar()
			literal += l.readIdentifier()
		case peek == '_' || peek == '.' || isDigit(l.char):
			l.readChar()
			literal += string(l.char)
		// Array or map indexing as `["..."]`
		case peek == '[':
			l.readChar()
			literal += string(l.char)
			switch {
			case l.peekChar() == '"': // string - object key indexing
				l.readChar()
				literal += `"` + l.readString() + `"`
			case isDigit(l.peekChar()): // digit - array indexing
				l.readChar()
				literal += l.readNumber()
			default: // illegal
				return "", false
			}

			if l.peekChar() != ']' {
				return "", false
			}
			l.readChar()
			literal += string(l.char)
		default:
			return literal, true
		}
	}
}

func (l *Lexer) readIdentifier() string {
	rs := []rune{l.char}
	for isLetter(l.peekChar()) {
		l.readChar()
		rs = append(rs, l.char)
	}
	return string(rs)
}

func (l *Lexer) readString() string {
	var rs []rune
	l.readChar()
	for {
		if l.char == '"' || l.char == 0x00 {
			break
		}
		rs = append(rs, l.char)
		l.readChar()
	}

	return string(rs)
}

func (l *Lexer) readNumber() string {
	rs := []rune{l.char}
	for isDigit(l.peekChar()) {
		l.readChar()
		rs = append(rs, l.char)
	}
	return string(rs)
}

func newToken(tokenType token.TokenType, literal string, line, index int) token.Token {
	return token.Token{
		Type:     tokenType,
		Literal:  literal,
		Line:     line,
		Position: index,
	}
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_'
}

func isDigit(r rune) bool {
	// Digit allows "." character to parse literal is INTEGER of FLOAT.
	return (r >= '0' && r <= '9') || r == '.'
}
