package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const unknown = "unknown"

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// FuncInfo returns the full path to the File and Line number of the source code that contains the
// function and its name for this Frame's program counter.
func (f Frame) FuncInfo() (file string, line int, name string) {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return unknown, 0, unknown
	}
	file, line = fn.FileLine(f.pc())
	return file, line, fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	file, line, name := f.FuncInfo()
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			if file == unknown {
				io.WriteString(s, file)
				return
			}
			io.WriteString(s, name)
			io.WriteString(s, "\n\t")
			io.WriteString(s, file)
		default:
			io.WriteString(s, path.Base(file))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(line))
	case 'n':
		io.WriteString(s, funcname(name))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// MarshalText formats a stacktrace Frame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (f Frame) MarshalText() ([]byte, error) {
	file, line, name := f.FuncInfo()
	if name == unknown {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, file, line)), nil
}

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				f.Format(s, verb)
				io.WriteString(s, "\n")
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}

func callers(extraSkip uint) StackTrace {
	// skip calls in stacktrace to ensure that runtime returns only func calls outside this package
	const defaultSkip uint = 2
	// maximum depth of stacktrace to save only important calls and to save some memory
	const depth = 32

	var pcs [depth]uintptr
	n := runtime.Callers(int(defaultSkip+extraSkip), pcs[:])

	stack := make(StackTrace, n)
	for i := 0; i < n; i++ { // not ranging to avoid allocating
		stack[i] = Frame(pcs[i])
	}

	return stack
}

// utils

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
