package value

import (
	"reflect"
)

type fieldSyntax int

const (
	none fieldSyntax = iota + 1
	dot
	sliceBracket
	mapBracket
)

type Field struct {
	name   string
	syntax fieldSyntax
}

func (f Field) String() string {
	switch f.syntax {
	case sliceBracket:
		return "[" + f.name + "]"
	case mapBracket:
		return "[\"" + f.name + "\"]"
	case dot:
		return "." + f.name
	default:
		return f.name
	}
}

func parseFields(ident string) (Field, []Field) {
	var parsed []Field

	runes := []rune(ident)
	stack := Field{syntax: none}

	for i := 0; i < len(runes); i++ {
		switch ident[i] {
		case '.':
			if stack.name != "" {
				parsed = append(parsed, stack)
				stack = Field{syntax: dot}
			}
		case '[':
			if stack.name != "" {
				parsed = append(parsed, stack)
			}
			stack = Field{syntax: sliceBracket}
			j := i + 1
			for ; j < len(ident); j++ {
				if ident[j] == '"' {
					stack.syntax = mapBracket
					continue
				}
				if ident[j] == ']' {
					break
				}
				stack.name += string(runes[j])
			}
			i = j
			parsed = append(parsed, stack)
			stack = Field{syntax: none}
		default:
			stack.name += string(runes[i])
		}
	}

	if stack.name != "" {
		parsed = append(parsed, stack)
	}
	return parsed[0], parsed[1:]
}

func deref(v reflect.Value) reflect.Value {
	if v.Type().Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func IsSlice(v reflect.Value) bool {
	return deref(v).Type().Kind() == reflect.Slice
}

func IsMap(v reflect.Value) bool {
	return deref(v).Type().Kind() == reflect.Map
}

func IsStruct(v reflect.Value) bool {
	return deref(v).Type().Kind() == reflect.Struct
}

func IsBool(v reflect.Value) bool {
	return deref(v).Type().Kind() == reflect.Bool
}

func IsNumeric(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func IsThuthy(v reflect.Value) (bool, error) {
	switch v.Type().Kind() {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		return v.String() != "", nil
	default:
		return false, NotTruthy(v.Type().Kind().String())
	}
}

// In this project, any comparison types should treat as the following:
//
// bool -> bool
// string -> string
// struct -> struct
// int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64 -> int64
// float32, float64 -> float64
//
// It means that, for example, you can compare int vs uint without type casting.
// The reason why is the main logics must have acculate type conversions,
// but a template is just "view", so types should be flexible to avoid annying type conversions.
func toComparableTypes(left, right reflect.Value) (reflect.Value, reflect.Value, error) {
	if !left.Comparable() {
		return left, right, NotComparable("left expression")
	}
	if !right.Comparable() {
		return left, right, NotComparable("right expression")
	}

	left = deref(left)
	switch left.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		left = reflect.ValueOf(left.Int())
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		left = reflect.ValueOf(int64(left.Uint()))
	case reflect.Float32, reflect.Float64:
		left = reflect.ValueOf(left.Float())
	}

	right = deref(right)
	switch right.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		right = reflect.ValueOf(right.Int())
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		right = reflect.ValueOf(int64(right.Uint()))
	case reflect.Float32, reflect.Float64:
		right = reflect.ValueOf(right.Float())
	}

	return left, right, nil
}
