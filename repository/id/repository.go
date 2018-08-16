package id

import (
	"fmt"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/values/consts"
	"github.com/rs/xid"
)

// Repository id repository implementation
type Repository struct{}

// New return new IDRepository
func New() cashew.IDRepository {
	return new(Repository)
}

// Issue creates new id
func (r *Repository) Issue(t cashew.TicketType) string {
	return fmt.Sprintf("%s-%s", prefix(t), xid.New().String())
}

// make prefix string from ticket type
func prefix(t cashew.TicketType) string {
	switch t {
	case cashew.TicketTypeLogin:
		return consts.PrefixLoginTicket
	case cashew.TicketTypeService:
		return consts.PrefixServiceTicket
	case cashew.TicketTypeTicketGranting:
		return consts.PrefixTicketGrantingCookie
	case cashew.TicketTypeProxy:
		return consts.PrefixProxyTicket
	case cashew.TicketTypeProxyGranting:
		return consts.PrefixProxyGrantingTicket
	case cashew.TicketTypeProxyGrantingIOU:
		return consts.PrefixProxyGrantingTicketIOU
	}
	panic("invalid ticket type may be creating")
}
