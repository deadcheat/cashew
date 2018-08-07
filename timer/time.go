package timer

import "time"

// Local Local variable for TimeWrapper
var Local TimeWrapper

func init() {
	Local = new(LocalTimer)
}

// TimeWrapper is a wrapper for time
type TimeWrapper interface {
	Now() time.Time
}

// LocalTimer implement TimeWrapper
type LocalTimer struct{}

// Now return time.Now() simply
func (l *LocalTimer) Now() time.Time {
	return time.Now()
}
