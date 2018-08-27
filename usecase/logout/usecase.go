package logout

import "github.com/deadcheat/cashew"

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r cashew.TicketRepository
}

// Terminate delete all ticket related given ticket
func (u *UseCase) Terminate(t *cashew.Ticket) error {
	return u.r.DeleteRelatedTicket(t)
}
