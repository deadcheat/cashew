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
	ValidateTicket(ticketType TicketType, id string) error
	ServiceTicket(service string) (*Ticket, error)
	Login() error
}

// TicketRepository repository for ticket
type TicketRepository interface {
	Issue(t TicketType) (*Ticket, error)
	Find(t TicketType, id string) (*Ticket, error)
	Create(*Ticket) error
}

// AuthenticateUseCase interface for authenticate
type AuthenticateUseCase interface {
	Authenticate(id, pass string) error
}

// AuthenticatedUserRepository repository interface for authenticate
type AuthenticatedUserRepository interface {
	Find(id, pass string) error
}
