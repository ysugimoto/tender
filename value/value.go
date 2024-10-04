package value

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Value map[string]reflect.Value

var Null = reflect.ValueOf(nil)
var zero = reflect.Value{}

func (v Value) Resolve(ident string) (reflect.Value, error) {
	first, subFields := parseFields(ident)

	variable, ok := v[first.name]
	if !ok {
		return Null, UndefinedVariable(first.name)
	}

	names := pool.Get().(*bytes.Buffer) // nolint:errcheck
	defer pool.Put(names)

	names.Reset()

	child := deref(variable)
	for _, field := range subFields {
		switch {
		case IsMap(child):
			child = child.MapIndex(reflect.ValueOf(field.name))
			if child == zero {
				return Null, UndefinedKey(names.String(), field.name)
			}
			child = reflect.ValueOf(child.Interface())
		case IsSlice(child):
			idx, err := strconv.Atoi(field.name)
			if err != nil {
				return Null, UnaccessibleIndex(names.String(), field.name)
			}
			if idx > child.Len()-1 {
				return Null, UndefinedIndex(names.String(), field.name)
			}
			child = reflect.ValueOf(child.Index(idx).Interface())
		case IsStruct(child):
			// Struct field must start with Upper-case alphabet, valid field name
			// otherwise reflect.Value.FieldByName will cause panic
			if field.name[0] < 0x41 || field.name[0] > 0x5A {
				return Null, InvalidFieldAccess(names.String(), field.name)
			}
			child = child.FieldByName(field.name)
			if child == zero {
				return Null, UndefinedField(names.String(), field.name)
			}
			child = reflect.ValueOf(child.Interface())
		default:
			return Null, UndefinedVariable(field.name)
		}
		child = deref(child)
		names.WriteString(field.String())
	}

	return child, nil
}

// Compare values with "==" operator
func Equal(left, right reflect.Value) (bool, error) {
	left, right, err := toComparableTypes(left, right)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if left.Type().Kind() != right.Type().Kind() {
		return false, TypeMismatch(left.Type().String(), right.Type().String())
	}
	return left.Equal(right), nil
}

// Compare values with "!=" operator
func NotEqual(left, right reflect.Value) (bool, error) {
	eq, err := Equal(left, right)
	if err != nil {
		return false, err
	}
	return !eq, nil
}

// Compare values with ">" operator
func GreaterThan(left, right reflect.Value) (bool, error) {
	left, right, err := toComparableTypes(left, right)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if left.Type().Kind() != right.Type().Kind() {
		return false, TypeMismatch(left.Type().Kind().String(), right.Type().Kind().String())
	}
	if !IsNumeric(left) {
		return false, NotNumeric(left.Type().Kind().String())
	} else if !IsNumeric(right) {
		return false, NotNumeric(right.Type().Kind().String())
	}

	switch left.Type().Kind() {
	case reflect.Int64:
		return left.Int() > right.Int(), nil
	case reflect.Float64:
		return left.Float() > right.Float(), nil
	default:
		return false, &ComparisonError{
			Message: "Unknown",
		}
	}
}

// Compare values with ">=" operator
func GreaterThanEqual(left, right reflect.Value) (bool, error) {
	left, right, err := toComparableTypes(left, right)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if left.Type().Kind() != right.Type().Kind() {
		return false, TypeMismatch(left.Type().Kind().String(), right.Type().Kind().String())
	}
	if !IsNumeric(left) {
		return false, NotNumeric(left.Type().Kind().String())
	} else if !IsNumeric(right) {
		return false, NotNumeric(right.Type().Kind().String())
	}

	switch left.Type().Kind() {
	case reflect.Int64:
		return left.Int() >= right.Int(), nil
	case reflect.Float64:
		return left.Float() >= right.Float(), nil
	default:
		return false, &ComparisonError{
			Message: "Unknown",
		}
	}
}

// Compare values with "<" operator
func LessThan(left, right reflect.Value) (bool, error) {
	left, right, err := toComparableTypes(left, right)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if left.Type().Kind() != right.Type().Kind() {
		return false, TypeMismatch(left.Type().Kind().String(), right.Type().Kind().String())
	}
	if !IsNumeric(left) {
		return false, NotNumeric(left.Type().Kind().String())
	} else if !IsNumeric(right) {
		return false, NotNumeric(right.Type().Kind().String())
	}

	switch left.Type().Kind() {
	case reflect.Int64:
		return left.Int() < right.Int(), nil
	case reflect.Float64:
		return left.Float() < right.Float(), nil
	default:
		return false, &ComparisonError{
			Message: "Unknown",
		}
	}
}

// Compare values with "<=" operator
func LessThanEqual(left, right reflect.Value) (bool, error) {
	left, right, err := toComparableTypes(left, right)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if left.Type().Kind() != right.Type().Kind() {
		return false, TypeMismatch(left.Type().Kind().String(), right.Type().Kind().String())
	}
	if !IsNumeric(left) {
		return false, NotNumeric(left.Type().Kind().String())
	} else if !IsNumeric(right) {
		return false, NotNumeric(right.Type().Kind().String())
	}

	switch left.Type().Kind() {
	case reflect.Int64:
		return left.Int() <= right.Int(), nil
	case reflect.Float64:
		return left.Float() <= right.Float(), nil
	default:
		return false, &ComparisonError{
			Message: "Unknown",
		}
	}
}

// Stringify reflect.Value
func ToString(v reflect.Value) string {
	v = deref(v)

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		return v.String()
	case reflect.Slice:
		values := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			values[i] = ToString(v.Index(i))
		}
		return "[" + strings.Join(values, ", ") + "]"
	case reflect.Map:
		values := make([]string, v.Len())
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			a := ToString(keys[i])
			b := ToString(keys[j])
			return a < b
		})
		for i := 0; i < len(keys); i++ {
			values[i] = ToString(keys[i]) + ": " + ToString(v.MapIndex(keys[i]))
		}
		return "{" + strings.Join(values, ", ") + "}"
	case reflect.Struct:
		values := make([]string, v.NumField())
		index := 0

		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			fv := v.FieldByName(f.Name)
			if fv.IsZero() {
				continue
			}
			values[index] = f.Name + ": " + ToString(fv)
			index++
		}
		return "{" + strings.Join(values[:index], ", ") + "}"
	default:
		return fmt.Sprintf("%v", v)
	}
}
