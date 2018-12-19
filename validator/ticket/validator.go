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
	now := timer.Local.Now()
	// tickets will be expired in granting-default-expire-time(default: 8 hour) from the last time ticket has been referenced
	// this is a rule for proxy granting ticket and ticket granting ticket
	switch t.Type {
	case cashew.TicketTypeProxyGranting, cashew.TicketTypeTicketGranting:
		if t.LastReferencedAt != nil && now.After(t.LastReferencedAt.Add(time.Duration(foundation.App().GrantingDefaultExpire)*time.Second)) {
			return errors.ErrTicketHasBeenExpired
		}
	default:
		// hard time out will happen in default-hard-expire time(default 2hour) from ticket has been created
		if now.After(t.CreatedAt.Add(time.Duration(foundation.App().GrantingHardTimeout) * time.Second)) {
			return errors.ErrHardTimeoutTicket
		}
	}
	return nil
}
