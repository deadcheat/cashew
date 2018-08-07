package timer

import "time"

// Clock is a wrapper for time
type Clock interface {
	Now() time.Time
}
