package cashew

import "time"

type Ticket struct {
	ID        string
	Type      string
	Service   string
	CreatedAt time.Time
}
