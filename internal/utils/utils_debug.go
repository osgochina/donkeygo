package utils

import "donkeygo/internal/command"

const (
	debugKey                 = "dk.debug"                        // Debug key for checking if in debug mode.
	StackFilterKeyForGoFrame = "/github.com/osgochina/donkeygo/" // Stack filtering key for all GoFrame module paths.
)

var (
	// isDebugEnabled marks whether GoFrame debug mode is enabled.
	isDebugEnabled = false
)

func init() {
	// Debugging configured.
	value := command.GetOptWithEnv(debugKey)
	if value == "" || value == "0" || value == "false" {
		isDebugEnabled = false
	} else {
		isDebugEnabled = true
	}
}

// IsDebugEnabled checks and returns whether debug mode is enabled.
// The debug mode is enabled when command argument "dk.debug" or environment "DK_DEBUG" is passed.
func IsDebugEnabled() bool {
	return isDebugEnabled
}
