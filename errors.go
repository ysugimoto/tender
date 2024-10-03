package tender

import (
	"fmt"

	"github.com/ysugimoto/tender/token"
)

// RenderError is kind of error type for rendering
type RenderError struct {
	Token   token.Token
	Message string
}

func (r *RenderError) Error() string {
	return fmt.Sprintf(
		`Rendering Error: %s at line %d, position %d`,
		r.Message,
		r.Token.Line,
		r.Token.Position,
	)
}

func UndefinedVariable(t token.Token, name string) *RenderError {
	return &RenderError{
		Token:   t,
		Message: fmt.Sprintf(`Undefined variable "%s"`, name),
	}
}

func NotIterable(t token.Token, name string) *RenderError {
	return &RenderError{
		Token:   t,
		Message: fmt.Sprintf(`Variable "%s" is not iterable, must be slice or map`, name),
	}
}

func UnexpectedType(t token.Token, actual, expect string) *RenderError {
	return &RenderError{
		Token:   t,
		Message: fmt.Sprintf(`Unexpected Type found "%s", expects "%s"`, actual, expect),
	}
}
