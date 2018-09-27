package cashew

import (
	"net/http"
	"net/url"
)

// Deliver delivery interface
type Deliver interface {
	Mount()
}

// Executer is an interface as cli processor
type Executer interface {
	Execute()
}

// TicketUseCase define behaviors about finding and generating ticket
type TicketUseCase interface {
	Find(id string) (*Ticket, error)
	NewLogin(r *http.Request) (*Ticket, error)
	NewProxyGranting(r *http.Request, callbackURL *url.URL, st *Ticket) (*Ticket, error)
	NewService(r *http.Request, service *url.URL, tgt *Ticket, primary bool) (*Ticket, error)
	NewGranting(r *http.Request, username string, extraAttributes interface{}) (*Ticket, error)
	NewProxy(r *http.Request, service string, grantedBy *Ticket) (*Ticket, error)
	Consume(*Ticket) error
}

// TerminateUseCase define behaviors for terminate ticket
type TerminateUseCase interface {
	Terminate(*Ticket) error
}

// LogoutUseCase define behaviors for logout by tgt
type LogoutUseCase interface {
	Terminate(*Ticket) error
}

// ValidateUseCase define behaviors for validation
type ValidateUseCase interface {
	ValidateLogin(ticket *Ticket) error
	ValidateGranting(ticket *Ticket) error
	ValidateService(ticket *Ticket, service *url.URL, renew bool) error
	ValidateProxy(ticket *Ticket, service *url.URL) error
	ValidateProxyGranting(ticket *Ticket) error
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
	Authenticate(id, pass string) (map[string]interface{}, error)
}
