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
)
