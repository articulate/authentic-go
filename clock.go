package authentic

import (
	"time"
)

type (
	Clock interface {
		IsBeforeNow(time.Time) bool
	}

	clock struct{}
)

func (c *clock) IsBeforeNow(t time.Time) bool {
	return time.Now().UTC().Before(t)
}
