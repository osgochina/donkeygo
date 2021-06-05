package dtime

import "time"

type wrapper struct {
	time.Time
}

func (that wrapper) String() string {
	if that.IsZero() {
		return ""
	}
	return that.Format("2006-01-02 15:04:05")
}
