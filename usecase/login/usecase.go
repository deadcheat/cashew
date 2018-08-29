package login

import (
	"net/http"
	"net/url"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/validator/ticket"
	"github.com/deadcheat/cashew/values/errs"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r   cashew.TicketRepository
	idr cashew.IDRepository
	chr cashew.ClientHostNameRepository
	tv  ticket.Validator
}

// New return new usecase
func New(r cashew.TicketRepository, idr cashew.IDRepository, chr cashew.ClientHostNameRepository) cashew.LoginUseCase {
	tv := ticket.New()
	return &UseCase{r, idr, chr, tv}
}

// ValidateTicket validate ticket identified id
func (u *UseCase) ValidateTicket(ticketType cashew.TicketType, t *cashew.Ticket) error {
	if t.Type != ticketType {
		return errs.ErrTicketTypeNotMatched
	}

	return u.tv.Validate(t)
}

// FindTicket find ticket by id
func (u *UseCase) FindTicket(id string) (*cashew.Ticket, error) {
	if id == "" {
		return nil, errs.ErrNoTicketID
	}
	return u.r.Find(id)
}

// ServiceTicket create new ServiceTicket
func (u *UseCase) ServiceTicket(r *http.Request, service *url.URL, tgt *cashew.Ticket) (t *cashew.Ticket, err error) {
	if service == nil {
		return nil, errs.ErrNoServiceDetected
	}
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeService
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	t.Service = service.String()
	t.UserName = tgt.UserName
	t.GrantedBy = tgt
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// TicketGrantingTicket create new ServiceTicket
func (u *UseCase) TicketGrantingTicket(r *http.Request, username string, extraAttributes interface{}) (t *cashew.Ticket, err error) {
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeTicketGranting
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	t.UserName = username
	t.ExtraAttributes = extraAttributes
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// LoginTicket create new LoginTicket
func (u *UseCase) LoginTicket(r *http.Request) (t *cashew.Ticket, err error) {
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeLogin
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// Login login auth
func (u *UseCase) Login() error {
	return nil
}

// TerminateLoginTicket delete login ticket
func (u *UseCase) TerminateLoginTicket(t *cashew.Ticket) error {
	if t.Type != cashew.TicketTypeLogin {
		return errs.ErrInvalidTicketType
	}
	return u.r.Delete(t)
}
