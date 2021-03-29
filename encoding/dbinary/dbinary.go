package dbinary

func DecodeToInt64(b []byte) int64 {
	return LeDecodeToInt64(b)
}

func DecodeToFloat64(b []byte) float64 {
	return LeDecodeToFloat64(b)
}
