package tender

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIsEnvironmentVariable(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{input: "foobarBaz", expect: false},
		{input: "Foobar-baz", expect: false},
		{input: "FOO-BAR", expect: false},
		{input: "FOO_BAR", expect: true},
	}

	for _, tt := range tests {
		if diff := cmp.Diff(tt.expect, isEnvironmentVariable(tt.input)); diff != "" {
			t.Errorf("isEnvironmentVariable result mismatch, diff=%s", diff)
		}
	}
}
