package dtime

import (
	"sync"
	"time"
)

var (
	// locationMap is time zone name to its location object.
	// Time zone name is like: Asia/Shanghai.
	locationMap = make(map[string]*time.Location)
	// locationMu is used for concurrent safety for `locationMap`.
	locationMu = sync.RWMutex{}
)

// ToLocation 把当前时间转换成指定时区的时间
func (that *Time) ToLocation(location *time.Location) *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.In(location)
	return newTime
}

// ToZone 把当前时间转换成指定时区的时间，如: Asia/Shanghai
func (that *Time) ToZone(zone string) (*Time, error) {
	if location, err := that.getLocationByZoneName(zone); err == nil {
		return that.ToLocation(location), nil
	} else {
		return nil, err
	}
}

func (that *Time) getLocationByZoneName(name string) (location *time.Location, err error) {
	locationMu.RLock()
	location = locationMap[name]
	locationMu.RUnlock()
	if location == nil {
		location, err = time.LoadLocation(name)
		if err == nil && location != nil {
			locationMu.Lock()
			locationMap[name] = location
			locationMu.Unlock()
		}
	}
	return
}

// Local 将时间转换成当前时区
func (that *Time) Local() *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.Local()
	return newTime
}
