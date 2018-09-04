package cashew

import (
	"net/http"
	"net/url"
)

// Deliver delivery interface
type Deliver interface {
	Mount()
}

// TicketUseCase define behaviors about finding and generating ticket
type TicketUseCase interface {
	FindTicket(id string) (*Ticket, error)
	LoginTicket(r *http.Request) (*Ticket, error)
	ProxyGrantingTicket(r *http.Request, callbackURL *url.URL, st *Ticket) (*Ticket, error)
	ServiceTicket(r *http.Request, service *url.URL, tgt *Ticket, primary bool) (*Ticket, error)
	TicketGrantingTicket(r *http.Request, username string, extraAttributes interface{}) (*Ticket, error)
}

// LoginUseCase define behaviors for Cas Server
type LoginUseCase interface {
	Validate(TicketType, *Ticket) error
	TerminateLogin(*Ticket) error
}

// LogoutUseCase define behaviors for logout by tgt
type LogoutUseCase interface {
	Terminate(*Ticket) error
}

// ValidateUseCase define behaviors for validation
type ValidateUseCase interface {
	Validate(ticket string, service *url.URL, renew bool) (*Ticket, error)
}

// TicketRepository repository for ticket
type TicketRepository interface {
	Find(id string) (*Ticket, error)
	Delete(*Ticket) error
	DeleteRelatedTicket(*Ticket) error
	Create(*Ticket) error
	Consume(*Ticket) error
}

// IDRepository is an interface to issue an ID
type IDRepository interface {
	Issue(t TicketType) string
}

// ClientHostNameRepository is an interface to find real-ip/hostname
type ClientHostNameRepository interface {
	Ensure(r *http.Request) string
}

// ProxyCallBackRepository is an interface to dial proxy-callback-url
type ProxyCallBackRepository interface {
	Dial(u *url.URL, pgt, iou string) error
}

// AuthenticateUseCase interface for authenticate
type AuthenticateUseCase interface {
	Authenticate(id, pass string) error
}
