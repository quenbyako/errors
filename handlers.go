package errors

import (
	"fmt"
	"reflect"
)

// ErrConverter is a function which converts one type of error into a different one. You can use it to convert
// errors in remapper functions to throw more correct and explicit errors
type ErrConverter = func(error) error

// ConstConverter is an alias function which "converts" any error to provided one in the `err` argument.
// Technically, it just throw provided error in argument without any modifications.
func ConstConverter(err error) ErrConverter {
	return func(error) error { return err }
}

// ErrRemapperFunc is a function which detects provided error type (or value) and returns positive response,
// if provided error matched this remapper. The implementation of function MUST return error value and true, if
// provided error matched this remapper, and return nil and false, if error is not matched as well.
//
// It's better to use already predefined functions `ValueRemapper`, `ValueRemapperFunc`, `TypeRemapperLegacy`,
// `TypeRemapperLegacyFunc`, `TypeRemapper` or `TypeRemapperFunc`, but if you want, you can create custom
// remapper.
type ErrRemapperFunc = func(error) (error, bool)

// Remap is a function, which can remap provided error to a different one.
func Remap(err error, remappers []ErrRemapperFunc) error {
	for _, remapper := range remappers {
		if e, ok := remapper(err); ok {
			return e
		}
	}
	return err
}

func ValueRemapper(comparedErr, convertTo error) ErrRemapperFunc {
	return ValueRemapperFunc(comparedErr, ConstConverter(convertTo))
}

func ValueRemapperFunc(comparedErr error, converter ErrConverter) ErrRemapperFunc {
	return func(err error) (error, bool) {
		if err == comparedErr {
			return converter(err), true
		}
		return nil, false
	}
}

func TypeRemapperLegacy(T, convertTo error) ErrRemapperFunc {
	return TypeRemapperLegacyF(T, ConstConverter(convertTo))
}

func TypeRemapperLegacyF(T error, converter ErrConverter) ErrRemapperFunc {
	t := reflect.TypeOf(T)

	return func(err error) (error, bool) {
		if reflect.TypeOf(err) == t {
			return converter(err), true
		}
		return nil, false
	}
}

// ErrConstantWrap must be used only as last wrapper, when no any other remaper didn't work. It always returns
// true as converter response and wrapping handling errors into errors.Wrap function. If error already
// contains stack trace, new stack trace will not provided here. Otherwise, there will be added stack trace to
// ErrConstantWrap call (not internal logic).
func ErrConstantWrap(message string, args ...interface{}) ErrRemapperFunc {
	return func(err error) (error, bool) {
		return wrap(err, fmt.Sprintf(message, args...), 1), true
	}
}
