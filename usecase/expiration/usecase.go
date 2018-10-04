package expiration

import (
	"github.com/deadcheat/cashew"
)

// UseCase is a struct implemnts cashew.ExpirationUseCase
type UseCase struct {
	er cashew.ExpirationRepository
	tr cashew.TicketRepository
}

// New return new implement of cashew.ExpirationUseCase
func New(er cashew.ExpirationRepository, tr cashew.TicketRepository) cashew.ExpirationUseCase {
	return &UseCase{er, tr}
}

// RevokeAll remove all time-out tickets
func (u *UseCase) RevokeAll() error {
	// remove all expired ticket
	tickets, err := u.er.FindAll()
	if err != nil {
		return err
	}
	for i := range tickets {
		ticket := tickets[i]
		if err = u.tr.DeleteRelatedTicket(ticket); err != nil {
			return
		}
	}
	return nil
}
