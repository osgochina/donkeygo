package utils

func IsLetterUpper(b byte) bool {
	if b >= byte('A') && b <= byte('Z') {
		return true
	}
	return false
}
