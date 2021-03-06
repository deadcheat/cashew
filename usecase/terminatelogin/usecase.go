package terminatelogin

import (
	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/validator/ticket"
)

// UseCase implemented cashew.TerminateUseCase
type UseCase struct {
	r  cashew.TicketRepository
	tv ticket.Validator
}

// New return new usecase
func New(r cashew.TicketRepository) cashew.TerminateUseCase {
	tv := ticket.New()
	return &UseCase{r, tv}
}

// Terminate delete login ticket
func (u *UseCase) Terminate(t *cashew.Ticket) error {
	if t.Type != cashew.TicketTypeLogin {
		return errors.NewTicketTypeError(t.ID, t.Type)
	}
	return u.r.Delete(t)
}
