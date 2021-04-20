package dstr

import "donkeygo/internal/utils"

func Trim(str string, characterMask ...string) string {
	return utils.Trim(str, characterMask...)
}
