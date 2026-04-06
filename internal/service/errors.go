package service

import "fmt"

const (
	CodeInvalidInput        = "invalid_input"
	CodeEnvironmentNotReady = "environment_not_ready"
	CodeConfigurationError  = "configuration_error"
	CodeRuntimeFailure      = "runtime_failure"
)

// AppError is a typed error with a stable machine-readable code and exit code.
type AppError struct {
	code     string
	message  string
	cause    error
	exitCode int
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *AppError) Code() string {
	if e == nil {
		return ""
	}
	return e.code
}

func (e *AppError) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *AppError) ExitCode() int {
	if e == nil {
		return 0
	}
	return e.exitCode
}

func invalidInput(format string, args ...any) error {
	return &AppError{
		code:     CodeInvalidInput,
		message:  fmt.Sprintf(format, args...),
		exitCode: 2,
	}
}

// InvalidInput exposes a stable invalid-input error for command wiring.
func InvalidInput(message string) error {
	return &AppError{
		code:     CodeInvalidInput,
		message:  message,
		exitCode: 2,
	}
}

func environmentNotReady(format string, args ...any) error {
	return &AppError{
		code:     CodeEnvironmentNotReady,
		message:  fmt.Sprintf(format, args...),
		exitCode: 3,
	}
}

func configurationError(err error, format string, args ...any) error {
	return &AppError{
		code:     CodeConfigurationError,
		message:  fmt.Sprintf(format, args...),
		cause:    err,
		exitCode: 4,
	}
}

func runtimeFailure(err error, format string, args ...any) error {
	return &AppError{
		code:     CodeRuntimeFailure,
		message:  fmt.Sprintf(format, args...),
		cause:    err,
		exitCode: 5,
	}
}

// ExitCodeForError returns the stable CLI exit code for the given error.
func ExitCodeForError(err error) int {
	if err == nil {
		return 0
	}

	var appErr *AppError
	if ok := AsAppError(err, &appErr); ok {
		return appErr.ExitCode()
	}

	return 5
}

// ErrorCodeForError returns the stable machine code for the given error.
func ErrorCodeForError(err error) string {
	if err == nil {
		return ""
	}

	var appErr *AppError
	if ok := AsAppError(err, &appErr); ok {
		return appErr.Code()
	}

	return CodeRuntimeFailure
}

// ErrorMessageForError extracts a user-facing message from the given error.
func ErrorMessageForError(err error) string {
	if err == nil {
		return ""
	}

	var appErr *AppError
	if ok := AsAppError(err, &appErr); ok {
		return appErr.Message()
	}

	return err.Error()
}

// AsAppError wraps errors.As without leaking the errors package everywhere.
func AsAppError(err error, target **AppError) bool {
	if err == nil {
		return false
	}

	appErr, ok := err.(*AppError)
	if ok {
		*target = appErr
		return true
	}

	type unwrapper interface{ Unwrap() error }
	for current := err; current != nil; {
		if ae, ok := current.(*AppError); ok {
			*target = ae
			return true
		}
		u, ok := current.(unwrapper)
		if !ok {
			return false
		}
		current = u.Unwrap()
	}

	return false
}
