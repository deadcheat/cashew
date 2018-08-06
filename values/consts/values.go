package consts

import "github.com/deadcheat/cashew"

const (
	// RFC2822 format of time
	RFC2822 = "Mon Jan 02 15:04:05 -0700 2006"
)

const (
	// ParamKeyService key string for param
	ParamKeyService = "service"
	// ParamKeyRenew key string for param
	ParamKeyRenew = "renew"
	// ParamKeyGateway key string for param
	ParamKeyGateway = "gateway"
)

const (
	// CookieKeyTGT key string for cookie
	CookieKeyTGT = "tgt"
)

const (
	// TicketTypeService type of service ticket
	TicketTypeService cashew.TicketType = iota + 1
	// TicketTypeProxy type of proxy ticket
	TicketTypeProxy
	// TicketTypeProxyGranting type of proxy-granting ticket
	TicketTypeProxyGranting
	// TicketTypeProxyGrantingIOU type of proxy-granting ticket iou
	TicketTypeProxyGrantingIOU
	// TicketTypeLogin type of login ticket
	TicketTypeLogin
)

const (
	// PrefixServiceTicket prefix string for service-ticket
	PrefixServiceTicket = "ST"
	// PrefixProxyTicket prefix string for proxy-ticket
	PrefixProxyTicket = "PT"
	// PrefixProxyGrantingTicket prefix string for proxy-ticket
	PrefixProxyGrantingTicket = "PGT"
	// PrefixProxyGrantingTicketIOU prefix string for proxy-ticket-IOU
	PrefixProxyGrantingTicketIOU = "PGTIOU"
	// PrefixLoginTicket prefix string for login-ticket
	PrefixLoginTicket = "LT"
	// PrefixTicketGrantingCookie prefix string for login-ticket
	PrefixTicketGrantingCookie = "TGC"
)
