package cashew

import "net/http"

// Deliver delivery interface
type Deliver interface {
	Mount()
	GetLogin(w http.ResponseWriter, r *http.Request)
	PostLogin(w http.ResponseWriter, r *http.Request)
}

// UseCase define behaviors for Cas Server
type UseCase interface {
	ValidateTicket(id string) error
	ServiceTicket() *Ticket
	Login() error
}

// TicketRepository repository for ticket
type TicketRepository interface {
	Find(id string) (*Ticket, error)
	Create(*Ticket) error
}

// Authenticator interface for authenticate
type Authenticator interface {
}
