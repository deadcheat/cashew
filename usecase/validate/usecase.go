package validate

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
func (u *UseCase) Validate(ticket string, service *url.URL, renew, allowProxy bool) (t *cashew.Ticket, err error) {
	t, err = u.r.Find(ticket)
	if err == sql.ErrNoRows {
		return nil, errs.ErrTicketNotFound
	}
	if err != nil {
		return
	}
	if err = validateTicket(t, service, renew, allowProxy); err != nil {
		return
	}
	if err = u.r.Consume(t); err != nil {
		return
	}
	return
}

func validateTicket(t *cashew.Ticket, service *url.URL, renew, allowProxy bool) error {
	if renew && !primary(t) {
		return errs.ErrServiceTicketIsNoPrimary
	}
	if !allowProxy && t.Type == cashew.TicketTypeProxy {
		return errs.ErrTicketIsProxyTicket
	}
	if t.Type != cashew.TicketTypeService && t.Type != cashew.TicketTypeProxy {
		return errs.ErrTicketTypeNotMatched
	}
	if t.Service != service.String() {
		return errs.ErrServiceURLNotMatched
	}
	return nil
}

func primary(t *cashew.Ticket) bool {
	if t == nil {
		return false
	}
	switch t.Type {
	case cashew.TicketTypeService:
		return t.Primary
	case cashew.TicketTypeProxy:
		g := t.GrantedBy
		for g != nil {
			if t.Type == cashew.TicketTypeService {
				return t.Primary
			}
			g = t.GrantedBy
		}
		return false
	default:
		return false
	}
}
