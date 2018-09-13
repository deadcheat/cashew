package ticket

import (
	"time"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/timer"
)

// Validator check ticket is validate or not
type Validator interface {
	Validate(*cashew.Ticket) error
}

type v struct {
}

// New return new provider
func New() Validator {
	return &v{}
}

// Check ticket validateion
func (p *v) Validate(t *cashew.Ticket) error {
	if timer.Local.Now().After(t.CreatedAt.Add(time.Duration(foundation.App().GrantingHardTimeout))) {
		return errors.ErrHardTimeoutTicket
	}
	switch t.Type {
	case cashew.TicketTypeProxyGranting, cashew.TicketTypeTicketGranting:
		if t.LastReferencedAt != nil && t.LastReferencedAt.After(timer.Local.Now().Add(time.Duration(foundation.App().GrantingDefaultExpire)*time.Second)) {
			return errors.ErrTicketHasBeenExpired
		}
	}
	return nil
}
