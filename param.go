package cli

import (
	"fmt"

	"github.com/ab36245/go-bindings"
)

type Param struct {
	Binding     ParamBinding
	Description string
	Max         int
	Min         int
	Name        string
}

type ParamBinding interface {
	bindings.Binding
}

func (p *Param) Init() {
	if p.Name == "" {
		panic(fmt.Errorf("param without a name"))
	}
}
