package value

type ValueError struct {
	Message string
}

func (e *ValueError) Error() string {
	return e.Message
}

func UndefinedVariable(name string) *ValueError {
	return &ValueError{
		Message: `Undefined variable "` + name + `"`,
	}
}

func UndefinedIndex(name, index string) *ValueError {
	return &ValueError{
		Message: `Undefined index "` + index + `" for slice value of "` + name + `"`,
	}
}

func UndefinedKey(name, key string) *ValueError {
	return &ValueError{
		Message: `Undefined key "` + key + `" for map value of "` + name + `"`,
	}
}

func UndefinedField(name, field string) *ValueError {
	return &ValueError{
		Message: `Undefined field "` + field + `" for struct value of "` + name + `"`,
	}
}

func InvalidFieldAccess(name, field string) *ValueError {
	return &ValueError{
		Message: `Invalid struct field accessing (unexported or invalid syntax) "` + field + `" for struct value of "` + name + `"`,
	}
}

func UnaccessibleIndex(name, index string) *ValueError {
	return &ValueError{
		Message: `Unaccessible index "` + index + `" for slice "` + name + `"`,
	}
}

func NotIterable(name string) *ValueError {
	return &ValueError{
		Message: name + " is not iterable",
	}
}

type ComparisonError struct {
	Message string
}

func (e *ComparisonError) Error() string {
	return e.Message
}

func NotComparable(name string) *ValueError {
	return &ValueError{
		Message: name + " is not comparable value",
	}
}

func TypeMismatch(left, right string) *ValueError {
	return &ValueError{
		Message: `Comparison type is mismatch, left="` + left + `" and right="` + right + `"`,
	}
}

func NotNumeric(name string) *ValueError {
	return &ValueError{
		Message: name + ` is not numeric value`,
	}
}

func NotTruthy(name string) *ValueError {
	return &ValueError{
		Message: name + ` is not truthy type. it must be bool or string`,
	}
}
