package cli

type Params []*Param

func (ps *Params) Init() {
	for _, param := range *ps {
		param.Init()
	}
}

func (ps *Params) Parse(args *[]string) error {
	for _, p := range *ps {
		err := p.Binding.Consume(args, p.Min, p.Max)
		if err != nil {
			return err
		}
	}
	return nil
}
