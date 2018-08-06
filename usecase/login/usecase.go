package login

import (
	"errors"
	"time"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/values/consts"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r cashew.TicketRepository
}

// New return new usecase
func New(r cashew.TicketRepository) cashew.LoginUseCase {
	return &UseCase{r}
}

// ValidateTicket validate ticket identified id
func (u *UseCase) ValidateTicket(t cashew.TicketType, id string) error {
	ticket, err := u.r.Find(t, id)
	if err != nil {
		return err
	}
	// TODO fix to use interfaced timer
	if ticket.ExpiredAt.Before(time.Now()) {
		// TODO define this error in values package
		return errors.New("ticket has already been expired")
	}
	return nil
}

// ServiceTicket create new ServiceTicket
func (u *UseCase) ServiceTicket(service string) (t *cashew.Ticket, err error) {
	if t, err = u.r.Issue(consts.TicketTypeService); err != nil {
		return
	}
	t.Service = service
	if err = u.r.Create(t); err != nil {
		return nil, err
	}
	return
}

// Login login auth
func (u *UseCase) Login() error {
	return nil
}
