package consts

const (
	// RFC2822 format of time
	RFC2822 = "Mon Jan 02 15:04:05 -0700 2006"
)

const (
	// ParamKeyService key string for param
	ParamKeyService = "service"
	// ParamKeyURL key string for param
	ParamKeyURL = "url"
	// ParamKeyPgtURL key string for param
	ParamKeyPgtURL = "pgtUrl"
	// ParamKeyRenew key string for param
	ParamKeyRenew = "renew"
	// ParamKeyGateway key string for param
	ParamKeyGateway = "gateway"
	// ParamKeyPgt key string for param
	ParamKeyPgt = "pgt"
	// ParamKeyTargetService key string for param
	ParamKeyTargetService = "targetService"
)

const (
	// CookieKeyTGT key string for cookie
	CookieKeyTGT = "tgt"
)

const (
	// TicketTypeStrService type of service ticket
	TicketTypeStrService = "service-ticket"
	// TicketTypeStrProxy type of proxy ticket
	TicketTypeStrProxy = "proxy-ticket"
	// TicketTypeStrTicketGranting type of proxy-granting ticket
	TicketTypeStrTicketGranting = "ticket-granting-ticket"
	// TicketTypeStrProxyGranting type of proxy-granting ticket
	TicketTypeStrProxyGranting = "proxy-granting-ticket"
	// TicketTypeStrProxyGrantingIOU type of proxy-granting ticket iou
	// TODO check is this really necessary
	TicketTypeStrProxyGrantingIOU = "proxy-granting-ticket-iou"
	// TicketTypeStrLogin type of login ticket
	TicketTypeStrLogin = "login-ticket"
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
