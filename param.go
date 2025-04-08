package cli

import "fmt"

type Param struct {
	Binding     any
	Description string
	Name        string
	Optional    bool

	value *value
}

type paramUsage struct {
	description string
	dflt        string
	kind        string
	name        string
}

func (p *Param) begin() error {
	if p.Name == "" {
		return fmt.Errorf("param without a name")
	}
	value, err := newValue(p.Binding, false)
	if err != nil {
		return fmt.Errorf("\"%s\" param error: %s", p.Name, err)
	}
	p.value = value
	return nil
}

func (p *Param) parse(args *[]string) error {
	var err error
	if p.value.array {
		for len(*args) > 0 {
			err := p.value.set(shiftArgs(args))
			if err != nil {
				break
			}
		}
	} else if len(*args) > 0 {
		err = p.value.set(shiftArgs(args))
	} else if !p.Optional {
		err = fmt.Errorf("requires a value")
	}
	if err != nil {
		err = fmt.Errorf("param \"%s\" %s", p.Name, err)
	}
	return err
}

func (p *Param) usage() paramUsage {
	u := paramUsage{}
	u.description = p.Description
	u.dflt = p.value.dflt
	u.kind = p.value.kind
	u.name = p.Name
	return u
}
