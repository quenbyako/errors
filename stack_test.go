package errors_test

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/k0kubun/pp"
	"github.com/quenbyako/errors"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() errors.Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() errors.Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		errors.Frame
		format string
		want   string
	}{{
		initpc,
		"%s",
		"stack_test.go",
	}, {
		initpc,
		"%+s",
		errors.PkgName + ".init\n" +
			`.+/` + errors.PkgNameRaw + `/stack_test\.go`,
	}, {
		0,
		"%s",
		"unknown",
	}, {
		0,
		"%+s",
		"unknown",
	}, {
		initpc,
		"%d",
		"12",
	}, {
		0,
		"%d",
		"0",
	}, {
		initpc,
		"%n",
		"init",
	}, {
		func() errors.Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`\(\*X\).ptr`,
	}, {
		func() errors.Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		0,
		"%n",
		"",
	}, {
		initpc,
		"%v",
		"stack_test.go:12",
	}, {
		initpc,
		"%+v",
		errors.PkgName + ".init\n" +
			"\t.+/" + errors.PkgNameRaw + "/stack_test.go:12",
	}, {
		0,
		"%v",
		"unknown:0",
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.Frame))
		})

	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/pkg/errors.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		errors.StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]errors.Frame\(nil\)`,
	}, {
		errors.StackTrace{},
		"%s",
		`\[\]`,
	}, {
		errors.StackTrace{},
		"%v",
		`\[\]`,
	}, {
		errors.StackTrace{},
		"%+v",
		"",
	}, {
		errors.StackTrace{},
		"%#v",
		`\[\]errors.Frame{}`,
	}, {
		stackyCaller()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stackyCaller()[:2],
		"%v",
		`\[stack_test.go:189 stack_test.go:164\]`,
	}, {
		stackyCaller()[:2],
		"%+v",
		errors.PkgName + ".stackyCaller\n" +
			"\t.+/" + errors.PkgNameRaw + "/stack_test.go:189\n" +
			errors.PkgName + ".TestStackTraceFormat\n" +
			"\t.+/" + errors.PkgNameRaw + "/stack_test.go:168\n",
	}, {
		stackyCaller()[:2],
		"%#v",
		`\[\]errors.Frame{stack_test.go:189, stack_test.go:175}`,
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			requireMultilineRegexp(t, tt.want, fmt.Sprintf(tt.format, tt.StackTrace))
		})
	}
}

// a version of runtime.Caller that returns a single Frame, not a full stacktrace.
func caller() errors.Frame            { return errors.Callers(1)[0] }
func stackyCaller() errors.StackTrace { return errors.Callers(0) }

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
