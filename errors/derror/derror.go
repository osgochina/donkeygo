package derror

import "fmt"

// apiCode is the interface for Code feature.
type apiCode interface {
	Error() string // It should be an error.
	Code() int
}

// apiStack is the interface for Stack feature.
type apiStack interface {
	Error() string // It should be an error.
	Stack() string
}

// apiCause is the interface for Cause feature.
type apiCause interface {
	Error() string // It should be an error.
	Cause() error
}

// apiCurrent is the interface for Current feature.
type apiCurrent interface {
	Error() string // It should be an error.
	Current() error
}

// apiNext is the interface for Next feature.
type apiNext interface {
	Error() string // It should be an error.
	Next() error
}

func New(text string) error {
	return &Error{
		code:  -1,
		text:  text,
		stack: callers(),
	}
}

func Newf(format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  -1,
	}
}

func NewSkip(skip int, text string) error {
	return &Error{
		stack: callers(skip),
		text:  text,
		code:  -1,
	}
}

func NewSkipf(skip int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  -1,
	}
}

func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  text,
		code:  Code(err),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}

func WrapSkip(skip int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  text,
		code:  Code(err),
	}
}

func WrapSkipf(skip int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}

func NewCode(code int, text string) error {
	return &Error{
		stack: callers(),
		text:  text,
		code:  code,
	}
}

// NewCodef returns an error that has error code and formats as the given format and args.
func NewCodef(code int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// NewCodeSkip creates and returns an error which has error code and is formatted from given text.
// The parameter <skip> specifies the stack callers skipped amount.
func NewCodeSkip(code, skip int, text string) error {
	return &Error{
		stack: callers(skip),
		text:  text,
		code:  code,
	}
}

// NewCodeSkipf returns an error that has error code and formats as the given format and args.
// The parameter <skip> specifies the stack callers skipped amount.
func NewCodeSkipf(code, skip int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// WrapCode wraps error with code and text.
// It returns nil if given err is nil.
func WrapCode(code int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  text,
		code:  code,
	}
}

// WrapCodef wraps error with code and format specifier.
// It returns nil if given <err> is nil.
func WrapCodef(code int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// WrapCodeSkip wraps error with code and text.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapCodeSkip(code, skip int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  text,
		code:  code,
	}
}

// WrapCodeSkipf wraps error with code and text that is formatted with given format and args.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapCodeSkipf(code, skip int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// Cause returns the error code of current error.
// It returns -1 if it has no error code or it does not implements interface Code.
func Code(err error) int {
	if err != nil {
		if e, ok := err.(apiCode); ok {
			return e.Code()
		}
	}
	return -1
}

// Cause returns the root cause error of <err>.
func Cause(err error) error {
	if err != nil {
		if e, ok := err.(apiCause); ok {
			return e.Cause()
		}
	}
	return err
}

// Stack returns the stack callers as string.
// It returns the error string directly if the <err> does not support stacks.
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(apiStack); ok {
		return e.Stack()
	}
	return err.Error()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func Current(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(apiCurrent); ok {
		return e.Current()
	}
	return err
}

// Next returns the next level error.
// It returns nil if current level error or the next level error is nil.
func Next(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(apiNext); ok {
		return e.Next()
	}
	return nil
}
