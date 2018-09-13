package terminatelogout

import "github.com/deadcheat/cashew"

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r cashew.TicketRepository
}

// New return new logout usecase
func New(r cashew.TicketRepository) cashew.TerminateUseCase {
	return &UseCase{r}
}

// Terminate delete all ticket related given ticket
func (u *UseCase) Terminate(t *cashew.Ticket) error {
	return u.r.DeleteRelatedTicket(t)
}
