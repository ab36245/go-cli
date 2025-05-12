package cli

import "fmt"

type Option struct {
	Binding     OptionBinding
	Description string
	Name        string
	Short       string

	OnUsage func(option Option)

	defaultValue string
}

type OptionBinding interface {
	Assign(string) error
	NonZero() string
	Reset()
	String() string
	Type() string
}

type OptionFlag interface {
	Update()
}

func (o *Option) Init() {
	if o.Name == "" {
		panic(fmt.Errorf("option without a name"))
	}
	o.defaultValue = o.Binding.NonZero()
}
