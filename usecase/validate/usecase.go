package validate

import (
	"fmt"
	"net/url"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/validator/ticket"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	tr cashew.TicketRepository
	tv ticket.Validator
}

// New return new logout usecase
func New(r cashew.TicketRepository) cashew.ValidateUseCase {
	tv := ticket.New()
	return &UseCase{r, tv}
}

// ValidateLogin execute validation for login
func (u *UseCase) ValidateLogin(ticket *cashew.Ticket) (err error) {
	if ticket.Type != cashew.TicketTypeLogin {
		return errors.NewTicketTypeError(ticket.ID, ticket.Type)
	}
	return u.tv.Validate(ticket)
}

// ValidateService execute validation service and ticket
func (u *UseCase) ValidateService(ticket *cashew.Ticket, service *url.URL, renew bool) (err error) {
	return u.validateTicket(ticket, service, renew, false)
}

// ValidateProxy execute validation service and ticket
func (u *UseCase) ValidateProxy(ticket *cashew.Ticket, service *url.URL) (err error) {
	return u.validateTicket(ticket, service, false, true)
}

func (u *UseCase) validateTicket(t *cashew.Ticket, service fmt.Stringer, renew, allowProxy bool) error {
	if renew && !primary(t) {
		return errors.ErrServiceTicketIsNotPrimary
	}
	if !allowProxy && t.Type == cashew.TicketTypeProxy {
		return errors.NewTicketTypeError(t.ID, t.Type)
	}
	if t.Type != cashew.TicketTypeService && t.Type != cashew.TicketTypeProxy {
		return errors.NewTicketTypeError(t.ID, t.Type)
	}
	if t.Service != service.String() {
		return errors.ErrServiceURLNotMatched
	}
	return u.tv.Validate(t)
}

// ValidateProxyGranting execute validation service and ticket
func (u *UseCase) ValidateProxyGranting(ticket *cashew.Ticket) (err error) {
	if ticket.Type != cashew.TicketTypeProxyGranting {
		return errors.NewTicketTypeError(ticket.ID, ticket.Type)
	}

	return u.tv.Validate(ticket)
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
