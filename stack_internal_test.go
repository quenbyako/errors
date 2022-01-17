package errors

func Callers(extraSkip uint) StackTrace { return callers(extraSkip + 1) }
