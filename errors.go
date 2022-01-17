package errors

import (
	"fmt"
	"io"
)

// fundamental is an error that has a message and a stack, but no caller.
//
// reason is to create custom fundamental error instead using stdlib error is to speed up benchmarks
type fundamental struct {
	msg   string
	stack StackTrace
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(text string) error { return newFundamental(text, 1) }

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return newFundamental(fmt.Sprintf(format, args...), 1)
}

func newFundamental(text string, extraSkip uint) error {
	return &fundamental{
		msg:   text,
		stack: callers(1 + extraSkip),
	}
}

func (f *fundamental) Error() string          { return f.msg }
func (f *fundamental) stackTrace() StackTrace { return f.stack }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg+"\n")
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

type withStack struct {
	error
	stack StackTrace
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error { return wStack(err, 1) }

func wStack(err error, extraSkip uint) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(1 + extraSkip),
	}
}

func (w *withStack) Unwrap() error          { return w.error }
func (w *withStack) stackTrace() StackTrace { return w.stack }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.error)
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

type withMessage struct {
	cause error
	msg   string
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	return wMessage(err, message)
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	return wMessage(err, fmt.Sprintf(format, args...))
}

func wMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return wrap(err, message, 1)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return wrap(err, fmt.Sprintf(format, args...), 1)
}

func wrap(err error, message string, extraSkip uint) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	if Stack(err) != nil {
		return err
	}
	return &withStack{
		err,
		callers(1 + extraSkip),
	}
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s: %+v", w.msg, w.cause)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		io.WriteString(s, "\""+w.Error()+"\"")
	}
}

// Stack returns stack trace of error
func Stack(err error) StackTrace {
	if err == nil {
		return nil
	}
	cause, ok := err.(interface{ stackTrace() StackTrace })
	if ok {
		return cause.stackTrace()
	}
	return Stack(Unwrap(err))
}

// Cause returns the underlying cause of the error, if possible (looking for the deepest error).
//
// If the error does not implement Unwrap, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		maybeErr := cause.Unwrap()
		if maybeErr == nil {
			return err
		}
		err = maybeErr
	}
	return err
}
