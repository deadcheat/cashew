package cashew

import (
	"time"
)

// Ticket struct for ticket
type Ticket struct {
	ID        string
	Type      TicketType
	Service   string
	CreatedAt time.Time
	ExpiredAt time.Time
}

// TicketType types of ticket
type TicketType int8
