package tender

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/value"
)

// Evaluate expression inside if condition
func (t *Template) evaluateExpression(expr ast.Expression) (reflect.Value, error) {
	switch tt := expr.(type) {
	case *ast.Ident:
		v, err := t.lookupVariable(tt.Value)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return v, nil
	case *ast.String:
		return reflect.ValueOf(tt.Value), nil
	case *ast.Int:
		return reflect.ValueOf(tt.Value), nil
	case *ast.Float:
		return reflect.ValueOf(tt.Value), nil
	case *ast.Bool:
		return reflect.ValueOf(tt.Value), nil

	case *ast.PrefixExpression:
		return t.evaluatePrefixExpression(tt)
	case *ast.InfixExpression:
		return t.evaluateInfixExpression(tt)
	case *ast.GroupedExpression:
		return t.evaluateGroupedExpression(tt)
	}

	return value.Null, errors.WithStack(&RenderError{
		Token:   expr.GetToken(),
		Message: "Unexpected expression found",
	})
}

func (t *Template) evaluatePrefixExpression(expr *ast.PrefixExpression) (reflect.Value, error) {
	right, err := t.evaluateExpression(expr.Right)
	if err != nil {
		return value.Null, errors.WithStack(err)
	}
	switch expr.Operator {
	case "!":
		switch right.Type().Kind() {
		case reflect.Bool:
			return reflect.ValueOf(!right.Bool()), nil
		case reflect.String:
			return reflect.ValueOf(right.String() == ""), nil
		default:
			return value.Null, errors.WithStack(&RenderError{
				Token:   expr.GetToken(),
				Message: `Unexpected "!" prefix operator for ` + right.Type().Kind().String(),
			})
		}
	case "-":
		switch right.Type().Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			return reflect.ValueOf(-right.Int()), nil
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			return reflect.ValueOf(-right.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(-right.Float()), nil
		default:
			return value.Null, errors.WithStack(&RenderError{
				Token:   expr.GetToken(),
				Message: `Unexpected "-" prefix operator for ` + right.Type().Kind().String(),
			})
		}
	default:
		return value.Null, errors.WithStack(&RenderError{
			Token:   expr.GetToken(),
			Message: `Unexpected prefix operator "` + expr.Operator + `" found`,
		})
	}
}

func (t *Template) evaluateInfixExpression(expr *ast.InfixExpression) (reflect.Value, error) {
	left, err := t.evaluateExpression(expr.Left)
	if err != nil {
		return value.Null, errors.WithStack(err)
	}
	right, err := t.evaluateExpression(expr.Right)
	if err != nil {
		return value.Null, errors.WithStack(err)
	}
	switch expr.Operator {
	case "==":
		cmp, err := value.Equal(left, right)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return reflect.ValueOf(cmp), nil
	case "!=":
		cmp, err := value.NotEqual(left, right)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return reflect.ValueOf(cmp), nil
	case ">":
		cmp, err := value.GreaterThan(left, right)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return reflect.ValueOf(cmp), nil
	case ">=":
		cmp, err := value.GreaterThanEqual(left, right)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return reflect.ValueOf(cmp), nil
	case "<":
		cmp, err := value.LessThan(left, right)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return reflect.ValueOf(cmp), nil
	case "<=":
		cmp, err := value.LessThanEqual(left, right)
		if err != nil {
			return value.Null, errors.WithStack(err)
		}
		return reflect.ValueOf(cmp), nil
	case "&&":
		if !value.IsBool(left) {
			if err != nil {
				return value.Null, errors.WithStack(
					UnexpectedType(expr.Left.GetToken(), left.Type().Kind().String(), "bool"),
				)
			}
		}
		if !value.IsBool(right) {
			if err != nil {
				return value.Null, errors.WithStack(
					UnexpectedType(expr.Right.GetToken(), right.Type().Kind().String(), "bool"),
				)
			}
		}
		return reflect.ValueOf(left.Bool() && right.Bool()), nil
	case "||":
		if !value.IsBool(left) {
			if err != nil {
				return value.Null, errors.WithStack(
					UnexpectedType(expr.Left.GetToken(), left.Type().Kind().String(), "bool"),
				)
			}
		}
		if !value.IsBool(right) {
			if err != nil {
				return value.Null, errors.WithStack(
					UnexpectedType(expr.Right.GetToken(), right.Type().Kind().String(), "bool"),
				)
			}
		}
		return reflect.ValueOf(left.Bool() || right.Bool()), nil
	default:
		return value.Null, errors.WithStack(&RenderError{
			Token:   expr.GetToken(),
			Message: `Unexpected operation "` + expr.Operator + `" found`,
		})
	}
}

func (t *Template) evaluateGroupedExpression(expr *ast.GroupedExpression) (reflect.Value, error) {
	v, err := t.evaluateExpression(expr.Right)
	if err != nil {
		return value.Null, errors.WithStack(err)
	}
	return v, nil
}
