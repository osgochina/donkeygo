package dstr

import (
	"donkeygo/internal/utils"
	"strings"
)

func Trim(str string, characterMask ...string) string {
	return utils.Trim(str, characterMask...)
}

// TrimLeft strips whitespace (or other characters) from the beginning of a string.
func TrimLeft(str string, characterMask ...string) string {
	trimChars := utils.DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.TrimLeft(str, trimChars)
}

// TrimRight strips whitespace (or other characters) from the end of a string.
func TrimRight(str string, characterMask ...string) string {
	trimChars := utils.DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.TrimRight(str, trimChars)
}
