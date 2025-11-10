package cli

import (
	"fmt"

	"github.com/ab36245/go-bindings"
)

type Option struct {
	Binding     OptionBinding
	Description string
	Name        string
	Short       string

	OnUsage func(option Option)

	defaultValue string
}

type OptionBinding interface {
	bindings.Binding
}

type OptionFlag interface {
	Update()
}

func (o *Option) Init() {
	if o.Name == "" {
		panic(fmt.Errorf("option without a name"))
	}
	if !o.Binding.IsZero() {
		o.defaultValue = o.Binding.String()
	}
}
