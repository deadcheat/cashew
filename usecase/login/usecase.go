package login

import (
	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/validator/ticket"
	"github.com/deadcheat/cashew/values/errs"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r  cashew.TicketRepository
	tv ticket.Validator
}

// New return new usecase
func New(r cashew.TicketRepository) cashew.LoginUseCase {
	tv := ticket.New()
	return &UseCase{r, tv}
}

// Validate validate ticket identified id
func (u *UseCase) Validate(ticketType cashew.TicketType, t *cashew.Ticket) error {
	if t.Type != ticketType {
		return errs.ErrTicketTypeNotMatched
	}

	return u.tv.Validate(t)
}

// TerminateLogin delete login ticket
func (u *UseCase) TerminateLogin(t *cashew.Ticket) error {
	if t.Type != cashew.TicketTypeLogin {
		return errs.ErrInvalidTicketType
	}
	return u.r.Delete(t)
}
