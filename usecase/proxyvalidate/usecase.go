package proxyvalidate

import (
	"database/sql"
	"net/url"

	"github.com/deadcheat/cashew/values/errs"

	"github.com/deadcheat/cashew"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r cashew.TicketRepository
}

// New return new logout usecase
func New(r cashew.TicketRepository) cashew.ValidateUseCase {
	return &UseCase{r}
}

// Validate execute validation service and ticket
func (u *UseCase) Validate(ticket string, service *url.URL, renew bool) (t *cashew.Ticket, err error) {
	t, err = u.r.Find(ticket)
	if err == sql.ErrNoRows {
		return nil, errs.ErrTicketNotFound
	}
	if err != nil {
		return
	}
	if renew && !t.Primary {
		return nil, errs.ErrServiceTicketIsNoPrimary
	}
	if t.Type != cashew.TicketTypeService && t.Type != cashew.TicketTypeProxyGranting {
		return nil, errs.ErrTicketTypeNotMatched
	}
	if t.Service != service.String() {
		return nil, errs.ErrServiceURLNotMatched
	}
	if err = u.r.Consume(t); err != nil {
		return
	}
	return
}
