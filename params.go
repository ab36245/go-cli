package cli

import "fmt"

type Params []*Param

func (p *Params) Init() {
	for _, param := range *p {
		if param.Name == "" {
			panic(fmt.Errorf("param without a name"))
		}
	}
}

func (p *Params) Parse(args *[]string) error {
	// for _, param := range *p {
	// 	err := param.parse(args)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
