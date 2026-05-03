package cli

import (
	"fmt"
)

type Option struct {
	Binding     Binding
	Description string
	Name        string
	Short       string

	OnUsage func(option Option)

	defaultValue string
}

func (o *Option) Init() {
	if o.Name == "" {
		panic(fmt.Errorf("option without a name"))
	}
	if !o.Binding.IsZero() {
		o.defaultValue = o.Binding.String()
	}
}
