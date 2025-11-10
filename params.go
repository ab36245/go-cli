package cli

import "fmt"

type Params []*Param

func (p *Params) Init() {
	for _, param := range *p {
		param.Init()
	}
}

func (p *Params) Parse(args *[]string) error {
	// for _, param := range *p {
	// 	err := param.Binding.Consume(args, param.Min, param.Max)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	return fmt.Errorf("unimplemented")
}
