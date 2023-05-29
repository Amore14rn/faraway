package timestamp

import "time"

type SystemTime struct {
}

func (s SystemTime) Now() time.Time {
	return time.Now()
}
