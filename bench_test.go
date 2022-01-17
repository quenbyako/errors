package errors_test

import (
	"fmt"
	"testing"

	stderrors "errors"
	"github.com/quenbyako/errors"
)

func stdErrors(at, depth int) error {
	if at >= depth {
		return stderrors.New("no error")
	}
	return stdErrors(at+1, depth)
}

func ownErrors(at, depth int) error {
	if at >= depth {
		return errors.New("ye error")
	}
	return ownErrors(at+1, depth)
}

// GlobalE is an exported global to store the result of benchmark results,
// preventing the compiler from optimising the benchmark functions away.
var GlobalE interface{}

func BenchmarkErrors(b *testing.B) {
	type run struct {
		stack int
		std   bool
	}
	runs := []run{
		{10, false},
		{10, true},
		{100, false},
		{100, true},
		{1000, false},
		{1000, true},
	}
	for _, r := range runs {
		part := "quenbyako/errors"
		if r.std {
			part = "errors"
		}
		name := fmt.Sprintf("%s-stack-%d", part, r.stack)
		b.Run(name, func(b *testing.B) {
			var err error
			f := ownErrors
			if r.std {
				f = stdErrors
			}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = f(0, r.stack)
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}

func BenchmarkStackFormatting(b *testing.B) {
	type run struct {
		stack  int
		format string
	}
	runs := []run{
		{10, "%s"},
		{10, "%v"},
		{10, "%+v"},
		{30, "%s"},
		{30, "%v"},
		{30, "%+v"},
		{60, "%s"},
		{60, "%v"},
		{60, "%+v"},
	}

	var stackStr string
	for _, r := range runs {
		name := fmt.Sprintf("%s-stack-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := ownErrors(0, r.stack)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, err)
			}
			b.StopTimer()
		})
	}

	for _, r := range runs {
		name := fmt.Sprintf("%s-stacktrace-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := ownErrors(0, r.stack)
			stack := errors.Stack(err)
			if stack == nil {
				b.FailNow()
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, stack)
			}
			b.StopTimer()
		})
	}
	GlobalE = stackStr
}
