package cli

import "fmt"

type Param struct {
	Binding     ParamBinding
	Description string
	Name        string
}

type ParamBinding interface {
	Assign(string) error
	Consume(*[]string) error
	Reset()
	String() string
	Type() string
}

func (p *Param) Init() {
	if p.Name == "" {
		panic(fmt.Errorf("param without a name"))
	}
}
