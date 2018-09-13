package errors

// AppError is an interface aggregates error and ErrorView
type AppError interface {
	ErrorView
	error
}

// ErrorView is an interface to be used in view
type ErrorView interface {
	Code() string
	Message() string
	Is(ErrorCode) bool
	Err() error
}

// core will implement two interfaces, ErrorView and error.
type core struct {
	code ErrorCode
	err  error
}

// New create with specified errorcode and error
func New(ec ErrorCode, err error) AppError {
	return &core{
		ec,
		err,
	}
}

// AsErrorView return ErrorView if interface given is implemented as ErrorView
func AsErrorView(i interface{}) (ev ErrorView, ok bool) {
	ev, ok = i.(ErrorView)
	return
}

// IsAppError inspects interface given is error or not
func IsAppError(i interface{}) bool {
	switch i.(type) {
	case AppError:
		return true
	default:
		return false
	}
}

// IsErrorView inspectsA interface given is implemented as ErrorView or not
func IsErrorView(i interface{}) bool {
	switch i.(type) {
	case ErrorView:
		return true
	default:
		return false
	}
}

// IsError inspects interface given is error or not
func IsError(i interface{}) bool {
	switch i.(type) {
	case error:
		return true
	default:
		return false
	}
}

// Error implement for interface 'error'
func (c *core) Error() string {
	return c.Message()
}

// Code return ErrorCodeString
func (c *core) Code() string {
	switch c.code {
	case ErrorCodeInvalidRequest:
		return XMLErrorCodeInvalidRequest
	case ErrorCodeInvalidTicket:
		return XMLErrorCodeInvalidTicket
	case ErrorCodeInvalidTicketSpec:
		return XMLErrorCodeInvalidTicketSpec
	case ErrorCodeInvalidProxyCallback:
		return XMLErrorCodeInvalidProxyCallback
	case ErrorCodeInvalidService:
		return XMLErrorCodeInvalidService
	case ErrorCodeUnAuthorizedServiceProxy:
		return XMLErrorCodeUnAuthorizedServiceProxy
	default:
		return XMLErrorCodeInternalError
	}
}

// Message return message
func (c *core) Message() string {

	return c.err.Error()
}

// Is return whether is errorcode same
func (c *core) Is(ec ErrorCode) bool {
	return (c.code == ec)
}

// Err return message
func (c *core) Err() error {
	return c.err
}
