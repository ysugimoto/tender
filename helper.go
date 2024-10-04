package tender

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Must is a helper function that panics if render() method raised an error.
// Usage:
// rendered := tender.Must(tender.New(tmpl).With(...).Render())
func Must(v string, err error) string {
	if err != nil {
		panic(err)
	}
	return v
}

// Following function is just divided strings.TrimSpace function for trimming space left-only or right-only
var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func trimLeftSpace(s string) string {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return strings.TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	return s[start:]
}

func trimRightSpace(s string) string {
	stop := len(s)
	for ; stop > 0; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			// start has been already trimmed above, should trim end only
			return strings.TrimRightFunc(s[:stop], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	return s[:stop]
}

// strings.TrimRight implementation for *bytes.Buffer
func trimRightSpaceBuffer(buf *bytes.Buffer) {
	str := buf.String()
	buf.Reset()
	buf.WriteString(trimRightSpace(str))
}

// Check ident indicates environment variable reference.
// If variable name is constructed only "[A-Z_]*", returns true
func isEnvironmentVariable(ident string) bool {
	for i := range ident {
		if (ident[i] < 0x41 || ident[i] > 0x5A) && ident[i] != 0x5F {
			return false
		}
	}
	return true
}

var escapeMap = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'&':  "&amp;",
	'"':  "&quot;",
	'\'': "&apos;",
}

func escapeHTML(v string) string {
	buf := pool.Get().(*bytes.Buffer) // nolint:errcheck
	defer pool.Put(buf)

	buf.Reset()

	for i := range v {
		c := v[i]
		if c >= utf8.RuneSelf {
			buf.WriteByte(c)
			continue
		}
		if rep, ok := escapeMap[v[i]]; ok {
			buf.WriteString(rep)
			continue
		}
		buf.WriteByte(c)
	}

	return buf.String()
}
