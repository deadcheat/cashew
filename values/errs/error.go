package errs

import "errors"

var (
	// ErrNoServiceDetected error that request had no service url
	ErrNoServiceDetected = errors.New("no service is detected")

	// ErrNoTicketID error that ticket id is empty
	ErrNoTicketID = errors.New("no ticket id is passed")

	// ErrTicketTypeNotMatched error that ticket what has different ticket type was found
	ErrTicketTypeNotMatched = errors.New("this ticket has not matched ticket type")

	// ErrTicketHasBeenExpired error that ticket has been expired already
	ErrTicketHasBeenExpired = errors.New("ticket has already been expired")

	// ErrInvalidCredentials error that inputed credentials are not validated successfully
	ErrInvalidCredentials = errors.New("your credential is not validated")

	// ErrTicketGrantedTicketIsNotFound error when ticket has no granter
	ErrTicketGrantedTicketIsNotFound = errors.New("granting ticket is not found")

	// ErrInvalidMethodCall error when invalid argument for method
	ErrInvalidMethodCall = errors.New("you may failed how to invoke this method")

	// ErrInvalidTicketType error when unexpected ticket passed as arguments
	ErrInvalidTicketType = errors.New("this method may not be scheduled to receive this type of ticket")

	// ErrHardTimeoutTicket error when ticket reached hard-timeout
	ErrHardTimeoutTicket = errors.New("this ticket may be a hard-timed-out one")

	// ErrMultipleUserFound defined error when multiple users matched identification
	ErrMultipleUserFound = errors.New("there are many users to match user/password")

	// ErrServiceURLNotMatched defined error that service url is not matched granted one
	ErrServiceURLNotMatched = errors.New("ticket grant other service not this")
)
