package cli

import (
	"fmt"
)

type Param struct {
	Binding     Binding
	Brief       string
	Description string
	Max         int
	Min         int
	Name        string
}

func (p *Param) Init() {
	if p.Name == "" {
		panic(fmt.Errorf("param without a name"))
	}
}
