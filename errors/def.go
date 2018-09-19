package errors

import (
	"errors"
	"fmt"

	"github.com/deadcheat/cashew"
)

var (
	// --- ErrorCodeInvalidRequest

	// ErrAuthenticateFailed authentication failed
	ErrAuthenticateFailed = New(ErrorCodeInvalidRequest, errors.New("authentication failed"))

	// ErrRequiredParameterMissed request has been sent without rquired parameters
	ErrRequiredParameterMissed = New(ErrorCodeInvalidRequest, errors.New("required parameters are not satisfied"))

	// ErrNoServiceDetected error that request had no service url
	ErrNoServiceDetected = New(ErrorCodeInvalidRequest, errors.New("no service is detected"))

	// ErrNoTicketID error that ticket id is empty
	ErrNoTicketID = New(ErrorCodeInvalidRequest, errors.New("no ticket id is passed"))

	// ErrInvalidCredentials error that inputed credentials are not validated successfully
	ErrInvalidCredentials = New(ErrorCodeInvalidRequest, errors.New("your credential is not validated"))

	// ---  ErrorCodeInvalidTicket

	// ErrTicketNotFound error that request had no ticket id
	ErrTicketNotFound = New(ErrorCodeInvalidTicket, errors.New("ticket is not found"))

	// --- ErrorCodeInvalidTicketSpec

	// ErrTicketHasBeenExpired error that ticket has been expired already
	ErrTicketHasBeenExpired = New(ErrorCodeInvalidTicketSpec, errors.New("ticket has already been expired"))

	// ErrTicketGrantedTicketIsNotFound error when ticket has no granter
	ErrTicketGrantedTicketIsNotFound = New(ErrorCodeInvalidTicketSpec, errors.New("granting ticket is not found"))

	// ErrHardTimeoutTicket error when ticket reached hard-timeout
	ErrHardTimeoutTicket = New(ErrorCodeInvalidTicketSpec, errors.New("this ticket may be a hard-timed-out one"))

	// ErrServiceURLNotMatched defined error that service url is not matched granted one
	ErrServiceURLNotMatched = New(ErrorCodeInvalidTicketSpec, errors.New("ticket grant other service not this"))

	// ErrServiceTicketIsNotPrimary defined error that service ticket is not primary one
	ErrServiceTicketIsNotPrimary = New(ErrorCodeInvalidTicketSpec, errors.New("ticket is not authenticated primary"))

	// --- ErrorCodeUnAuthorizedServiceProxy

	// ErrProxyGrantingURLUnexpectedStatus defined error that pgt request is failed
	ErrProxyGrantingURLUnexpectedStatus = New(ErrorCodeUnAuthorizedServiceProxy, errors.New("requested to given pgtUrl  response is unexpected status"))

	// --- ErrorCodeInvalidProxyCallback

	// ErrProxyCallBackURLMissing defined error that paramter pgtUrl is missing
	ErrProxyCallBackURLMissing = New(ErrorCodeInvalidProxyCallback, errors.New("a URL for proxy callback is missing"))
)

// TicketTypeError needs to output ticket id and type
type TicketTypeError struct {
	id string
	t  cashew.TicketType
	AppError
}

// Message TicketTypeError needs to output id and ticket type as message
func (tte *TicketTypeError) Message() string {
	return fmt.Sprintf("ticket \"%s\" is %s", tte.id, tte.t.String())
}

// NewTicketTypeError create new ticket type error as AppError
func NewTicketTypeError(id string, t cashew.TicketType) AppError {
	return &TicketTypeError{
		id, t,
		New(ErrorCodeInvalidTicket, errors.New("ticket type is not valid")),
	}
}

// InternalError implements messages for internal server error
type InternalError struct {
	AppError
}

// Message internal server error hides message into own
func (ie *InternalError) Message() string {
	return "the server encountered an internal error"
}

// NewInternalError return new internal error implements
func NewInternalError(err error) AppError {
	return &InternalError{
		New(ErrorCodeInternalError, err),
	}
}
