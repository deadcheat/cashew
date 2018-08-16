package login

import (
	"net/http"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/timer"
	"github.com/deadcheat/cashew/values/errs"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r   cashew.TicketRepository
	idr cashew.IDRepository
	chr cashew.ClientHostNameRepository
}

// New return new usecase
func New(r cashew.TicketRepository, idr cashew.IDRepository, chr cashew.ClientHostNameRepository) cashew.LoginUseCase {
	return &UseCase{r, idr, chr}
}

// ValidateTicket validate ticket identified id
func (u *UseCase) ValidateTicket(t cashew.TicketType, id string) (*cashew.Ticket, error) {
	if id == "" {
		return nil, errs.ErrNoTicketID
	}
	ticket, err := u.r.Find(id)
	if err != nil {
		return nil, err
	}

	if ticket.Type != t {
		return nil, errs.ErrTicketTypeNotMatched
	}

	if ticket.ExpiresAt.Before(timer.Local.Now()) {
		return nil, errs.ErrTicketHasBeenExpired
	}
	return ticket, nil
}

// ServiceTicket create new ServiceTicket
func (u *UseCase) ServiceTicket(r *http.Request, service string, tgt *cashew.Ticket) (t *cashew.Ticket, err error) {
	t.Type = cashew.TicketTypeService
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	t.Service = service
	t.UserName = tgt.UserName
	t.GrantedBy = tgt
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// LoginTicket create new LoginTicket
func (u *UseCase) LoginTicket(r *http.Request) (t *cashew.Ticket, err error) {
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
