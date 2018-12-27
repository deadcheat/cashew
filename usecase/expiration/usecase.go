package expiration

import (
	"log"

	"github.com/deadcheat/cashew"
)

// UseCase is a struct implemnts cashew.ExpirationUseCase
type UseCase struct {
	// ExpirationRepo for granting ticket
	egr cashew.ExpirationRepository
	// ExpirationRepo for login ticket
	elr cashew.ExpirationRepository
	// ExpirationRepo for service/proxy ticket
	esr cashew.ExpirationRepository
	tr  cashew.TicketRepository
}

// New return new implement of cashew.ExpirationUseCase
func New(egr cashew.ExpirationRepository, elr cashew.ExpirationRepository, esr cashew.ExpirationRepository, tr cashew.TicketRepository) cashew.ExpirationUseCase {
	return &UseCase{egr, elr, esr, tr}
}

// RevokeAll remove all time-out tickets
func (u *UseCase) RevokeAll() error {
	// remove all expired ticket
	grantingTickets, err := u.egr.FindAll()
	if err != nil {
		return err
	}
	// remove all expired ticket
	loginTickets, err := u.elr.FindAll()
	if err != nil {
		return err
	}
	// remove all expired ticket
	serviceTickets, err := u.esr.FindAll()
	if err != nil {
		return err
	}
	tickets := make([]*cashew.Ticket, 0)
	tickets = append(tickets, grantingTickets...)
	tickets = append(tickets, loginTickets...)
	tickets = append(tickets, serviceTickets...)
	for i := range tickets {
		ticket := tickets[i]
		if ticket == nil {
			continue
		}
		log.Printf("delete timed-out ticket %s \n", ticket.ID)
		if err = u.tr.DeleteRelatedTicket(ticket); err != nil {
			return err
		}
	}
	return nil
}
