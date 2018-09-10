package errors

import (
	"fmt"

	"github.com/deadcheat/cashew/values/errs"
)

// ErrorCode error code
type ErrorCode int

const (
	// ErrorCodeInvalidRequest INVALID_REQUEST
	ErrorCodeInvalidRequest ErrorCode = iota
	// ErrorCodeInvalidTicket INVALID_TICKET
	ErrorCodeInvalidTicket
	// ErrorCodeInvalidTicketSpec INVALID_TICKET_SPEC
	ErrorCodeInvalidTicketSpec
	// ErrorCodeInvalidService INVALID_SERVICE
	ErrorCodeInvalidService
	// ErrorCodeInternalError INTERNAL_ERROR
	ErrorCodeInternalError
	// ErrorCodeUnAuthorizedServiceProxy BAD_PGT
	ErrorCodeUnAuthorizedServiceProxy
	// ErrorCodeInvalidProxyCallback INVALID_PROXY_CALLBACK
	ErrorCodeInvalidProxyCallback
)

const (
	// XMLErrorCodeInvalidRequest  not all of the required request parameters were present
	XMLErrorCodeInvalidRequest = "INVALID_REQUEST"
	// XMLErrorCodeInvalidTicket the ticket provided was not valid,
	// or the ticket did not come from an initial login and “renew” was set on validation.
	// The body of the block of the XML response SHOULD describe the exact details.
	XMLErrorCodeInvalidTicket = "INVALID_TICKET"
	// XMLErrorCodeInvalidTicketSpec failure to meet the requirements of validation specification
	XMLErrorCodeInvalidTicketSpec = "INVALID_TICKET_SPEC"
	// XMLErrorCodeInvalidService the ticket provided was valid,
	// but the service specified did not match the service associated with the ticket.
	// CAS MUST invalidate the ticket and disallow future validation of that same ticket.
	XMLErrorCodeInvalidService = "INVALID_SERVICE"
	// XMLErrorCodeInternalError an internal error occurred during ticket validation
	XMLErrorCodeInternalError = "INTERNAL_ERROR"
	// XMLErrorCodeUnAuthorizedServiceProxy the pgt provided was invalid
	XMLErrorCodeUnAuthorizedServiceProxy = "UNAUTHORIZED_SERVICE_PROXY"
	// XMLErrorCodeInvalidProxyCallback  The proxy callback specified is invalid. The credentials specified for proxy authentication do not meet the security requirements
	XMLErrorCodeInvalidProxyCallback = "INVALID_PROXY_CALLBACK"
)

// Wrapper provides error code and message for xml output
type Wrapper interface {
	Code() string
	Message() string
	Is(ErrorCode) bool
	Err() error
}

// ErrorWrapper is hold error-code (what will be returned as xml)
// and actual error
type wrapper struct {
	errorCode ErrorCode
	err       error
}

// Code return ErrorCodeString
func (w *wrapper) Code() string {
	switch w.errorCode {
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
func (w *wrapper) Message() string {
	return w.err.Error()
}

// Err return message
func (w *wrapper) Err() error {
	return w.err
}

// Is return whether is errorcode same
func (w *wrapper) Is(ec ErrorCode) bool {
	return (w.errorCode == ec)
}

// New create with specified errorcode and error
func New(ec ErrorCode, err error) Wrapper {
	return &wrapper{
		ec,
		err,
	}
}

// NewInvalidRequest return invalid request wrapper
func NewInvalidRequest(err error) Wrapper {
	return New(ErrorCodeInvalidRequest, err)
}

// NewInvalidService return invalid service wrapper
func NewInvalidService(err error) Wrapper {
	return New(ErrorCodeInvalidService, err)
}

// NewInvalidTicket return invalid ticket wrapper
func NewInvalidTicket(ticket string, err error) Wrapper {
	return &invalidTicket{ticket, New(ErrorCodeInvalidTicket, err)}
}

type invalidTicket struct {
	ticket string
	Wrapper
}

func (i *invalidTicket) Message() string {
	switch i.Err() {
	case errs.ErrTicketIsProxyTicket:
		return fmt.Sprintf("Ticket %s is proxy ticket", i.ticket)
	case errs.ErrTicketTypeNotMatched:
		return fmt.Sprintf("Ticket %s is not service ticket", i.ticket)
	}

	return fmt.Sprintf("Ticket %s not recognized", i.ticket)
}

// NewInvalidTicketSpec return invalid ticket wrapper
func NewInvalidTicketSpec(err error) Wrapper {
	return New(ErrorCodeInvalidTicketSpec, err)
}

// NewUnAuthorizedServiceProxy return invalid ticket wrapper
func NewUnAuthorizedServiceProxy(err error) Wrapper {
	return New(ErrorCodeUnAuthorizedServiceProxy, err)
}

// NewInternalError return invalid ticket wrapper
func NewInternalError(err error) Wrapper {
	return New(ErrorCodeInternalError, err)
}

// NewInvalidProxyCallback return invalid proxy callback wrapper
func NewInvalidProxyCallback(err error) Wrapper {
	return New(ErrorCodeInvalidProxyCallback, err)
}
