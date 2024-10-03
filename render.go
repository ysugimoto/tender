package tender

import (
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/value"
)

// Lookup variables from current scope local variables,
// or global assigned variables if local variable is not found.
func (t *Template) lookupVariable(name string) (reflect.Value, error) {
	if len(t.locals) > 0 {
		local := t.locals[len(t.locals)-1]
		if v, err := local.Resolve(name); err == nil {
			return v, nil
		}
	}
	return t.global.Resolve(name)
}

// render the template from parsed AST Nodes.
func (t *Template) render(nodes []ast.Node) (string, error) {
	var stack []string

	for i := range nodes {
		switch n := nodes[i].(type) {
		case *ast.Literal:
			stack = append(stack, n.Value)
		case *ast.If:
			if n.Token.LeftTrim {
				trimRightSpaceLastString(stack)
			}

			v, err := t.renderIfControl(n)
			if err != nil {
				return "", errors.WithStack(err)
			}
			stack = append(stack, v)

			if n.End.Token.RightTrim {
				if i+1 < len(nodes)-1 {
					if l, ok := nodes[i+i].(*ast.Literal); ok {
						l.Value = trimLeftSpace(l.Value)
					}
				}
			}
		case *ast.For:
			if n.Token.LeftTrim {
				trimRightSpaceLastString(stack)
			}

			v, err := t.renderForControl(n)
			if err != nil {
				return "", errors.WithStack(err)
			}
			stack = append(stack, v)

			if n.End.Token.RightTrim {
				if i+1 < len(nodes)-1 {
					if l, ok := nodes[i+i].(*ast.Literal); ok {
						l.Value = trimLeftSpace(l.Value)
					}
				}
			}
		case *ast.Interporation:
			if isEnvironmentVariable(n.Value.Value) {
				v, ok := os.LookupEnv(n.Value.Value)
				if !ok {
					return "", errors.WithStack(errors.New(
						`environment variable "` + n.Value.Value + `" is not specified"`,
					))
				}
				stack = append(stack, v)
			} else {
				v, err := t.lookupVariable(n.Value.Value)
				if err != nil {
					return "", errors.WithStack(err)
				}
				stack = append(stack, value.ToString(v))
			}
		default:
			return "", errors.New("Unexpected node found")
		}
	}

	return strings.Join(stack, ""), nil
}

// Render the for control syntax
func (t *Template) renderForControl(node *ast.For) (string, error) {
	var stack []string

	// Check iterator variable is assigned
	iterator, err := t.lookupVariable(node.Iterator.Value)
	if err != nil {
		return "", errors.WithStack(UndefinedVariable(node.Iterator.Token, node.Iterator.Value))
	}

	// For loop iterator value must be a slice of map
	switch {
	case value.IsMap(iterator):
		keys := iterator.MapKeys()
		// Map look key is unordered so we will sort alphabetically
		sort.Slice(keys, func(i, j int) bool {
			a := value.ToString(keys[i])
			b := value.ToString(keys[j])
			return a < b
		})

		for i := 0; i < len(keys); i++ {
			iteration, err := t.renderForIteration(node, keys[i], iterator.MapIndex(keys[i]))
			if err != nil {
				return "", errors.WithStack(err)
			}
			if node.Token.RightTrim {
				iteration = trimLeftSpace(iteration)
			}
			if node.End.Token.LeftTrim {
				iteration = trimRightSpace(iteration)
			}
			stack = append(stack, iteration)
		}
	case value.IsSlice(iterator):
		for i := 0; i < iterator.Len(); i++ {
			iteration, err := t.renderForIteration(node, reflect.ValueOf(i), iterator.Index(i))
			if err != nil {
				return "", errors.WithStack(err)
			}
			if node.Token.RightTrim {
				iteration = trimLeftSpace(iteration)
			}
			if node.End.Token.LeftTrim {
				iteration = trimRightSpace(iteration)
			}
			stack = append(stack, iteration)
		}
	default:
		// Otherwise, raise NotIterable error
		return "", errors.WithStack(NotIterable(node.Iterator.Token, node.Iterator.Value))
	}

	return strings.Join(stack, ""), nil
}

// Process the one interation for the "for" block
func (t *Template) renderForIteration(node *ast.For, key, val reflect.Value) (string, error) {
	// Assign key and value to local variable
	local := value.Value{
		node.Arg1.Value: key,
	}
	if node.Arg2 != nil {
		local[node.Arg2.Value] = val
	}

	// Push current local scoped values
	t.locals = append(t.locals, local)
	defer func() {
		// Pop current local scoped values after the iteration
		t.locals = t.locals[0 : len(t.locals)-1]
	}()

	ret, err := t.render(node.Block)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return ret, nil
}

// Render the "if" syntax
func (t *Template) renderIfControl(node *ast.If) (string, error) {
	cond, err := t.evaluateExpression(node.Condition)
	if err != nil {
		return "", errors.WithStack(err)
	}

	truthy, err := value.IsThuthy(cond)
	if err != nil {
		return "", errors.WithStack(&RenderError{
			Token:   node.Condition.GetToken(),
			Message: err.Error(),
		})
	}
	// If first if condition could evaluate as "true", render consequence block
	if truthy {
		v, err := t.render(node.Consequence)
		if err != nil {
			return "", errors.WithStack(err)
		}
		if node.Token.RightTrim {
			v = trimLeftSpace(v)
		}
		switch {
		case len(node.Another) > 0:
			if node.Another[0].Token.LeftTrim {
				v = trimRightSpace(v)
			}
		case node.Alternative != nil:
			if node.Alternative.Token.LeftTrim {
				v = trimRightSpace(v)
			}
		case node.End.Token.LeftTrim:
			v = trimRightSpace(v)
		}
		return v, nil
	}

	// Evaluate else if syntax as possible as we find
	for i := range node.Another {
		n := node.Another[i]
		cond, err := t.evaluateExpression(n.Condition)
		if err != nil {
			return "", errors.WithStack(err)
		}
		truthy, err := value.IsThuthy(cond)
		if err != nil {
			return "", errors.WithStack(&RenderError{
				Token:   node.Condition.GetToken(),
				Message: err.Error(),
			})
		}
		if truthy {
			v, err := t.render(n.Consequence)
			if err != nil {
				return "", errors.WithStack(err)
			}
			if node.Token.RightTrim {
				v = trimLeftSpace(v)
			}

			switch {
			case i+1 < len(node.Another)-1:
				if node.Another[i+1].Token.LeftTrim {
					v = trimRightSpace(v)
				}
			case node.Alternative != nil:
				if node.Alternative.Token.LeftTrim {
					v = trimRightSpace(v)
				}
			case node.End.Token.LeftTrim:
				v = trimRightSpace(v)
			}
			return v, nil
		}
	}

	// Evaluate else syntax if found
	if node.Alternative != nil {
		v, err := t.render(node.Alternative.Consequence)
		if err != nil {
			return "", errors.WithStack(err)
		}
		if node.Alternative.Token.RightTrim {
			v = trimLeftSpace(v)
		}
		if node.End.Token.LeftTrim {
			v = trimRightSpace(v)
		}
		return v, nil
	}

	return "", nil
}
