package errors

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
