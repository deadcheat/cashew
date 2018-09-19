package tickets

import (
	"github.com/deadcheat/cashew"
)

// ProxyServices aggregate proxies
func ProxyServices(t *cashew.Ticket) []string {
	proxies := make([]string, 0)
	if t == nil {
		return proxies
	}

	if t.Service != "" {
		proxies = append(proxies, t.Service)
	}
	grant := t.GrantedBy
	for grant != nil {
		service := grant.Service
		ticketType := grant.Type
		grant = grant.GrantedBy
		if (ticketType != cashew.TicketTypeService && ticketType != cashew.TicketTypeProxy) || service == "" {
			continue
		}
		proxies = append(proxies, service)
	}
	return proxies
}
