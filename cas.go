package cashew

import (
	"net/http"
	"net/url"
)

// Deliver delivery interface
type Deliver interface {
	Mount()
}

// LoginUseCase define behaviors for Cas Server
type LoginUseCase interface {
	ValidateTicket(ticketType TicketType, id string) (*Ticket, error)
	ServiceTicket(r *http.Request, service *url.URL, tgt *Ticket) (*Ticket, error)
	TicketGrantingTicket(r *http.Request, username string, extraAttributes interface{}) (*Ticket, error)
	LoginTicket(r *http.Request) (*Ticket, error)
}

// TicketRepository repository for ticket
type TicketRepository interface {
	Find(id string) (*Ticket, error)
	Create(t *Ticket) error
}

// IDRepository is an interface to issue an ID
type IDRepository interface {
	Issue(t TicketType) string
}

// ClientHostNameRepository is an interface to find real-ip/hostname
type ClientHostNameRepository interface {
	Ensure(r *http.Request) string
}

// AuthenticateUseCase interface for authenticate
type AuthenticateUseCase interface {
	Authenticate(id, pass string) error
}
