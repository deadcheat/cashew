package cashew

import (
	"fmt"
	"time"

	"github.com/deadcheat/cashew/values/consts"
)

// Ticket struct for ticket
type Ticket struct {
	// required fields
	ID             string
	Type           TicketType
	ClientHostName string
	CreatedAt      time.Time

	// required for service-ticket, proxy-ticket, proxy-granting-ticket and ticket-granting-ticket
	LastReferencedAt *time.Time

	// required for service ticket
	Service string

	// required for service ticket
	// that granted by proxy granting ticket and ticket granting ticket
	GrantedBy *Ticket

	// required for service or ticket granting ticket
	UserName string

	// required for proxy granting ticket
	IOU string

	// set true if this is service-ticket and primary
	Primary bool

	// required for ticket granting ticket, usually serialized json is set
	ExtraAttributes interface{}
}

// TicketType types of ticket
type TicketType int8

const (
	// TicketTypeService type of service ticket
	TicketTypeService TicketType = iota + 1
	// TicketTypeProxy type of proxy ticket
	TicketTypeProxy
	// TicketTypeTicketGranting type of ticket-granting ticket
	TicketTypeTicketGranting
	// TicketTypeProxyGranting type of proxy-granting ticket
	TicketTypeProxyGranting
	// TicketTypeProxyGrantingIOU type of proxy-granting ticket iou
	// TODO check is this really necessary
	TicketTypeProxyGrantingIOU
	// TicketTypeLogin type of login ticket
	TicketTypeLogin
)

// String implements Stringer
func (t TicketType) String() string {
	switch t {
	case TicketTypeLogin:
		return consts.TicketTypeStrLogin
	case TicketTypeService:
		return consts.TicketTypeStrService
	case TicketTypeProxy:
		return consts.TicketTypeStrProxy
	case TicketTypeTicketGranting:
		return consts.TicketTypeStrTicketGranting
	case TicketTypeProxyGranting:
		return consts.TicketTypeStrProxyGranting
	case TicketTypeProxyGrantingIOU:
		return consts.TicketTypeStrProxyGrantingIOU
	}
	// panic when
	panic("invalid ticket will be created, don't do that")
}

// Prefix return ticket-type prefix string
func (t TicketType) Prefix() string {
	switch t {
	case TicketTypeService:
		return consts.PrefixServiceTicket
	case TicketTypeLogin:
		return consts.PrefixLoginTicket
	case TicketTypeProxy:
		return consts.PrefixProxyTicket
	case TicketTypeTicketGranting:
		return consts.PrefixTicketGrantingCookie
	case TicketTypeProxyGranting:
		return consts.PrefixProxyGrantingTicket
	case TicketTypeProxyGrantingIOU:
		return consts.PrefixProxyGrantingTicketIOU
	}
	// panic when
	panic(fmt.Sprintf("invalid ticket will be created, don't do that. type: %#+v", t))
}

// ParseTicketType parse to ticket type from string
func ParseTicketType(s string) TicketType {
	switch s {
	case consts.TicketTypeStrLogin:
		return TicketTypeLogin
	case consts.TicketTypeStrService:
		return TicketTypeService
	case consts.TicketTypeStrProxy:
		return TicketTypeProxy
	case consts.TicketTypeStrTicketGranting:
		return TicketTypeTicketGranting
	case consts.TicketTypeStrProxyGranting:
		return TicketTypeProxyGranting
	case consts.TicketTypeStrProxyGrantingIOU:
		return TicketTypeProxyGrantingIOU
	}
	// panic when
	panic(fmt.Sprintf("invalid ticket was found, it can't be, you've gotta kidding me type: %#+v", s))
}
