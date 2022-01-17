//go:build go1.18

package errors

func TypeRemapper[T error](convertTo error) ErrRemapperFunc {
	return TypeRemapperFunc[T](ConstConverter(convertTo))
}

func TypeRemapperFunc[T error](converter ErrConverter) ErrRemapperFunc {
	return func(err error) (error, bool) {
		if _, ok := err.(T); ok {
			return converter(err), true
		}
		return nil, false
	}
}
