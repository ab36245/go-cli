package cli

type Params []*Param

func (p *Params) Init() {
	for _, param := range *p {
		if err := param.begin(); err != nil {
			panic(err)
		}
	}
}

func (p *Params) Parse(args *[]string) error {
	for _, param := range *p {
		err := param.parse(args)
		if err != nil {
			return err
		}
	}
	return nil
}
