package login

import (
	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/timer"
	"github.com/deadcheat/cashew/values/errs"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r   cashew.TicketRepository
	idr cashew.IDRepository
}

// New return new usecase
func New(r cashew.TicketRepository, idr cashew.IDRepository) cashew.LoginUseCase {
	return &UseCase{r, idr}
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
func (u *UseCase) ServiceTicket(service string, tgt *cashew.Ticket) (t *cashew.Ticket, err error) {
	t.ID = u.idr.Issue(cashew.TicketTypeService)
	t.Type = cashew.TicketTypeService
	t.Service = service
	t.UserName = tgt.UserName
	t.GrantedBy = tgt
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// Login login auth
func (u *UseCase) Login() error {
	return nil
}
