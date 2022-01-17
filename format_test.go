package errors_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quenbyako/errors"
)

func wrappedNew(message string) error { // This function will be mid-stack inlined in go 1.12+
	return errors.New(message)
}

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.New("error"),
		"%s",
		"error",
	}, {
		errors.New("error"),
		"%v",
		"error",
	}, {
		errors.New("error"),
		"%+v",
		"error\n" +
			errors.PkgName + ".TestFormatNew\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:33\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}, {
		errors.New("error"),
		"%q",
		`"error"`,
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		errors.Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		errors.Errorf("%s", "error"),
		"%+v",
		"error\n" +
			errors.PkgName + ".TestFormatErrorf\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:69\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Wrap(errors.New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		errors.Wrap(errors.New("error"), "error2"),
		"%v",
		"error2: error",
	}, {
		errors.Wrap(errors.New("error"), "error2"),
		"%+v",
		"error2: error\n" +
			errors.PkgName + ".TestFormatWrap\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:101\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%v",
		"error: EOF",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%+v",
		"error: EOF\n" +
			errors.PkgName + ".TestFormatWrap\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:119\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}, {
		errors.Wrap(errors.Wrap(io.EOF, "error1"), "error2"),
		"%+v",
		"error2: error1: EOF\n" +
			errors.PkgName + ".TestFormatWrap\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:129\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}, {
		errors.Wrap(errors.New("error with space"), "context"),
		"%q",
		`"context: error with space"`,
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2: EOF",
	}, {
		errors.Wrapf(io.EOF, "error%d", 2),
		"%v",
		"error2: EOF",
	}, {
		errors.Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"error2: EOF\n" +
			errors.PkgName + ".TestFormatWrapf\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:165\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}, {
		errors.Wrapf(errors.New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		errors.Wrapf(errors.New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		errors.Wrapf(errors.New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			errors.PkgName + ".TestFormatWrapf\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:183\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrappedNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+v",
		"error\n" +
			errors.PkgName + ".wrappedNew\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:16\n" +
			errors.PkgName + ".TestFormatWrappedNew\n" +
			"\t.+/" + errors.PkgNameRaw + "/format_test.go:207\n" +
			"testing.tRunner\n" +
			"\t.+/src/testing/testing.go:1194\n" +
			"runtime.goexit\n" +
			"\t.+/src/runtime/asm_amd64.s:1371\n",
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

type wrapper struct {
	wrap func(err error) error
	want []string
}

func prettyBlocks(blocks []string) string {
	var out []string

	for _, b := range blocks {
		out = append(out, fmt.Sprintf("%v", b))
	}

	return "   " + strings.Join(out, "\n   ")
}

func requireMultilineRegexp(t *testing.T, want, got string) {
	t.Helper()
	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")

	if len(wantLines) != len(gotLines) {
		assert.FailNow(t, fmt.Sprintf("mismatched lines count:\n"+
			"expected: %v\n"+
			"actual:   %v\n", len(wantLines), len(gotLines)))
	}

	for i, wantLine := range wantLines {
		require.Regexp(t, wantLine, gotLines[i])
	}
}
