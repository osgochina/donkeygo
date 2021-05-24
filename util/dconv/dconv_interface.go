package dconv

type apiString interface {
	String() string
}

type apiBool interface {
	Bool() bool
}

type apiInt64 interface {
	Int64() int64
}

type apiUint64 interface {
	Uint64() uint64
}

type apiFloat32 interface {
	Float32() float32
}

type apiFloat64 interface {
	Float64() float64
}

type apiError interface {
	Error() string
}

type apiBytes interface {
	Bytes() []byte
}

type apiInterfaces interface {
	Interfaces() []interface{}
}

type apiFloats interface {
	Floats() []float64
}

type apiInts interface {
	Ints() []int
}

type apiStrings interface {
	Strings() []string
}

type apiUints interface {
	Uints() []uint
}

type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

type apiUnmarshalValue interface {
	UnmarshalValue(interface{}) error
}
type apiUnmarshalText interface {
	UnmarshalText(text []byte) error
}

type apiSet interface {
	Set(value interface{}) (old interface{})
}
