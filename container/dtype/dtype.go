package dtype

// Type is alias of Interface.
type Type = Interface

// New is alias of NewInterface.
// See NewInterface.
func New(value ...interface{}) *Type {
	return NewInterface(value...)
}
