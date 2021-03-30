package dconv

type apiString interface {
	String() string
}

type apiError interface {
	Error() string
}

type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}
