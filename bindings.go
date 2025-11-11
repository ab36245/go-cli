package cli

import "github.com/ab36245/go-bindings"

var (
	Bool        = bindings.Bool
	BoolFlag    = bindings.BoolFlag
	BoolSlice   = bindings.BoolSlice
	Date        = bindings.Date
	DateSlice   = bindings.DateSlice
	Int         = bindings.Int
	IntFlag     = bindings.IntFlag
	IntSlice    = bindings.IntSlice
	String      = bindings.String
	StringSlice = bindings.StringSlice
)

func Enum[T comparable](binding *T) bindings.EnumBinding[T] {
	return bindings.Enum(binding)
}

func EnumSlice[T comparable](binding *[]T) bindings.EnumBinding[T] {
	return bindings.EnumSlice(binding)
}
