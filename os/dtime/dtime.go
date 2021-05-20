package dtime

import (
	"fmt"
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/text/dregex"
	"strconv"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	if utils.IsNumeric(s) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(v), nil
	}
	match, err := dregex.MatchString(`^([\-\d]+)[dD](.*)$`, s)
	if err != nil {
		return 0, err
	}
	if len(match) == 3 {
		v, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.ParseDuration(fmt.Sprintf(`%dh%s`, v*24, match[2]))
	}
	return time.ParseDuration(s)
}
