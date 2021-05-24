package dtime

import "time"

type wrapper struct {
	time.Time
}

func (t wrapper) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
