package value

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestResolveValue(t *testing.T) {
	global := Value{
		"string": reflect.ValueOf("bar"),
		"int":    reflect.ValueOf(10),
		"bool":   reflect.ValueOf(true),
		"slice":  reflect.ValueOf([]string{"a", "b", "c"}),
		"map": reflect.ValueOf(map[string]any{
			"key":    "Value",
			"nested": []int{1},
		}),
		"struct": reflect.ValueOf(struct {
			Foo string
		}{Foo: "bar"}),
	}

	tests := []struct {
		index  string
		expect reflect.Value
	}{
		{
			index:  "string",
			expect: reflect.ValueOf("bar"),
		},
		{
			index:  "int",
			expect: reflect.ValueOf(10),
		},
		{
			index:  "bool",
			expect: reflect.ValueOf(true),
		},
		{
			index:  "slice[0]",
			expect: reflect.ValueOf("a"),
		},
		{
			index:  `map["key"]`,
			expect: reflect.ValueOf("Value"),
		},
		{
			index:  `map.key`,
			expect: reflect.ValueOf("Value"),
		},
		{
			index:  `map["nested"][0]`,
			expect: reflect.ValueOf(1),
		},
		{
			index:  `struct.Foo`,
			expect: reflect.ValueOf("bar"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.index, func(t *testing.T) {
			v, err := global.Resolve(tt.index)
			if err != nil {
				t.Errorf("Undexpected error: %s", err)
				return
			}
			if diff := cmp.Diff(tt.expect, v); diff != "" {
				t.Errorf("Resolved Value unmatch, diff=%s", diff)
			}
		})
	}
}

type testStruct struct {
	name string
}

func TestCompareValuesEqual(t *testing.T) {
	tests := []struct {
		left    any
		right   any
		expect  bool
		isError bool
	}{
		// left is bool
		{left: true, right: true, expect: true},
		{left: true, right: int(1), isError: true},
		{left: true, right: int8(1), isError: true},
		{left: true, right: int16(1), isError: true},
		{left: true, right: int32(1), isError: true},
		{left: true, right: int64(1), isError: true},
		{left: true, right: uint(1), isError: true},
		{left: true, right: uint8(1), isError: true},
		{left: true, right: uint16(1), isError: true},
		{left: true, right: uint32(1), isError: true},
		{left: true, right: uint64(1), isError: true},
		{left: true, right: float32(1), isError: true},
		{left: true, right: float64(1), isError: true},
		{left: true, right: "foo", isError: true},
		{left: true, right: testStruct{name: "foo"}, isError: true},
		{left: true, right: []string{"foo"}, isError: true},

		// left is int
		{left: int(1), right: true, isError: true},
		{left: int(1), right: int(1), expect: true},
		{left: int(1), right: int8(1), expect: true},
		{left: int(1), right: int16(1), expect: true},
		{left: int(1), right: int32(1), expect: true},
		{left: int(1), right: int64(1), expect: true},
		{left: int(1), right: uint(1), expect: true},
		{left: int(1), right: uint8(1), expect: true},
		{left: int(1), right: uint16(1), expect: true},
		{left: int(1), right: uint32(1), expect: true},
		{left: int(1), right: uint64(1), expect: true},
		{left: int(1), right: float32(1), isError: true},
		{left: int(1), right: float64(1), isError: true},
		{left: int(1), right: "foo", isError: true},
		{left: int(1), right: testStruct{name: "foo"}, isError: true},
		{left: int(1), right: []string{"foo"}, isError: true},

		// left is int8
		{left: int8(1), right: true, isError: true},
		{left: int8(1), right: int(1), expect: true},
		{left: int8(1), right: int8(1), expect: true},
		{left: int8(1), right: int16(1), expect: true},
		{left: int8(1), right: int32(1), expect: true},
		{left: int8(1), right: int64(1), expect: true},
		{left: int8(1), right: uint(1), expect: true},
		{left: int8(1), right: uint8(1), expect: true},
		{left: int8(1), right: uint16(1), expect: true},
		{left: int8(1), right: uint32(1), expect: true},
		{left: int8(1), right: uint64(1), expect: true},
		{left: int8(1), right: float32(1), isError: true},
		{left: int8(1), right: float64(1), isError: true},
		{left: int8(1), right: "foo", isError: true},
		{left: int8(1), right: testStruct{name: "foo"}, isError: true},
		{left: int8(1), right: []string{"foo"}, isError: true},

		// left is int16
		{left: int16(1), right: true, isError: true},
		{left: int16(1), right: int(1), expect: true},
		{left: int16(1), right: int8(1), expect: true},
		{left: int16(1), right: int16(1), expect: true},
		{left: int16(1), right: int32(1), expect: true},
		{left: int16(1), right: int64(1), expect: true},
		{left: int16(1), right: uint(1), expect: true},
		{left: int16(1), right: uint8(1), expect: true},
		{left: int16(1), right: uint16(1), expect: true},
		{left: int16(1), right: uint32(1), expect: true},
		{left: int16(1), right: uint64(1), expect: true},
		{left: int16(1), right: float32(1), isError: true},
		{left: int16(1), right: float64(1), isError: true},
		{left: int16(1), right: "foo", isError: true},
		{left: int16(1), right: testStruct{name: "foo"}, isError: true},
		{left: int16(1), right: []string{"foo"}, isError: true},

		// left is int32
		{left: int32(1), right: true, isError: true},
		{left: int32(1), right: int(1), expect: true},
		{left: int32(1), right: int8(1), expect: true},
		{left: int32(1), right: int16(1), expect: true},
		{left: int32(1), right: int32(1), expect: true},
		{left: int32(1), right: int64(1), expect: true},
		{left: int32(1), right: uint(1), expect: true},
		{left: int32(1), right: uint8(1), expect: true},
		{left: int32(1), right: uint16(1), expect: true},
		{left: int32(1), right: uint32(1), expect: true},
		{left: int32(1), right: uint64(1), expect: true},
		{left: int32(1), right: float32(1), isError: true},
		{left: int32(1), right: float64(1), isError: true},
		{left: int32(1), right: "foo", isError: true},
		{left: int32(1), right: testStruct{name: "foo"}, isError: true},
		{left: int32(1), right: []string{"foo"}, isError: true},

		// left is int64
		{left: int64(1), right: true, isError: true},
		{left: int64(1), right: int(1), expect: true},
		{left: int64(1), right: int8(1), expect: true},
		{left: int64(1), right: int16(1), expect: true},
		{left: int64(1), right: int32(1), expect: true},
		{left: int64(1), right: int64(1), expect: true},
		{left: int64(1), right: uint(1), expect: true},
		{left: int64(1), right: uint8(1), expect: true},
		{left: int64(1), right: uint16(1), expect: true},
		{left: int64(1), right: uint32(1), expect: true},
		{left: int64(1), right: uint64(1), expect: true},
		{left: int64(1), right: float32(1), isError: true},
		{left: int64(1), right: float64(1), isError: true},
		{left: int64(1), right: "foo", isError: true},
		{left: int64(1), right: testStruct{name: "foo"}, isError: true},
		{left: int64(1), right: []string{"foo"}, isError: true},

		// left is uint
		{left: uint(1), right: true, isError: true},
		{left: uint(1), right: int(1), expect: true},
		{left: uint(1), right: int8(1), expect: true},
		{left: uint(1), right: int16(1), expect: true},
		{left: uint(1), right: int32(1), expect: true},
		{left: uint(1), right: int64(1), expect: true},
		{left: uint(1), right: uint(1), expect: true},
		{left: uint(1), right: uint8(1), expect: true},
		{left: uint(1), right: uint16(1), expect: true},
		{left: uint(1), right: uint32(1), expect: true},
		{left: uint(1), right: uint64(1), expect: true},
		{left: uint(1), right: float32(1), isError: true},
		{left: uint(1), right: float64(1), isError: true},
		{left: uint(1), right: "foo", isError: true},
		{left: uint(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint(1), right: []string{"foo"}, isError: true},

		// left is uint8
		{left: uint8(1), right: true, isError: true},
		{left: uint8(1), right: int(1), expect: true},
		{left: uint8(1), right: int8(1), expect: true},
		{left: uint8(1), right: int16(1), expect: true},
		{left: uint8(1), right: int32(1), expect: true},
		{left: uint8(1), right: int64(1), expect: true},
		{left: uint8(1), right: uint(1), expect: true},
		{left: uint8(1), right: uint8(1), expect: true},
		{left: uint8(1), right: uint16(1), expect: true},
		{left: uint8(1), right: uint32(1), expect: true},
		{left: uint8(1), right: uint64(1), expect: true},
		{left: uint8(1), right: float32(1), isError: true},
		{left: uint8(1), right: float64(1), isError: true},
		{left: uint8(1), right: "foo", isError: true},
		{left: uint8(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint8(1), right: []string{"foo"}, isError: true},

		// left is uint16
		{left: uint16(1), right: true, isError: true},
		{left: uint16(1), right: int(1), expect: true},
		{left: uint16(1), right: int8(1), expect: true},
		{left: uint16(1), right: int16(1), expect: true},
		{left: uint16(1), right: int32(1), expect: true},
		{left: uint16(1), right: int64(1), expect: true},
		{left: uint16(1), right: uint(1), expect: true},
		{left: uint16(1), right: uint8(1), expect: true},
		{left: uint16(1), right: uint16(1), expect: true},
		{left: uint16(1), right: uint32(1), expect: true},
		{left: uint16(1), right: uint64(1), expect: true},
		{left: uint16(1), right: float32(1), isError: true},
		{left: uint16(1), right: float64(1), isError: true},
		{left: uint16(1), right: "foo", isError: true},
		{left: uint16(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint16(1), right: []string{"foo"}, isError: true},

		// left is uint32
		{left: uint32(1), right: true, isError: true},
		{left: uint32(1), right: int(1), expect: true},
		{left: uint32(1), right: int8(1), expect: true},
		{left: uint32(1), right: int16(1), expect: true},
		{left: uint32(1), right: int32(1), expect: true},
		{left: uint32(1), right: int64(1), expect: true},
		{left: uint32(1), right: uint(1), expect: true},
		{left: uint32(1), right: uint8(1), expect: true},
		{left: uint32(1), right: uint16(1), expect: true},
		{left: uint32(1), right: uint32(1), expect: true},
		{left: uint32(1), right: uint64(1), expect: true},
		{left: uint32(1), right: float32(1), isError: true},
		{left: uint32(1), right: float64(1), isError: true},
		{left: uint32(1), right: "foo", isError: true},
		{left: uint32(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint32(1), right: []string{"foo"}, isError: true},

		// left is uint64
		{left: uint64(1), right: true, isError: true},
		{left: uint64(1), right: int(1), expect: true},
		{left: uint64(1), right: int8(1), expect: true},
		{left: uint64(1), right: int16(1), expect: true},
		{left: uint64(1), right: int32(1), expect: true},
		{left: uint64(1), right: int64(1), expect: true},
		{left: uint64(1), right: uint(1), expect: true},
		{left: uint64(1), right: uint8(1), expect: true},
		{left: uint64(1), right: uint16(1), expect: true},
		{left: uint64(1), right: uint32(1), expect: true},
		{left: uint64(1), right: uint64(1), expect: true},
		{left: uint64(1), right: float32(1), isError: true},
		{left: uint64(1), right: float64(1), isError: true},
		{left: uint64(1), right: "foo", isError: true},
		{left: uint64(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint64(1), right: []string{"foo"}, isError: true},

		// left is float32
		{left: float32(1), right: true, isError: true},
		{left: float32(1), right: int(1), isError: true},
		{left: float32(1), right: int8(1), isError: true},
		{left: float32(1), right: int16(1), isError: true},
		{left: float32(1), right: int32(1), isError: true},
		{left: float32(1), right: int64(1), isError: true},
		{left: float32(1), right: uint(1), isError: true},
		{left: float32(1), right: uint8(1), isError: true},
		{left: float32(1), right: uint16(1), isError: true},
		{left: float32(1), right: uint32(1), isError: true},
		{left: float32(1), right: uint64(1), isError: true},
		{left: float32(1), right: float32(1), expect: true},
		{left: float32(1), right: float64(1), expect: true},
		{left: float32(1), right: "foo", isError: true},
		{left: float32(1), right: testStruct{name: "foo"}, isError: true},
		{left: float32(1), right: []string{"foo"}, isError: true},

		// left is float64
		{left: float64(1), right: true, isError: true},
		{left: float64(1), right: int(1), isError: true},
		{left: float64(1), right: int8(1), isError: true},
		{left: float64(1), right: int16(1), isError: true},
		{left: float64(1), right: int32(1), isError: true},
		{left: float64(1), right: int64(1), isError: true},
		{left: float64(1), right: uint(1), isError: true},
		{left: float64(1), right: uint8(1), isError: true},
		{left: float64(1), right: uint16(1), isError: true},
		{left: float64(1), right: uint32(1), isError: true},
		{left: float64(1), right: uint64(1), isError: true},
		{left: float64(1), right: float32(1), expect: true},
		{left: float64(1), right: float64(1), expect: true},
		{left: float64(1), right: "foo", isError: true},
		{left: float64(1), right: testStruct{name: "foo"}, isError: true},
		{left: float64(1), right: []string{"foo"}, isError: true},

		// left is string
		{left: "foo", right: true, isError: true},
		{left: "foo", right: int(1), isError: true},
		{left: "foo", right: int8(1), isError: true},
		{left: "foo", right: int16(1), isError: true},
		{left: "foo", right: int32(1), isError: true},
		{left: "foo", right: int64(1), isError: true},
		{left: "foo", right: uint(1), isError: true},
		{left: "foo", right: uint8(1), isError: true},
		{left: "foo", right: uint16(1), isError: true},
		{left: "foo", right: uint32(1), isError: true},
		{left: "foo", right: uint64(1), isError: true},
		{left: "foo", right: float32(1), isError: true},
		{left: "foo", right: float64(1), isError: true},
		{left: "foo", right: "foo", expect: true},
		{left: "foo", right: testStruct{name: "foo"}, isError: true},
		{left: "foo", right: []string{"foo"}, isError: true},

		// left is struct
		{left: testStruct{name: "foo"}, right: true, isError: true},
		{left: testStruct{name: "foo"}, right: int(1), isError: true},
		{left: testStruct{name: "foo"}, right: int8(1), isError: true},
		{left: testStruct{name: "foo"}, right: int16(1), isError: true},
		{left: testStruct{name: "foo"}, right: int32(1), isError: true},
		{left: testStruct{name: "foo"}, right: int64(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint8(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint16(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint32(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint64(1), isError: true},
		{left: testStruct{name: "foo"}, right: float32(1), isError: true},
		{left: testStruct{name: "foo"}, right: float64(1), isError: true},
		{left: testStruct{name: "foo"}, right: "foo", isError: true},
		{left: testStruct{name: "foo"}, right: testStruct{name: "foo"}, expect: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "foo"}, expect: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "bar"}, expect: false},
		{left: testStruct{name: "foo"}, right: []string{"foo"}, isError: true},
	}

	for i, tt := range tests {
		ret, err := Equal(reflect.ValueOf(tt.left), reflect.ValueOf(tt.right))

		if tt.isError {
			if err == nil {
				t.Errorf("[%d] Expects error, got nil", i)
			}
			return
		}
		if err != nil {
			t.Errorf("[%d] Expects no error, got error %s", i, err)
			return
		}

		if diff := cmp.Diff(tt.expect, ret); diff != "" {
			t.Errorf("[%d], Equal result mismatch, diff=%s", i, diff)
		}
	}
}

func TestCompareValuesNotEqual(t *testing.T) {
	tests := []struct {
		left    any
		right   any
		expect  bool
		isError bool
	}{
		// left is bool
		{left: true, right: true, expect: false},
		{left: true, right: int(1), isError: true},
		{left: true, right: int8(1), isError: true},
		{left: true, right: int16(1), isError: true},
		{left: true, right: int32(1), isError: true},
		{left: true, right: int64(1), isError: true},
		{left: true, right: uint(1), isError: true},
		{left: true, right: uint8(1), isError: true},
		{left: true, right: uint16(1), isError: true},
		{left: true, right: uint32(1), isError: true},
		{left: true, right: uint64(1), isError: true},
		{left: true, right: float32(1), isError: true},
		{left: true, right: float64(1), isError: true},
		{left: true, right: "foo", isError: true},
		{left: true, right: testStruct{name: "foo"}, isError: true},
		{left: true, right: []string{"foo"}, isError: true},

		// left is int
		{left: int(1), right: true, isError: true},
		{left: int(1), right: int(1), expect: false},
		{left: int(1), right: int8(1), expect: false},
		{left: int(1), right: int16(1), expect: false},
		{left: int(1), right: int32(1), expect: false},
		{left: int(1), right: int64(1), expect: false},
		{left: int(1), right: uint(1), expect: false},
		{left: int(1), right: uint8(1), expect: false},
		{left: int(1), right: uint16(1), expect: false},
		{left: int(1), right: uint32(1), expect: false},
		{left: int(1), right: uint64(1), expect: false},
		{left: int(1), right: float32(1), isError: true},
		{left: int(1), right: float64(1), isError: true},
		{left: int(1), right: "foo", isError: true},
		{left: int(1), right: testStruct{name: "foo"}, isError: true},
		{left: int(1), right: []string{"foo"}, isError: true},

		// left is int8
		{left: int8(1), right: true, isError: true},
		{left: int8(1), right: int(1), expect: false},
		{left: int8(1), right: int8(1), expect: false},
		{left: int8(1), right: int16(1), expect: false},
		{left: int8(1), right: int32(1), expect: false},
		{left: int8(1), right: int64(1), expect: false},
		{left: int8(1), right: uint(1), expect: false},
		{left: int8(1), right: uint8(1), expect: false},
		{left: int8(1), right: uint16(1), expect: false},
		{left: int8(1), right: uint32(1), expect: false},
		{left: int8(1), right: uint64(1), expect: false},
		{left: int8(1), right: float32(1), isError: true},
		{left: int8(1), right: float64(1), isError: true},
		{left: int8(1), right: "foo", isError: true},
		{left: int8(1), right: testStruct{name: "foo"}, isError: true},
		{left: int8(1), right: []string{"foo"}, isError: true},

		// left is int16
		{left: int16(1), right: true, isError: true},
		{left: int16(1), right: int(1), expect: false},
		{left: int16(1), right: int8(1), expect: false},
		{left: int16(1), right: int16(1), expect: false},
		{left: int16(1), right: int32(1), expect: false},
		{left: int16(1), right: int64(1), expect: false},
		{left: int16(1), right: uint(1), expect: false},
		{left: int16(1), right: uint8(1), expect: false},
		{left: int16(1), right: uint16(1), expect: false},
		{left: int16(1), right: uint32(1), expect: false},
		{left: int16(1), right: uint64(1), expect: false},
		{left: int16(1), right: float32(1), isError: true},
		{left: int16(1), right: float64(1), isError: true},
		{left: int16(1), right: "foo", isError: true},
		{left: int16(1), right: testStruct{name: "foo"}, isError: true},
		{left: int16(1), right: []string{"foo"}, isError: true},

		// left is int32
		{left: int32(1), right: true, isError: true},
		{left: int32(1), right: int(1), expect: false},
		{left: int32(1), right: int8(1), expect: false},
		{left: int32(1), right: int16(1), expect: false},
		{left: int32(1), right: int32(1), expect: false},
		{left: int32(1), right: int64(1), expect: false},
		{left: int32(1), right: uint(1), expect: false},
		{left: int32(1), right: uint8(1), expect: false},
		{left: int32(1), right: uint16(1), expect: false},
		{left: int32(1), right: uint32(1), expect: false},
		{left: int32(1), right: uint64(1), expect: false},
		{left: int32(1), right: float32(1), isError: true},
		{left: int32(1), right: float64(1), isError: true},
		{left: int32(1), right: "foo", isError: true},
		{left: int32(1), right: testStruct{name: "foo"}, isError: true},
		{left: int32(1), right: []string{"foo"}, isError: true},

		// left is int64
		{left: int64(1), right: true, isError: true},
		{left: int64(1), right: int(1), expect: false},
		{left: int64(1), right: int8(1), expect: false},
		{left: int64(1), right: int16(1), expect: false},
		{left: int64(1), right: int32(1), expect: false},
		{left: int64(1), right: int64(1), expect: false},
		{left: int64(1), right: uint(1), expect: false},
		{left: int64(1), right: uint8(1), expect: false},
		{left: int64(1), right: uint16(1), expect: false},
		{left: int64(1), right: uint32(1), expect: false},
		{left: int64(1), right: uint64(1), expect: false},
		{left: int64(1), right: float32(1), isError: true},
		{left: int64(1), right: float64(1), isError: true},
		{left: int64(1), right: "foo", isError: true},
		{left: int64(1), right: testStruct{name: "foo"}, isError: true},
		{left: int64(1), right: []string{"foo"}, isError: true},

		// left is uint
		{left: uint(1), right: true, isError: true},
		{left: uint(1), right: int(1), expect: false},
		{left: uint(1), right: int8(1), expect: false},
		{left: uint(1), right: int16(1), expect: false},
		{left: uint(1), right: int32(1), expect: false},
		{left: uint(1), right: int64(1), expect: false},
		{left: uint(1), right: uint(1), expect: false},
		{left: uint(1), right: uint8(1), expect: false},
		{left: uint(1), right: uint16(1), expect: false},
		{left: uint(1), right: uint32(1), expect: false},
		{left: uint(1), right: uint64(1), expect: false},
		{left: uint(1), right: float32(1), isError: true},
		{left: uint(1), right: float64(1), isError: true},
		{left: uint(1), right: "foo", isError: true},
		{left: uint(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint(1), right: []string{"foo"}, isError: true},

		// left is uint8
		{left: uint8(1), right: true, isError: true},
		{left: uint8(1), right: int(1), expect: false},
		{left: uint8(1), right: int8(1), expect: false},
		{left: uint8(1), right: int16(1), expect: false},
		{left: uint8(1), right: int32(1), expect: false},
		{left: uint8(1), right: int64(1), expect: false},
		{left: uint8(1), right: uint(1), expect: false},
		{left: uint8(1), right: uint8(1), expect: false},
		{left: uint8(1), right: uint16(1), expect: false},
		{left: uint8(1), right: uint32(1), expect: false},
		{left: uint8(1), right: uint64(1), expect: false},
		{left: uint8(1), right: float32(1), isError: true},
		{left: uint8(1), right: float64(1), isError: true},
		{left: uint8(1), right: "foo", isError: true},
		{left: uint8(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint8(1), right: []string{"foo"}, isError: true},

		// left is uint16
		{left: uint16(1), right: true, isError: true},
		{left: uint16(1), right: int(1), expect: false},
		{left: uint16(1), right: int8(1), expect: false},
		{left: uint16(1), right: int16(1), expect: false},
		{left: uint16(1), right: int32(1), expect: false},
		{left: uint16(1), right: int64(1), expect: false},
		{left: uint16(1), right: uint(1), expect: false},
		{left: uint16(1), right: uint8(1), expect: false},
		{left: uint16(1), right: uint16(1), expect: false},
		{left: uint16(1), right: uint32(1), expect: false},
		{left: uint16(1), right: uint64(1), expect: false},
		{left: uint16(1), right: float32(1), isError: true},
		{left: uint16(1), right: float64(1), isError: true},
		{left: uint16(1), right: "foo", isError: true},
		{left: uint16(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint16(1), right: []string{"foo"}, isError: true},

		// left is uint32
		{left: uint32(1), right: true, isError: true},
		{left: uint32(1), right: int(1), expect: false},
		{left: uint32(1), right: int8(1), expect: false},
		{left: uint32(1), right: int16(1), expect: false},
		{left: uint32(1), right: int32(1), expect: false},
		{left: uint32(1), right: int64(1), expect: false},
		{left: uint32(1), right: uint(1), expect: false},
		{left: uint32(1), right: uint8(1), expect: false},
		{left: uint32(1), right: uint16(1), expect: false},
		{left: uint32(1), right: uint32(1), expect: false},
		{left: uint32(1), right: uint64(1), expect: false},
		{left: uint32(1), right: float32(1), isError: true},
		{left: uint32(1), right: float64(1), isError: true},
		{left: uint32(1), right: "foo", isError: true},
		{left: uint32(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint32(1), right: []string{"foo"}, isError: true},

		// left is uint64
		{left: uint64(1), right: true, isError: true},
		{left: uint64(1), right: int(1), expect: false},
		{left: uint64(1), right: int8(1), expect: false},
		{left: uint64(1), right: int16(1), expect: false},
		{left: uint64(1), right: int32(1), expect: false},
		{left: uint64(1), right: int64(1), expect: false},
		{left: uint64(1), right: uint(1), expect: false},
		{left: uint64(1), right: uint8(1), expect: false},
		{left: uint64(1), right: uint16(1), expect: false},
		{left: uint64(1), right: uint32(1), expect: false},
		{left: uint64(1), right: uint64(1), expect: false},
		{left: uint64(1), right: float32(1), isError: true},
		{left: uint64(1), right: float64(1), isError: true},
		{left: uint64(1), right: "foo", isError: true},
		{left: uint64(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint64(1), right: []string{"foo"}, isError: true},

		// left is float32
		{left: float32(1), right: true, isError: true},
		{left: float32(1), right: int(1), isError: true},
		{left: float32(1), right: int8(1), isError: true},
		{left: float32(1), right: int16(1), isError: true},
		{left: float32(1), right: int32(1), isError: true},
		{left: float32(1), right: int64(1), isError: true},
		{left: float32(1), right: uint(1), isError: true},
		{left: float32(1), right: uint8(1), isError: true},
		{left: float32(1), right: uint16(1), isError: true},
		{left: float32(1), right: uint32(1), isError: true},
		{left: float32(1), right: uint64(1), isError: true},
		{left: float32(1), right: float32(1), expect: false},
		{left: float32(1), right: float64(1), expect: false},
		{left: float32(1), right: "foo", isError: true},
		{left: float32(1), right: testStruct{name: "foo"}, isError: true},
		{left: float32(1), right: []string{"foo"}, isError: true},

		// left is float64
		{left: float64(1), right: true, isError: true},
		{left: float64(1), right: int(1), isError: true},
		{left: float64(1), right: int8(1), isError: true},
		{left: float64(1), right: int16(1), isError: true},
		{left: float64(1), right: int32(1), isError: true},
		{left: float64(1), right: int64(1), isError: true},
		{left: float64(1), right: uint(1), isError: true},
		{left: float64(1), right: uint8(1), isError: true},
		{left: float64(1), right: uint16(1), isError: true},
		{left: float64(1), right: uint32(1), isError: true},
		{left: float64(1), right: uint64(1), isError: true},
		{left: float64(1), right: float32(1), expect: false},
		{left: float64(1), right: float64(1), expect: false},
		{left: float64(1), right: "foo", isError: true},
		{left: float64(1), right: testStruct{name: "foo"}, isError: true},
		{left: float64(1), right: []string{"foo"}, isError: true},

		// left is string
		{left: "foo", right: true, isError: true},
		{left: "foo", right: int(1), isError: true},
		{left: "foo", right: int8(1), isError: true},
		{left: "foo", right: int16(1), isError: true},
		{left: "foo", right: int32(1), isError: true},
		{left: "foo", right: int64(1), isError: true},
		{left: "foo", right: uint(1), isError: true},
		{left: "foo", right: uint8(1), isError: true},
		{left: "foo", right: uint16(1), isError: true},
		{left: "foo", right: uint32(1), isError: true},
		{left: "foo", right: uint64(1), isError: true},
		{left: "foo", right: float32(1), isError: true},
		{left: "foo", right: float64(1), isError: true},
		{left: "foo", right: "foo", expect: false},
		{left: "foo", right: testStruct{name: "foo"}, isError: true},
		{left: "foo", right: []string{"foo"}, isError: true},

		// left is struct
		{left: testStruct{name: "foo"}, right: true, isError: true},
		{left: testStruct{name: "foo"}, right: int(1), isError: true},
		{left: testStruct{name: "foo"}, right: int8(1), isError: true},
		{left: testStruct{name: "foo"}, right: int16(1), isError: true},
		{left: testStruct{name: "foo"}, right: int32(1), isError: true},
		{left: testStruct{name: "foo"}, right: int64(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint8(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint16(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint32(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint64(1), isError: true},
		{left: testStruct{name: "foo"}, right: float32(1), isError: true},
		{left: testStruct{name: "foo"}, right: float64(1), isError: true},
		{left: testStruct{name: "foo"}, right: "foo", isError: true},
		{left: testStruct{name: "foo"}, right: testStruct{name: "foo"}, expect: false},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "foo"}, expect: false},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "bar"}, expect: true},
		{left: testStruct{name: "foo"}, right: []string{"foo"}, isError: true},
	}

	for i, tt := range tests {
		ret, err := NotEqual(reflect.ValueOf(tt.left), reflect.ValueOf(tt.right))

		if tt.isError {
			if err == nil {
				t.Errorf("[%d] Expects error, got nil, %v", i, tt)
			}
			return
		}
		if err != nil {
			t.Errorf("[%d] Expects no error, got error %s", i, err)
			return
		}

		if diff := cmp.Diff(tt.expect, ret); diff != "" {
			t.Errorf("[%d] NotEqual result mismatch, diff=%s", i, diff)
		}
	}
}

func TestCompareValuesGreaterThan(t *testing.T) {
	tests := []struct {
		left    any
		right   any
		expect  bool
		isError bool
	}{
		// left is bool
		{left: true, right: true, isError: true},
		{left: true, right: int(1), isError: true},
		{left: true, right: int8(1), isError: true},
		{left: true, right: int16(1), isError: true},
		{left: true, right: int32(1), isError: true},
		{left: true, right: int64(1), isError: true},
		{left: true, right: uint(1), isError: true},
		{left: true, right: uint8(1), isError: true},
		{left: true, right: uint16(1), isError: true},
		{left: true, right: uint32(1), isError: true},
		{left: true, right: uint64(1), isError: true},
		{left: true, right: float32(1), isError: true},
		{left: true, right: float64(1), isError: true},
		{left: true, right: "foo", isError: true},
		{left: true, right: testStruct{name: "foo"}, isError: true},
		{left: true, right: []string{"foo"}, isError: true},

		// left is int
		{left: int(1), right: true, isError: true},
		{left: int(2), right: int(1), expect: true},
		{left: int(1), right: int(1), expect: false},
		{left: int(1), right: int8(1), expect: false},
		{left: int(1), right: int16(1), expect: false},
		{left: int(1), right: int32(1), expect: false},
		{left: int(1), right: int64(1), expect: false},
		{left: int(1), right: uint(1), expect: false},
		{left: int(1), right: uint8(1), expect: false},
		{left: int(1), right: uint16(1), expect: false},
		{left: int(1), right: uint32(1), expect: false},
		{left: int(1), right: uint64(1), expect: false},
		{left: int(1), right: float32(1), isError: true},
		{left: int(1), right: float64(1), isError: true},
		{left: int(1), right: "foo", isError: true},
		{left: int(1), right: testStruct{name: "foo"}, isError: true},
		{left: int(1), right: []string{"foo"}, isError: true},

		// left is int8
		{left: int8(1), right: true, isError: true},
		{left: int8(2), right: int(1), expect: true},
		{left: int8(1), right: int(1), expect: false},
		{left: int8(1), right: int8(1), expect: false},
		{left: int8(1), right: int16(1), expect: false},
		{left: int8(1), right: int32(1), expect: false},
		{left: int8(1), right: int64(1), expect: false},
		{left: int8(1), right: uint(1), expect: false},
		{left: int8(1), right: uint8(1), expect: false},
		{left: int8(1), right: uint16(1), expect: false},
		{left: int8(1), right: uint32(1), expect: false},
		{left: int8(1), right: uint64(1), expect: false},
		{left: int8(1), right: float32(1), isError: true},
		{left: int8(1), right: float64(1), isError: true},
		{left: int8(1), right: "foo", isError: true},
		{left: int8(1), right: testStruct{name: "foo"}, isError: true},
		{left: int8(1), right: []string{"foo"}, isError: true},

		// left is int16
		{left: int16(1), right: true, isError: true},
		{left: int16(2), right: int(1), expect: true},
		{left: int16(1), right: int(1), expect: false},
		{left: int16(1), right: int8(1), expect: false},
		{left: int16(1), right: int16(1), expect: false},
		{left: int16(1), right: int32(1), expect: false},
		{left: int16(1), right: int64(1), expect: false},
		{left: int16(1), right: uint(1), expect: false},
		{left: int16(1), right: uint8(1), expect: false},
		{left: int16(1), right: uint16(1), expect: false},
		{left: int16(1), right: uint32(1), expect: false},
		{left: int16(1), right: uint64(1), expect: false},
		{left: int16(1), right: float32(1), isError: true},
		{left: int16(1), right: float64(1), isError: true},
		{left: int16(1), right: "foo", isError: true},
		{left: int16(1), right: testStruct{name: "foo"}, isError: true},
		{left: int16(1), right: []string{"foo"}, isError: true},

		// left is int32
		{left: int32(1), right: true, isError: true},
		{left: int32(2), right: int(1), expect: true},
		{left: int32(1), right: int(1), expect: false},
		{left: int32(1), right: int8(1), expect: false},
		{left: int32(1), right: int16(1), expect: false},
		{left: int32(1), right: int32(1), expect: false},
		{left: int32(1), right: int64(1), expect: false},
		{left: int32(1), right: uint(1), expect: false},
		{left: int32(1), right: uint8(1), expect: false},
		{left: int32(1), right: uint16(1), expect: false},
		{left: int32(1), right: uint32(1), expect: false},
		{left: int32(1), right: uint64(1), expect: false},
		{left: int32(1), right: float32(1), isError: true},
		{left: int32(1), right: float64(1), isError: true},
		{left: int32(1), right: "foo", isError: true},
		{left: int32(1), right: testStruct{name: "foo"}, isError: true},
		{left: int32(1), right: []string{"foo"}, isError: true},

		// left is int64
		{left: int64(1), right: true, isError: true},
		{left: int64(2), right: int(1), expect: true},
		{left: int64(1), right: int(1), expect: false},
		{left: int64(1), right: int8(1), expect: false},
		{left: int64(1), right: int16(1), expect: false},
		{left: int64(1), right: int32(1), expect: false},
		{left: int64(1), right: int64(1), expect: false},
		{left: int64(1), right: uint(1), expect: false},
		{left: int64(1), right: uint8(1), expect: false},
		{left: int64(1), right: uint16(1), expect: false},
		{left: int64(1), right: uint32(1), expect: false},
		{left: int64(1), right: uint64(1), expect: false},
		{left: int64(1), right: float32(1), isError: true},
		{left: int64(1), right: float64(1), isError: true},
		{left: int64(1), right: "foo", isError: true},
		{left: int64(1), right: testStruct{name: "foo"}, isError: true},
		{left: int64(1), right: []string{"foo"}, isError: true},

		// left is uint
		{left: uint(1), right: true, isError: true},
		{left: uint(2), right: int(1), expect: true},
		{left: uint(1), right: int(1), expect: false},
		{left: uint(1), right: int8(1), expect: false},
		{left: uint(1), right: int16(1), expect: false},
		{left: uint(1), right: int32(1), expect: false},
		{left: uint(1), right: int64(1), expect: false},
		{left: uint(1), right: uint(1), expect: false},
		{left: uint(1), right: uint8(1), expect: false},
		{left: uint(1), right: uint16(1), expect: false},
		{left: uint(1), right: uint32(1), expect: false},
		{left: uint(1), right: uint64(1), expect: false},
		{left: uint(1), right: float32(1), isError: true},
		{left: uint(1), right: float64(1), isError: true},
		{left: uint(1), right: "foo", isError: true},
		{left: uint(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint(1), right: []string{"foo"}, isError: true},

		// left is uint8
		{left: uint8(1), right: true, isError: true},
		{left: uint8(2), right: int(1), expect: true},
		{left: uint8(1), right: int(1), expect: false},
		{left: uint8(1), right: int8(1), expect: false},
		{left: uint8(1), right: int16(1), expect: false},
		{left: uint8(1), right: int32(1), expect: false},
		{left: uint8(1), right: int64(1), expect: false},
		{left: uint8(1), right: uint(1), expect: false},
		{left: uint8(1), right: uint8(1), expect: false},
		{left: uint8(1), right: uint16(1), expect: false},
		{left: uint8(1), right: uint32(1), expect: false},
		{left: uint8(1), right: uint64(1), expect: false},
		{left: uint8(1), right: float32(1), isError: true},
		{left: uint8(1), right: float64(1), isError: true},
		{left: uint8(1), right: "foo", isError: true},
		{left: uint8(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint8(1), right: []string{"foo"}, isError: true},

		// left is uint16
		{left: uint16(1), right: true, isError: true},
		{left: uint16(2), right: int(1), expect: true},
		{left: uint16(1), right: int(1), expect: false},
		{left: uint16(1), right: int8(1), expect: false},
		{left: uint16(1), right: int16(1), expect: false},
		{left: uint16(1), right: int32(1), expect: false},
		{left: uint16(1), right: int64(1), expect: false},
		{left: uint16(1), right: uint(1), expect: false},
		{left: uint16(1), right: uint8(1), expect: false},
		{left: uint16(1), right: uint16(1), expect: false},
		{left: uint16(1), right: uint32(1), expect: false},
		{left: uint16(1), right: uint64(1), expect: false},
		{left: uint16(1), right: float32(1), isError: true},
		{left: uint16(1), right: float64(1), isError: true},
		{left: uint16(1), right: "foo", isError: true},
		{left: uint16(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint16(1), right: []string{"foo"}, isError: true},

		// left is uint32
		{left: uint32(1), right: true, isError: true},
		{left: uint32(2), right: int(1), expect: true},
		{left: uint32(1), right: int(1), expect: false},
		{left: uint32(1), right: int8(1), expect: false},
		{left: uint32(1), right: int16(1), expect: false},
		{left: uint32(1), right: int32(1), expect: false},
		{left: uint32(1), right: int64(1), expect: false},
		{left: uint32(1), right: uint(1), expect: false},
		{left: uint32(1), right: uint8(1), expect: false},
		{left: uint32(1), right: uint16(1), expect: false},
		{left: uint32(1), right: uint32(1), expect: false},
		{left: uint32(1), right: uint64(1), expect: false},
		{left: uint32(1), right: float32(1), isError: true},
		{left: uint32(1), right: float64(1), isError: true},
		{left: uint32(1), right: "foo", isError: true},
		{left: uint32(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint32(1), right: []string{"foo"}, isError: true},

		// left is uint64
		{left: uint64(1), right: true, isError: true},
		{left: uint64(1), right: int(1), expect: false},
		{left: uint64(1), right: int8(1), expect: false},
		{left: uint64(1), right: int16(1), expect: false},
		{left: uint64(1), right: int32(1), expect: false},
		{left: uint64(1), right: int64(1), expect: false},
		{left: uint64(1), right: uint(1), expect: false},
		{left: uint64(1), right: uint8(1), expect: false},
		{left: uint64(1), right: uint16(1), expect: false},
		{left: uint64(1), right: uint32(1), expect: false},
		{left: uint64(1), right: uint64(1), expect: false},
		{left: uint64(1), right: float32(1), isError: true},
		{left: uint64(1), right: float64(1), isError: true},
		{left: uint64(1), right: "foo", isError: true},
		{left: uint64(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint64(1), right: []string{"foo"}, isError: true},

		// left is float32
		{left: float32(1), right: true, isError: true},
		{left: float32(1), right: int(1), isError: true},
		{left: float32(1), right: int8(1), isError: true},
		{left: float32(1), right: int16(1), isError: true},
		{left: float32(1), right: int32(1), isError: true},
		{left: float32(1), right: int64(1), isError: true},
		{left: float32(1), right: uint(1), isError: true},
		{left: float32(1), right: uint8(1), isError: true},
		{left: float32(1), right: uint16(1), isError: true},
		{left: float32(1), right: uint32(1), isError: true},
		{left: float32(1), right: uint64(1), isError: true},
		{left: float32(1), right: float32(1), expect: false},
		{left: float32(1), right: float64(1), expect: false},
		{left: float64(1.1), right: float64(1), expect: true},
		{left: float32(1), right: "foo", isError: true},
		{left: float32(1), right: testStruct{name: "foo"}, isError: true},
		{left: float32(1), right: []string{"foo"}, isError: true},

		// left is float64
		{left: float64(1), right: true, isError: true},
		{left: float64(1), right: int(1), isError: true},
		{left: float64(1), right: int8(1), isError: true},
		{left: float64(1), right: int16(1), isError: true},
		{left: float64(1), right: int32(1), isError: true},
		{left: float64(1), right: int64(1), isError: true},
		{left: float64(1), right: uint(1), isError: true},
		{left: float64(1), right: uint8(1), isError: true},
		{left: float64(1), right: uint16(1), isError: true},
		{left: float64(1), right: uint32(1), isError: true},
		{left: float64(1), right: uint64(1), isError: true},
		{left: float64(1), right: float32(1), expect: false},
		{left: float64(1), right: float64(1), expect: false},
		{left: float64(1.1), right: float64(1), expect: true},
		{left: float64(1), right: "foo", isError: true},
		{left: float64(1), right: testStruct{name: "foo"}, isError: true},
		{left: float64(1), right: []string{"foo"}, isError: true},

		// left is string
		{left: "foo", right: true, isError: true},
		{left: "foo", right: int(1), isError: true},
		{left: "foo", right: int8(1), isError: true},
		{left: "foo", right: int16(1), isError: true},
		{left: "foo", right: int32(1), isError: true},
		{left: "foo", right: int64(1), isError: true},
		{left: "foo", right: uint(1), isError: true},
		{left: "foo", right: uint8(1), isError: true},
		{left: "foo", right: uint16(1), isError: true},
		{left: "foo", right: uint32(1), isError: true},
		{left: "foo", right: uint64(1), isError: true},
		{left: "foo", right: float32(1), isError: true},
		{left: "foo", right: float64(1), isError: true},
		{left: "foo", right: "foo", isError: true},
		{left: "foo", right: testStruct{name: "foo"}, isError: true},
		{left: "foo", right: []string{"foo"}, isError: true},

		// left is struct
		{left: testStruct{name: "foo"}, right: true, isError: true},
		{left: testStruct{name: "foo"}, right: int(1), isError: true},
		{left: testStruct{name: "foo"}, right: int8(1), isError: true},
		{left: testStruct{name: "foo"}, right: int16(1), isError: true},
		{left: testStruct{name: "foo"}, right: int32(1), isError: true},
		{left: testStruct{name: "foo"}, right: int64(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint8(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint16(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint32(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint64(1), isError: true},
		{left: testStruct{name: "foo"}, right: float32(1), isError: true},
		{left: testStruct{name: "foo"}, right: float64(1), isError: true},
		{left: testStruct{name: "foo"}, right: "foo", isError: true},
		{left: testStruct{name: "foo"}, right: testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "bar"}, isError: true},
		{left: testStruct{name: "foo"}, right: []string{"foo"}, isError: true},
	}

	for i, tt := range tests {
		ret, err := GreaterThan(reflect.ValueOf(tt.left), reflect.ValueOf(tt.right))

		if tt.isError {
			if err == nil {
				t.Errorf("[%d] Expects error, got nil, %v", i, tt)
			}
			return
		}
		if err != nil {
			t.Errorf("[%d] Expects no error, got error %s", i, err)
			return
		}

		if diff := cmp.Diff(tt.expect, ret); diff != "" {
			t.Errorf("[%d] GreaterThan result mismatch, diff=%s", i, diff)
		}
	}
}

func TestCompareValuesGreaterThanEqual(t *testing.T) {
	tests := []struct {
		left    any
		right   any
		expect  bool
		isError bool
	}{
		// left is bool
		{left: true, right: true, isError: true},
		{left: true, right: int(1), isError: true},
		{left: true, right: int8(1), isError: true},
		{left: true, right: int16(1), isError: true},
		{left: true, right: int32(1), isError: true},
		{left: true, right: int64(1), isError: true},
		{left: true, right: uint(1), isError: true},
		{left: true, right: uint8(1), isError: true},
		{left: true, right: uint16(1), isError: true},
		{left: true, right: uint32(1), isError: true},
		{left: true, right: uint64(1), isError: true},
		{left: true, right: float32(1), isError: true},
		{left: true, right: float64(1), isError: true},
		{left: true, right: "foo", isError: true},
		{left: true, right: testStruct{name: "foo"}, isError: true},
		{left: true, right: []string{"foo"}, isError: true},

		// left is int
		{left: int(1), right: true, isError: true},
		{left: int(2), right: int(1), expect: true},
		{left: int(1), right: int(1), expect: true},
		{left: int(1), right: int8(1), expect: true},
		{left: int(1), right: int16(1), expect: true},
		{left: int(1), right: int32(1), expect: true},
		{left: int(1), right: int64(1), expect: true},
		{left: int(1), right: uint(1), expect: true},
		{left: int(1), right: uint8(1), expect: true},
		{left: int(1), right: uint16(1), expect: true},
		{left: int(1), right: uint32(1), expect: true},
		{left: int(1), right: uint64(1), expect: true},
		{left: int(1), right: float32(1), isError: true},
		{left: int(1), right: float64(1), isError: true},
		{left: int(1), right: "foo", isError: true},
		{left: int(1), right: testStruct{name: "foo"}, isError: true},
		{left: int(1), right: []string{"foo"}, isError: true},

		// left is int8
		{left: int8(1), right: true, isError: true},
		{left: int8(2), right: int(1), expect: true},
		{left: int8(1), right: int(1), expect: true},
		{left: int8(1), right: int8(1), expect: true},
		{left: int8(1), right: int16(1), expect: true},
		{left: int8(1), right: int32(1), expect: true},
		{left: int8(1), right: int64(1), expect: true},
		{left: int8(1), right: uint(1), expect: true},
		{left: int8(1), right: uint8(1), expect: true},
		{left: int8(1), right: uint16(1), expect: true},
		{left: int8(1), right: uint32(1), expect: true},
		{left: int8(1), right: uint64(1), expect: true},
		{left: int8(1), right: float32(1), isError: true},
		{left: int8(1), right: float64(1), isError: true},
		{left: int8(1), right: "foo", isError: true},
		{left: int8(1), right: testStruct{name: "foo"}, isError: true},
		{left: int8(1), right: []string{"foo"}, isError: true},

		// left is int16
		{left: int16(1), right: true, isError: true},
		{left: int16(2), right: int(1), expect: true},
		{left: int16(1), right: int(1), expect: true},
		{left: int16(1), right: int8(1), expect: true},
		{left: int16(1), right: int16(1), expect: true},
		{left: int16(1), right: int32(1), expect: true},
		{left: int16(1), right: int64(1), expect: true},
		{left: int16(1), right: uint(1), expect: true},
		{left: int16(1), right: uint8(1), expect: true},
		{left: int16(1), right: uint16(1), expect: true},
		{left: int16(1), right: uint32(1), expect: true},
		{left: int16(1), right: uint64(1), expect: true},
		{left: int16(1), right: float32(1), isError: true},
		{left: int16(1), right: float64(1), isError: true},
		{left: int16(1), right: "foo", isError: true},
		{left: int16(1), right: testStruct{name: "foo"}, isError: true},
		{left: int16(1), right: []string{"foo"}, isError: true},

		// left is int32
		{left: int32(1), right: true, isError: true},
		{left: int32(2), right: int(1), expect: true},
		{left: int32(1), right: int(1), expect: true},
		{left: int32(1), right: int8(1), expect: true},
		{left: int32(1), right: int16(1), expect: true},
		{left: int32(1), right: int32(1), expect: true},
		{left: int32(1), right: int64(1), expect: true},
		{left: int32(1), right: uint(1), expect: true},
		{left: int32(1), right: uint8(1), expect: true},
		{left: int32(1), right: uint16(1), expect: true},
		{left: int32(1), right: uint32(1), expect: true},
		{left: int32(1), right: uint64(1), expect: true},
		{left: int32(1), right: float32(1), isError: true},
		{left: int32(1), right: float64(1), isError: true},
		{left: int32(1), right: "foo", isError: true},
		{left: int32(1), right: testStruct{name: "foo"}, isError: true},
		{left: int32(1), right: []string{"foo"}, isError: true},

		// left is int64
		{left: int64(1), right: true, isError: true},
		{left: int64(2), right: int(1), expect: true},
		{left: int64(1), right: int(1), expect: true},
		{left: int64(1), right: int8(1), expect: true},
		{left: int64(1), right: int16(1), expect: true},
		{left: int64(1), right: int32(1), expect: true},
		{left: int64(1), right: int64(1), expect: true},
		{left: int64(1), right: uint(1), expect: true},
		{left: int64(1), right: uint8(1), expect: true},
		{left: int64(1), right: uint16(1), expect: true},
		{left: int64(1), right: uint32(1), expect: true},
		{left: int64(1), right: uint64(1), expect: true},
		{left: int64(1), right: float32(1), isError: true},
		{left: int64(1), right: float64(1), isError: true},
		{left: int64(1), right: "foo", isError: true},
		{left: int64(1), right: testStruct{name: "foo"}, isError: true},
		{left: int64(1), right: []string{"foo"}, isError: true},

		// left is uint
		{left: uint(1), right: true, isError: true},
		{left: uint(2), right: int(1), expect: true},
		{left: uint(1), right: int(1), expect: true},
		{left: uint(1), right: int8(1), expect: true},
		{left: uint(1), right: int16(1), expect: true},
		{left: uint(1), right: int32(1), expect: true},
		{left: uint(1), right: int64(1), expect: true},
		{left: uint(1), right: uint(1), expect: true},
		{left: uint(1), right: uint8(1), expect: true},
		{left: uint(1), right: uint16(1), expect: true},
		{left: uint(1), right: uint32(1), expect: true},
		{left: uint(1), right: uint64(1), expect: true},
		{left: uint(1), right: float32(1), isError: true},
		{left: uint(1), right: float64(1), isError: true},
		{left: uint(1), right: "foo", isError: true},
		{left: uint(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint(1), right: []string{"foo"}, isError: true},

		// left is uint8
		{left: uint8(1), right: true, isError: true},
		{left: uint8(2), right: int(1), expect: true},
		{left: uint8(1), right: int(1), expect: true},
		{left: uint8(1), right: int8(1), expect: true},
		{left: uint8(1), right: int16(1), expect: true},
		{left: uint8(1), right: int32(1), expect: true},
		{left: uint8(1), right: int64(1), expect: true},
		{left: uint8(1), right: uint(1), expect: true},
		{left: uint8(1), right: uint8(1), expect: true},
		{left: uint8(1), right: uint16(1), expect: true},
		{left: uint8(1), right: uint32(1), expect: true},
		{left: uint8(1), right: uint64(1), expect: true},
		{left: uint8(1), right: float32(1), isError: true},
		{left: uint8(1), right: float64(1), isError: true},
		{left: uint8(1), right: "foo", isError: true},
		{left: uint8(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint8(1), right: []string{"foo"}, isError: true},

		// left is uint16
		{left: uint16(1), right: true, isError: true},
		{left: uint16(2), right: int(1), expect: true},
		{left: uint16(1), right: int(1), expect: true},
		{left: uint16(1), right: int8(1), expect: true},
		{left: uint16(1), right: int16(1), expect: true},
		{left: uint16(1), right: int32(1), expect: true},
		{left: uint16(1), right: int64(1), expect: true},
		{left: uint16(1), right: uint(1), expect: true},
		{left: uint16(1), right: uint8(1), expect: true},
		{left: uint16(1), right: uint16(1), expect: true},
		{left: uint16(1), right: uint32(1), expect: true},
		{left: uint16(1), right: uint64(1), expect: true},
		{left: uint16(1), right: float32(1), isError: true},
		{left: uint16(1), right: float64(1), isError: true},
		{left: uint16(1), right: "foo", isError: true},
		{left: uint16(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint16(1), right: []string{"foo"}, isError: true},

		// left is uint32
		{left: uint32(1), right: true, isError: true},
		{left: uint32(2), right: int(1), expect: true},
		{left: uint32(1), right: int(1), expect: true},
		{left: uint32(1), right: int8(1), expect: true},
		{left: uint32(1), right: int16(1), expect: true},
		{left: uint32(1), right: int32(1), expect: true},
		{left: uint32(1), right: int64(1), expect: true},
		{left: uint32(1), right: uint(1), expect: true},
		{left: uint32(1), right: uint8(1), expect: true},
		{left: uint32(1), right: uint16(1), expect: true},
		{left: uint32(1), right: uint32(1), expect: true},
		{left: uint32(1), right: uint64(1), expect: true},
		{left: uint32(1), right: float32(1), isError: true},
		{left: uint32(1), right: float64(1), isError: true},
		{left: uint32(1), right: "foo", isError: true},
		{left: uint32(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint32(1), right: []string{"foo"}, isError: true},

		// left is uint64
		{left: uint64(1), right: true, isError: true},
		{left: uint64(1), right: int(1), expect: true},
		{left: uint64(1), right: int8(1), expect: true},
		{left: uint64(1), right: int16(1), expect: true},
		{left: uint64(1), right: int32(1), expect: true},
		{left: uint64(1), right: int64(1), expect: true},
		{left: uint64(1), right: uint(1), expect: true},
		{left: uint64(1), right: uint8(1), expect: true},
		{left: uint64(1), right: uint16(1), expect: true},
		{left: uint64(1), right: uint32(1), expect: true},
		{left: uint64(1), right: uint64(1), expect: true},
		{left: uint64(1), right: float32(1), isError: true},
		{left: uint64(1), right: float64(1), isError: true},
		{left: uint64(1), right: "foo", isError: true},
		{left: uint64(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint64(1), right: []string{"foo"}, isError: true},

		// left is float32
		{left: float32(1), right: true, isError: true},
		{left: float32(1), right: int(1), isError: true},
		{left: float32(1), right: int8(1), isError: true},
		{left: float32(1), right: int16(1), isError: true},
		{left: float32(1), right: int32(1), isError: true},
		{left: float32(1), right: int64(1), isError: true},
		{left: float32(1), right: uint(1), isError: true},
		{left: float32(1), right: uint8(1), isError: true},
		{left: float32(1), right: uint16(1), isError: true},
		{left: float32(1), right: uint32(1), isError: true},
		{left: float32(1), right: uint64(1), isError: true},
		{left: float32(1), right: float32(1), expect: true},
		{left: float32(1), right: float64(1), expect: true},
		{left: float64(1.1), right: float64(1), expect: true},
		{left: float32(1), right: "foo", isError: true},
		{left: float32(1), right: testStruct{name: "foo"}, isError: true},
		{left: float32(1), right: []string{"foo"}, isError: true},

		// left is float64
		{left: float64(1), right: true, isError: true},
		{left: float64(1), right: int(1), isError: true},
		{left: float64(1), right: int8(1), isError: true},
		{left: float64(1), right: int16(1), isError: true},
		{left: float64(1), right: int32(1), isError: true},
		{left: float64(1), right: int64(1), isError: true},
		{left: float64(1), right: uint(1), isError: true},
		{left: float64(1), right: uint8(1), isError: true},
		{left: float64(1), right: uint16(1), isError: true},
		{left: float64(1), right: uint32(1), isError: true},
		{left: float64(1), right: uint64(1), isError: true},
		{left: float64(1), right: float32(1), expect: true},
		{left: float64(1), right: float64(1), expect: true},
		{left: float64(1.1), right: float64(1), expect: true},
		{left: float64(1), right: "foo", isError: true},
		{left: float64(1), right: testStruct{name: "foo"}, isError: true},
		{left: float64(1), right: []string{"foo"}, isError: true},

		// left is string
		{left: "foo", right: true, isError: true},
		{left: "foo", right: int(1), isError: true},
		{left: "foo", right: int8(1), isError: true},
		{left: "foo", right: int16(1), isError: true},
		{left: "foo", right: int32(1), isError: true},
		{left: "foo", right: int64(1), isError: true},
		{left: "foo", right: uint(1), isError: true},
		{left: "foo", right: uint8(1), isError: true},
		{left: "foo", right: uint16(1), isError: true},
		{left: "foo", right: uint32(1), isError: true},
		{left: "foo", right: uint64(1), isError: true},
		{left: "foo", right: float32(1), isError: true},
		{left: "foo", right: float64(1), isError: true},
		{left: "foo", right: "foo", isError: true},
		{left: "foo", right: testStruct{name: "foo"}, isError: true},
		{left: "foo", right: []string{"foo"}, isError: true},

		// left is struct
		{left: testStruct{name: "foo"}, right: true, isError: true},
		{left: testStruct{name: "foo"}, right: int(1), isError: true},
		{left: testStruct{name: "foo"}, right: int8(1), isError: true},
		{left: testStruct{name: "foo"}, right: int16(1), isError: true},
		{left: testStruct{name: "foo"}, right: int32(1), isError: true},
		{left: testStruct{name: "foo"}, right: int64(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint8(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint16(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint32(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint64(1), isError: true},
		{left: testStruct{name: "foo"}, right: float32(1), isError: true},
		{left: testStruct{name: "foo"}, right: float64(1), isError: true},
		{left: testStruct{name: "foo"}, right: "foo", isError: true},
		{left: testStruct{name: "foo"}, right: testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "bar"}, isError: true},
		{left: testStruct{name: "foo"}, right: []string{"foo"}, isError: true},
	}

	for i, tt := range tests {
		ret, err := GreaterThanEqual(reflect.ValueOf(tt.left), reflect.ValueOf(tt.right))

		if tt.isError {
			if err == nil {
				t.Errorf("[%d] Expects error, got nil, %v", i, tt)
			}
			return
		}
		if err != nil {
			t.Errorf("[%d] Expects no error, got error %s", i, err)
			return
		}

		if diff := cmp.Diff(tt.expect, ret); diff != "" {
			t.Errorf("[%d] GreaterThan result mismatch, diff=%s", i, diff)
		}
	}
}

func TestCompareValuesLessThan(t *testing.T) {
	tests := []struct {
		left    any
		right   any
		expect  bool
		isError bool
	}{
		// left is bool
		{left: true, right: true, isError: true},
		{left: true, right: int(1), isError: true},
		{left: true, right: int8(1), isError: true},
		{left: true, right: int16(1), isError: true},
		{left: true, right: int32(1), isError: true},
		{left: true, right: int64(1), isError: true},
		{left: true, right: uint(1), isError: true},
		{left: true, right: uint8(1), isError: true},
		{left: true, right: uint16(1), isError: true},
		{left: true, right: uint32(1), isError: true},
		{left: true, right: uint64(1), isError: true},
		{left: true, right: float32(1), isError: true},
		{left: true, right: float64(1), isError: true},
		{left: true, right: "foo", isError: true},
		{left: true, right: testStruct{name: "foo"}, isError: true},
		{left: true, right: []string{"foo"}, isError: true},

		// left is int
		{left: int(1), right: true, isError: true},
		{left: int(2), right: int(1), expect: false},
		{left: int(0), right: int(1), expect: true},
		{left: int(1), right: int(1), expect: false},
		{left: int(1), right: int8(1), expect: false},
		{left: int(1), right: int16(1), expect: false},
		{left: int(1), right: int32(1), expect: false},
		{left: int(1), right: int64(1), expect: false},
		{left: int(1), right: uint(1), expect: false},
		{left: int(1), right: uint8(1), expect: false},
		{left: int(1), right: uint16(1), expect: false},
		{left: int(1), right: uint32(1), expect: false},
		{left: int(1), right: uint64(1), expect: false},
		{left: int(1), right: float32(1), isError: true},
		{left: int(1), right: float64(1), isError: true},
		{left: int(1), right: "foo", isError: true},
		{left: int(1), right: testStruct{name: "foo"}, isError: true},
		{left: int(1), right: []string{"foo"}, isError: true},

		// left is int8
		{left: int8(1), right: true, isError: true},
		{left: int8(2), right: int(1), expect: false},
		{left: int8(0), right: int(1), expect: true},
		{left: int8(1), right: int(1), expect: false},
		{left: int8(1), right: int8(1), expect: false},
		{left: int8(1), right: int16(1), expect: false},
		{left: int8(1), right: int32(1), expect: false},
		{left: int8(1), right: int64(1), expect: false},
		{left: int8(1), right: uint(1), expect: false},
		{left: int8(1), right: uint8(1), expect: false},
		{left: int8(1), right: uint16(1), expect: false},
		{left: int8(1), right: uint32(1), expect: false},
		{left: int8(1), right: uint64(1), expect: false},
		{left: int8(1), right: float32(1), isError: true},
		{left: int8(1), right: float64(1), isError: true},
		{left: int8(1), right: "foo", isError: true},
		{left: int8(1), right: testStruct{name: "foo"}, isError: true},
		{left: int8(1), right: []string{"foo"}, isError: true},

		// left is int16
		{left: int16(1), right: true, isError: true},
		{left: int16(2), right: int(1), expect: false},
		{left: int16(0), right: int(1), expect: true},
		{left: int16(1), right: int(1), expect: false},
		{left: int16(1), right: int8(1), expect: false},
		{left: int16(1), right: int16(1), expect: false},
		{left: int16(1), right: int32(1), expect: false},
		{left: int16(1), right: int64(1), expect: false},
		{left: int16(1), right: uint(1), expect: false},
		{left: int16(1), right: uint8(1), expect: false},
		{left: int16(1), right: uint16(1), expect: false},
		{left: int16(1), right: uint32(1), expect: false},
		{left: int16(1), right: uint64(1), expect: false},
		{left: int16(1), right: float32(1), isError: true},
		{left: int16(1), right: float64(1), isError: true},
		{left: int16(1), right: "foo", isError: true},
		{left: int16(1), right: testStruct{name: "foo"}, isError: true},
		{left: int16(1), right: []string{"foo"}, isError: true},

		// left is int32
		{left: int32(1), right: true, isError: true},
		{left: int32(2), right: int(1), expect: false},
		{left: int32(0), right: int(1), expect: true},
		{left: int32(1), right: int(1), expect: false},
		{left: int32(1), right: int8(1), expect: false},
		{left: int32(1), right: int16(1), expect: false},
		{left: int32(1), right: int32(1), expect: false},
		{left: int32(1), right: int64(1), expect: false},
		{left: int32(1), right: uint(1), expect: false},
		{left: int32(1), right: uint8(1), expect: false},
		{left: int32(1), right: uint16(1), expect: false},
		{left: int32(1), right: uint32(1), expect: false},
		{left: int32(1), right: uint64(1), expect: false},
		{left: int32(1), right: float32(1), isError: true},
		{left: int32(1), right: float64(1), isError: true},
		{left: int32(1), right: "foo", isError: true},
		{left: int32(1), right: testStruct{name: "foo"}, isError: true},
		{left: int32(1), right: []string{"foo"}, isError: true},

		// left is int64
		{left: int64(1), right: true, isError: true},
		{left: int64(2), right: int(1), expect: false},
		{left: int64(0), right: int(1), expect: true},
		{left: int64(1), right: int(1), expect: false},
		{left: int64(1), right: int8(1), expect: false},
		{left: int64(1), right: int16(1), expect: false},
		{left: int64(1), right: int32(1), expect: false},
		{left: int64(1), right: int64(1), expect: false},
		{left: int64(1), right: uint(1), expect: false},
		{left: int64(1), right: uint8(1), expect: false},
		{left: int64(1), right: uint16(1), expect: false},
		{left: int64(1), right: uint32(1), expect: false},
		{left: int64(1), right: uint64(1), expect: false},
		{left: int64(1), right: float32(1), isError: true},
		{left: int64(1), right: float64(1), isError: true},
		{left: int64(1), right: "foo", isError: true},
		{left: int64(1), right: testStruct{name: "foo"}, isError: true},
		{left: int64(1), right: []string{"foo"}, isError: true},

		// left is uint
		{left: uint(1), right: true, isError: true},
		{left: uint(2), right: int(1), expect: false},
		{left: uint(0), right: int(1), expect: true},
		{left: uint(1), right: int(1), expect: false},
		{left: uint(1), right: int8(1), expect: false},
		{left: uint(1), right: int16(1), expect: false},
		{left: uint(1), right: int32(1), expect: false},
		{left: uint(1), right: int64(1), expect: false},
		{left: uint(1), right: uint(1), expect: false},
		{left: uint(1), right: uint8(1), expect: false},
		{left: uint(1), right: uint16(1), expect: false},
		{left: uint(1), right: uint32(1), expect: false},
		{left: uint(1), right: uint64(1), expect: false},
		{left: uint(1), right: float32(1), isError: true},
		{left: uint(1), right: float64(1), isError: true},
		{left: uint(1), right: "foo", isError: true},
		{left: uint(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint(1), right: []string{"foo"}, isError: true},

		// left is uint8
		{left: uint8(1), right: true, isError: true},
		{left: uint8(2), right: int(1), expect: false},
		{left: uint8(0), right: int(1), expect: true},
		{left: uint8(1), right: int(1), expect: false},
		{left: uint8(1), right: int8(1), expect: false},
		{left: uint8(1), right: int16(1), expect: false},
		{left: uint8(1), right: int32(1), expect: false},
		{left: uint8(1), right: int64(1), expect: false},
		{left: uint8(1), right: uint(1), expect: false},
		{left: uint8(1), right: uint8(1), expect: false},
		{left: uint8(1), right: uint16(1), expect: false},
		{left: uint8(1), right: uint32(1), expect: false},
		{left: uint8(1), right: uint64(1), expect: false},
		{left: uint8(1), right: float32(1), isError: true},
		{left: uint8(1), right: float64(1), isError: true},
		{left: uint8(1), right: "foo", isError: true},
		{left: uint8(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint8(1), right: []string{"foo"}, isError: true},

		// left is uint16
		{left: uint16(1), right: true, isError: true},
		{left: uint16(2), right: int(1), expect: false},
		{left: uint16(0), right: int(1), expect: true},
		{left: uint16(1), right: int(1), expect: false},
		{left: uint16(1), right: int8(1), expect: false},
		{left: uint16(1), right: int16(1), expect: false},
		{left: uint16(1), right: int32(1), expect: false},
		{left: uint16(1), right: int64(1), expect: false},
		{left: uint16(1), right: uint(1), expect: false},
		{left: uint16(1), right: uint8(1), expect: false},
		{left: uint16(1), right: uint16(1), expect: false},
		{left: uint16(1), right: uint32(1), expect: false},
		{left: uint16(1), right: uint64(1), expect: false},
		{left: uint16(1), right: float32(1), isError: true},
		{left: uint16(1), right: float64(1), isError: true},
		{left: uint16(1), right: "foo", isError: true},
		{left: uint16(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint16(1), right: []string{"foo"}, isError: true},

		// left is uint32
		{left: uint32(1), right: true, isError: true},
		{left: uint32(2), right: int(1), expect: false},
		{left: uint32(0), right: int(1), expect: true},
		{left: uint32(1), right: int(1), expect: false},
		{left: uint32(1), right: int8(1), expect: false},
		{left: uint32(1), right: int16(1), expect: false},
		{left: uint32(1), right: int32(1), expect: false},
		{left: uint32(1), right: int64(1), expect: false},
		{left: uint32(1), right: uint(1), expect: false},
		{left: uint32(1), right: uint8(1), expect: false},
		{left: uint32(1), right: uint16(1), expect: false},
		{left: uint32(1), right: uint32(1), expect: false},
		{left: uint32(1), right: uint64(1), expect: false},
		{left: uint32(1), right: float32(1), isError: true},
		{left: uint32(1), right: float64(1), isError: true},
		{left: uint32(1), right: "foo", isError: true},
		{left: uint32(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint32(1), right: []string{"foo"}, isError: true},

		// left is uint64
		{left: uint64(1), right: true, isError: true},
		{left: uint64(2), right: int(1), expect: false},
		{left: uint64(0), right: int(1), expect: true},
		{left: uint64(1), right: int(1), expect: false},
		{left: uint64(1), right: int8(1), expect: false},
		{left: uint64(1), right: int16(1), expect: false},
		{left: uint64(1), right: int32(1), expect: false},
		{left: uint64(1), right: int64(1), expect: false},
		{left: uint64(1), right: uint(1), expect: false},
		{left: uint64(1), right: uint8(1), expect: false},
		{left: uint64(1), right: uint16(1), expect: false},
		{left: uint64(1), right: uint32(1), expect: false},
		{left: uint64(1), right: uint64(1), expect: false},
		{left: uint64(1), right: float32(1), isError: true},
		{left: uint64(1), right: float64(1), isError: true},
		{left: uint64(1), right: "foo", isError: true},
		{left: uint64(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint64(1), right: []string{"foo"}, isError: true},

		// left is float32
		{left: float32(1), right: true, isError: true},
		{left: float32(1), right: int(1), isError: true},
		{left: float32(1), right: int8(1), isError: true},
		{left: float32(1), right: int16(1), isError: true},
		{left: float32(1), right: int32(1), isError: true},
		{left: float32(1), right: int64(1), isError: true},
		{left: float32(1), right: uint(1), isError: true},
		{left: float32(1), right: uint8(1), isError: true},
		{left: float32(1), right: uint16(1), isError: true},
		{left: float32(1), right: uint32(1), isError: true},
		{left: float32(1), right: uint64(1), isError: true},
		{left: float32(1), right: float32(1), expect: false},
		{left: float32(2), right: float32(1), expect: false},
		{left: float32(0), right: float32(1), expect: true},
		{left: float32(1), right: float64(1), expect: false},
		{left: float64(1.1), right: float64(1), expect: false},
		{left: float32(1), right: "foo", isError: true},
		{left: float32(1), right: testStruct{name: "foo"}, isError: true},
		{left: float32(1), right: []string{"foo"}, isError: true},

		// left is float64
		{left: float64(1), right: true, isError: true},
		{left: float64(1), right: int(1), isError: true},
		{left: float64(1), right: int8(1), isError: true},
		{left: float64(1), right: int16(1), isError: true},
		{left: float64(1), right: int32(1), isError: true},
		{left: float64(1), right: int64(1), isError: true},
		{left: float64(1), right: uint(1), isError: true},
		{left: float64(1), right: uint8(1), isError: true},
		{left: float64(1), right: uint16(1), isError: true},
		{left: float64(1), right: uint32(1), isError: true},
		{left: float64(1), right: uint64(1), isError: true},
		{left: float64(1), right: float32(1), expect: false},
		{left: float64(1), right: float64(1), expect: false},
		{left: float64(1.1), right: float64(1), expect: false},
		{left: float64(2), right: float64(1), expect: false},
		{left: float64(0), right: float64(1), expect: true},
		{left: float64(1), right: "foo", isError: true},
		{left: float64(1), right: testStruct{name: "foo"}, isError: true},
		{left: float64(1), right: []string{"foo"}, isError: true},

		// left is string
		{left: "foo", right: true, isError: true},
		{left: "foo", right: int(1), isError: true},
		{left: "foo", right: int8(1), isError: true},
		{left: "foo", right: int16(1), isError: true},
		{left: "foo", right: int32(1), isError: true},
		{left: "foo", right: int64(1), isError: true},
		{left: "foo", right: uint(1), isError: true},
		{left: "foo", right: uint8(1), isError: true},
		{left: "foo", right: uint16(1), isError: true},
		{left: "foo", right: uint32(1), isError: true},
		{left: "foo", right: uint64(1), isError: true},
		{left: "foo", right: float32(1), isError: true},
		{left: "foo", right: float64(1), isError: true},
		{left: "foo", right: "foo", isError: true},
		{left: "foo", right: testStruct{name: "foo"}, isError: true},
		{left: "foo", right: []string{"foo"}, isError: true},

		// left is struct
		{left: testStruct{name: "foo"}, right: true, isError: true},
		{left: testStruct{name: "foo"}, right: int(1), isError: true},
		{left: testStruct{name: "foo"}, right: int8(1), isError: true},
		{left: testStruct{name: "foo"}, right: int16(1), isError: true},
		{left: testStruct{name: "foo"}, right: int32(1), isError: true},
		{left: testStruct{name: "foo"}, right: int64(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint8(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint16(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint32(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint64(1), isError: true},
		{left: testStruct{name: "foo"}, right: float32(1), isError: true},
		{left: testStruct{name: "foo"}, right: float64(1), isError: true},
		{left: testStruct{name: "foo"}, right: "foo", isError: true},
		{left: testStruct{name: "foo"}, right: testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "bar"}, isError: true},
		{left: testStruct{name: "foo"}, right: []string{"foo"}, isError: true},
	}

	for i, tt := range tests {
		ret, err := LessThan(reflect.ValueOf(tt.left), reflect.ValueOf(tt.right))

		if tt.isError {
			if err == nil {
				t.Errorf("[%d] Expects error, got nil, %v", i, tt)
			}
			return
		}
		if err != nil {
			t.Errorf("[%d] Expects no error, got error %s", i, err)
			return
		}

		if diff := cmp.Diff(tt.expect, ret); diff != "" {
			t.Errorf("[%d] GreaterThan result mismatch, diff=%s", i, diff)
		}
	}
}

func TestCompareValuesLessThanEqual(t *testing.T) {
	tests := []struct {
		left    any
		right   any
		expect  bool
		isError bool
	}{
		// left is bool
		{left: true, right: true, isError: true},
		{left: true, right: int(1), isError: true},
		{left: true, right: int8(1), isError: true},
		{left: true, right: int16(1), isError: true},
		{left: true, right: int32(1), isError: true},
		{left: true, right: int64(1), isError: true},
		{left: true, right: uint(1), isError: true},
		{left: true, right: uint8(1), isError: true},
		{left: true, right: uint16(1), isError: true},
		{left: true, right: uint32(1), isError: true},
		{left: true, right: uint64(1), isError: true},
		{left: true, right: float32(1), isError: true},
		{left: true, right: float64(1), isError: true},
		{left: true, right: "foo", isError: true},
		{left: true, right: testStruct{name: "foo"}, isError: true},
		{left: true, right: []string{"foo"}, isError: true},

		// left is int
		{left: int(1), right: true, isError: true},
		{left: int(2), right: int(1), expect: false},
		{left: int(0), right: int(1), expect: true},
		{left: int(1), right: int(1), expect: true},
		{left: int(1), right: int8(1), expect: true},
		{left: int(1), right: int16(1), expect: true},
		{left: int(1), right: int32(1), expect: true},
		{left: int(1), right: int64(1), expect: true},
		{left: int(1), right: uint(1), expect: true},
		{left: int(1), right: uint8(1), expect: true},
		{left: int(1), right: uint16(1), expect: true},
		{left: int(1), right: uint32(1), expect: true},
		{left: int(1), right: uint64(1), expect: true},
		{left: int(1), right: float32(1), isError: true},
		{left: int(1), right: float64(1), isError: true},
		{left: int(1), right: "foo", isError: true},
		{left: int(1), right: testStruct{name: "foo"}, isError: true},
		{left: int(1), right: []string{"foo"}, isError: true},

		// left is int8
		{left: int8(1), right: true, isError: true},
		{left: int8(2), right: int(1), expect: false},
		{left: int8(0), right: int(1), expect: true},
		{left: int8(1), right: int(1), expect: true},
		{left: int8(1), right: int8(1), expect: true},
		{left: int8(1), right: int16(1), expect: true},
		{left: int8(1), right: int32(1), expect: true},
		{left: int8(1), right: int64(1), expect: true},
		{left: int8(1), right: uint(1), expect: true},
		{left: int8(1), right: uint8(1), expect: true},
		{left: int8(1), right: uint16(1), expect: true},
		{left: int8(1), right: uint32(1), expect: true},
		{left: int8(1), right: uint64(1), expect: true},
		{left: int8(1), right: float32(1), isError: true},
		{left: int8(1), right: float64(1), isError: true},
		{left: int8(1), right: "foo", isError: true},
		{left: int8(1), right: testStruct{name: "foo"}, isError: true},
		{left: int8(1), right: []string{"foo"}, isError: true},

		// left is int16
		{left: int16(1), right: true, isError: true},
		{left: int16(2), right: int(1), expect: false},
		{left: int16(0), right: int(1), expect: true},
		{left: int16(1), right: int(1), expect: true},
		{left: int16(1), right: int8(1), expect: true},
		{left: int16(1), right: int16(1), expect: true},
		{left: int16(1), right: int32(1), expect: true},
		{left: int16(1), right: int64(1), expect: true},
		{left: int16(1), right: uint(1), expect: true},
		{left: int16(1), right: uint8(1), expect: true},
		{left: int16(1), right: uint16(1), expect: true},
		{left: int16(1), right: uint32(1), expect: true},
		{left: int16(1), right: uint64(1), expect: true},
		{left: int16(1), right: float32(1), isError: true},
		{left: int16(1), right: float64(1), isError: true},
		{left: int16(1), right: "foo", isError: true},
		{left: int16(1), right: testStruct{name: "foo"}, isError: true},
		{left: int16(1), right: []string{"foo"}, isError: true},

		// left is int32
		{left: int32(1), right: true, isError: true},
		{left: int32(2), right: int(1), expect: false},
		{left: int32(0), right: int(1), expect: true},
		{left: int32(1), right: int(1), expect: true},
		{left: int32(1), right: int8(1), expect: true},
		{left: int32(1), right: int16(1), expect: true},
		{left: int32(1), right: int32(1), expect: true},
		{left: int32(1), right: int64(1), expect: true},
		{left: int32(1), right: uint(1), expect: true},
		{left: int32(1), right: uint8(1), expect: true},
		{left: int32(1), right: uint16(1), expect: true},
		{left: int32(1), right: uint32(1), expect: true},
		{left: int32(1), right: uint64(1), expect: true},
		{left: int32(1), right: float32(1), isError: true},
		{left: int32(1), right: float64(1), isError: true},
		{left: int32(1), right: "foo", isError: true},
		{left: int32(1), right: testStruct{name: "foo"}, isError: true},
		{left: int32(1), right: []string{"foo"}, isError: true},

		// left is int64
		{left: int64(1), right: true, isError: true},
		{left: int64(2), right: int(1), expect: false},
		{left: int64(0), right: int(1), expect: true},
		{left: int64(1), right: int(1), expect: true},
		{left: int64(1), right: int8(1), expect: true},
		{left: int64(1), right: int16(1), expect: true},
		{left: int64(1), right: int32(1), expect: true},
		{left: int64(1), right: int64(1), expect: true},
		{left: int64(1), right: uint(1), expect: true},
		{left: int64(1), right: uint8(1), expect: true},
		{left: int64(1), right: uint16(1), expect: true},
		{left: int64(1), right: uint32(1), expect: true},
		{left: int64(1), right: uint64(1), expect: true},
		{left: int64(1), right: float32(1), isError: true},
		{left: int64(1), right: float64(1), isError: true},
		{left: int64(1), right: "foo", isError: true},
		{left: int64(1), right: testStruct{name: "foo"}, isError: true},
		{left: int64(1), right: []string{"foo"}, isError: true},

		// left is uint
		{left: uint(1), right: true, isError: true},
		{left: uint(2), right: int(1), expect: false},
		{left: uint(0), right: int(1), expect: true},
		{left: uint(1), right: int(1), expect: true},
		{left: uint(1), right: int8(1), expect: true},
		{left: uint(1), right: int16(1), expect: true},
		{left: uint(1), right: int32(1), expect: true},
		{left: uint(1), right: int64(1), expect: true},
		{left: uint(1), right: uint(1), expect: true},
		{left: uint(1), right: uint8(1), expect: true},
		{left: uint(1), right: uint16(1), expect: true},
		{left: uint(1), right: uint32(1), expect: true},
		{left: uint(1), right: uint64(1), expect: true},
		{left: uint(1), right: float32(1), isError: true},
		{left: uint(1), right: float64(1), isError: true},
		{left: uint(1), right: "foo", isError: true},
		{left: uint(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint(1), right: []string{"foo"}, isError: true},

		// left is uint8
		{left: uint8(1), right: true, isError: true},
		{left: uint8(2), right: int(1), expect: false},
		{left: uint8(0), right: int(1), expect: true},
		{left: uint8(1), right: int(1), expect: true},
		{left: uint8(1), right: int8(1), expect: true},
		{left: uint8(1), right: int16(1), expect: true},
		{left: uint8(1), right: int32(1), expect: true},
		{left: uint8(1), right: int64(1), expect: true},
		{left: uint8(1), right: uint(1), expect: true},
		{left: uint8(1), right: uint8(1), expect: true},
		{left: uint8(1), right: uint16(1), expect: true},
		{left: uint8(1), right: uint32(1), expect: true},
		{left: uint8(1), right: uint64(1), expect: true},
		{left: uint8(1), right: float32(1), isError: true},
		{left: uint8(1), right: float64(1), isError: true},
		{left: uint8(1), right: "foo", isError: true},
		{left: uint8(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint8(1), right: []string{"foo"}, isError: true},

		// left is uint16
		{left: uint16(1), right: true, isError: true},
		{left: uint16(2), right: int(1), expect: false},
		{left: uint16(0), right: int(1), expect: true},
		{left: uint16(1), right: int(1), expect: true},
		{left: uint16(1), right: int8(1), expect: true},
		{left: uint16(1), right: int16(1), expect: true},
		{left: uint16(1), right: int32(1), expect: true},
		{left: uint16(1), right: int64(1), expect: true},
		{left: uint16(1), right: uint(1), expect: true},
		{left: uint16(1), right: uint8(1), expect: true},
		{left: uint16(1), right: uint16(1), expect: true},
		{left: uint16(1), right: uint32(1), expect: true},
		{left: uint16(1), right: uint64(1), expect: true},
		{left: uint16(1), right: float32(1), isError: true},
		{left: uint16(1), right: float64(1), isError: true},
		{left: uint16(1), right: "foo", isError: true},
		{left: uint16(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint16(1), right: []string{"foo"}, isError: true},

		// left is uint32
		{left: uint32(1), right: true, isError: true},
		{left: uint32(2), right: int(1), expect: false},
		{left: uint32(0), right: int(1), expect: true},
		{left: uint32(1), right: int(1), expect: true},
		{left: uint32(1), right: int8(1), expect: true},
		{left: uint32(1), right: int16(1), expect: true},
		{left: uint32(1), right: int32(1), expect: true},
		{left: uint32(1), right: int64(1), expect: true},
		{left: uint32(1), right: uint(1), expect: true},
		{left: uint32(1), right: uint8(1), expect: true},
		{left: uint32(1), right: uint16(1), expect: true},
		{left: uint32(1), right: uint32(1), expect: true},
		{left: uint32(1), right: uint64(1), expect: true},
		{left: uint32(1), right: float32(1), isError: true},
		{left: uint32(1), right: float64(1), isError: true},
		{left: uint32(1), right: "foo", isError: true},
		{left: uint32(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint32(1), right: []string{"foo"}, isError: true},

		// left is uint64
		{left: uint64(1), right: true, isError: true},
		{left: uint64(2), right: int(1), expect: false},
		{left: uint64(0), right: int(1), expect: true},
		{left: uint64(1), right: int(1), expect: true},
		{left: uint64(1), right: int8(1), expect: true},
		{left: uint64(1), right: int16(1), expect: true},
		{left: uint64(1), right: int32(1), expect: true},
		{left: uint64(1), right: int64(1), expect: true},
		{left: uint64(1), right: uint(1), expect: true},
		{left: uint64(1), right: uint8(1), expect: true},
		{left: uint64(1), right: uint16(1), expect: true},
		{left: uint64(1), right: uint32(1), expect: true},
		{left: uint64(1), right: uint64(1), expect: true},
		{left: uint64(1), right: float32(1), isError: true},
		{left: uint64(1), right: float64(1), isError: true},
		{left: uint64(1), right: "foo", isError: true},
		{left: uint64(1), right: testStruct{name: "foo"}, isError: true},
		{left: uint64(1), right: []string{"foo"}, isError: true},

		// left is float32
		{left: float32(1), right: true, isError: true},
		{left: float32(1), right: int(1), isError: true},
		{left: float32(1), right: int8(1), isError: true},
		{left: float32(1), right: int16(1), isError: true},
		{left: float32(1), right: int32(1), isError: true},
		{left: float32(1), right: int64(1), isError: true},
		{left: float32(1), right: uint(1), isError: true},
		{left: float32(1), right: uint8(1), isError: true},
		{left: float32(1), right: uint16(1), isError: true},
		{left: float32(1), right: uint32(1), isError: true},
		{left: float32(1), right: uint64(1), isError: true},
		{left: float32(1), right: float32(1), expect: true},
		{left: float32(2), right: float32(1), expect: false},
		{left: float32(0), right: float32(1), expect: true},
		{left: float32(1), right: float64(1), expect: true},
		{left: float64(1.1), right: float64(1), expect: false},
		{left: float32(1), right: "foo", isError: true},
		{left: float32(1), right: testStruct{name: "foo"}, isError: true},
		{left: float32(1), right: []string{"foo"}, isError: true},

		// left is float64
		{left: float64(1), right: true, isError: true},
		{left: float64(1), right: int(1), isError: true},
		{left: float64(1), right: int8(1), isError: true},
		{left: float64(1), right: int16(1), isError: true},
		{left: float64(1), right: int32(1), isError: true},
		{left: float64(1), right: int64(1), isError: true},
		{left: float64(1), right: uint(1), isError: true},
		{left: float64(1), right: uint8(1), isError: true},
		{left: float64(1), right: uint16(1), isError: true},
		{left: float64(1), right: uint32(1), isError: true},
		{left: float64(1), right: uint64(1), isError: true},
		{left: float64(1), right: float32(1), expect: true},
		{left: float64(1), right: float64(1), expect: true},
		{left: float64(1.1), right: float64(1), expect: false},
		{left: float64(2), right: float64(1), expect: false},
		{left: float64(0), right: float64(1), expect: true},
		{left: float64(1), right: "foo", isError: true},
		{left: float64(1), right: testStruct{name: "foo"}, isError: true},
		{left: float64(1), right: []string{"foo"}, isError: true},

		// left is string
		{left: "foo", right: true, isError: true},
		{left: "foo", right: int(1), isError: true},
		{left: "foo", right: int8(1), isError: true},
		{left: "foo", right: int16(1), isError: true},
		{left: "foo", right: int32(1), isError: true},
		{left: "foo", right: int64(1), isError: true},
		{left: "foo", right: uint(1), isError: true},
		{left: "foo", right: uint8(1), isError: true},
		{left: "foo", right: uint16(1), isError: true},
		{left: "foo", right: uint32(1), isError: true},
		{left: "foo", right: uint64(1), isError: true},
		{left: "foo", right: float32(1), isError: true},
		{left: "foo", right: float64(1), isError: true},
		{left: "foo", right: "foo", isError: true},
		{left: "foo", right: testStruct{name: "foo"}, isError: true},
		{left: "foo", right: []string{"foo"}, isError: true},

		// left is struct
		{left: testStruct{name: "foo"}, right: true, isError: true},
		{left: testStruct{name: "foo"}, right: int(1), isError: true},
		{left: testStruct{name: "foo"}, right: int8(1), isError: true},
		{left: testStruct{name: "foo"}, right: int16(1), isError: true},
		{left: testStruct{name: "foo"}, right: int32(1), isError: true},
		{left: testStruct{name: "foo"}, right: int64(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint8(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint16(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint32(1), isError: true},
		{left: testStruct{name: "foo"}, right: uint64(1), isError: true},
		{left: testStruct{name: "foo"}, right: float32(1), isError: true},
		{left: testStruct{name: "foo"}, right: float64(1), isError: true},
		{left: testStruct{name: "foo"}, right: "foo", isError: true},
		{left: testStruct{name: "foo"}, right: testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "foo"}, isError: true},
		{left: testStruct{name: "foo"}, right: &testStruct{name: "bar"}, isError: true},
		{left: testStruct{name: "foo"}, right: []string{"foo"}, isError: true},
	}

	for i, tt := range tests {
		ret, err := LessThanEqual(reflect.ValueOf(tt.left), reflect.ValueOf(tt.right))

		if tt.isError {
			if err == nil {
				t.Errorf("[%d] Expects error, got nil, %v", i, tt)
			}
			return
		}
		if err != nil {
			t.Errorf("[%d] Expects no error, got error %s", i, err)
			return
		}

		if diff := cmp.Diff(tt.expect, ret); diff != "" {
			t.Errorf("[%d] GreaterThan result mismatch, diff=%s", i, diff)
		}
	}
}
