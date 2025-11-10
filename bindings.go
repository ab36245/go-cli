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

func Enum[T any](mapping map[string]T, binding *T) bindings.Binding {
	return bindings.Enum(mapping, binding)
}

func EnumSlice[T any](mapping map[string]T, binding *[]T) bindings.Binding {
	return bindings.EnumSlice(mapping, binding)
}
