package cashew

import (
	"net/http"
)

// Deliver delivery interface
type Deliver interface {
	Mount()
	GetLogin(w http.ResponseWriter, r *http.Request)
	PostLogin(w http.ResponseWriter, r *http.Request)
}

// LoginUseCase define behaviors for Cas Server
type LoginUseCase interface {
	ValidateTicket(ticketType TicketType, id string) (*Ticket, error)
	ServiceTicket(r *http.Request, service string, tgt *Ticket) (*Ticket, error)
	LoginTicket(r *http.Request) (*Ticket, error)
	Login() error
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

// AuthenticatedUserRepository repository interface for authenticate
type AuthenticatedUserRepository interface {
	Find(id, pass string) error
}
