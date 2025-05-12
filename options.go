package cli

import (
	"fmt"
	"strings"
)

type Options []*Option

func (o *Options) Init() {
	for _, option := range *o {
		option.Init()
	}
}

func (o *Options) Parse(args *[]string) error {
	for len(*args) > 0 {
		arg := (*args)[0]
		if !strings.HasPrefix(arg, "-") {
			break
		}
		*args = (*args)[1:]
		if arg == "--" {
			break
		}
		var err error
		if strings.HasPrefix(arg, "--") {
			err = o.parseLong(arg[2:], args)
		} else {
			err = o.parseShort(arg[1:], args)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Options) parseLong(arg string, args *[]string) error {
	name := arg
	value := ""
	reset := false
	if strings.HasPrefix(arg, "no-") {
		name = arg[3:]
		reset = true
	} else if pos := strings.IndexRune(arg, '='); pos >= 0 {
		name = arg[:pos]
		value = arg[pos+1:]
	}
	for _, o := range *o {
		if o.Name != name {
			continue
		}
		b := o.Binding
		var err error
		if reset {
			b.Reset()
		} else if value != "" {
			err = b.Assign(value)
		} else if f, ok := b.(OptionFlag); ok {
			f.Update()
		} else if len(*args) > 0 {
			value = (*args)[0]
			*args = (*args)[1:]
			err = b.Assign(value)
		} else {
			err = fmt.Errorf("requires a value")
		}
		if err != nil {
			return fmt.Errorf("--%s: %w", name, err)
		}
		return nil
	}
	return fmt.Errorf("--%s: unknown option", arg)
}

func (o *Options) parseShort(arg string, args *[]string) error {
	for arg != "" {
		short := arg[0:1]
		arg = arg[1:]

		for _, option := range *o {
			if option.Short != short {
				continue
			}
			b := option.Binding
			var err error
			if f, ok := b.(OptionFlag); ok {
				f.Update()
			} else if arg != "" {
				value := arg
				arg = ""
				err = b.Assign(value)
			} else if len(*args) > 0 {
				value := (*args)[0]
				*args = (*args)[1:]
				err = b.Assign(value)
			} else {
				err = fmt.Errorf("requires a value")
			}
			if err != nil {
				return fmt.Errorf("-%s: %w", short, err)
			}
			return nil
		}
		return fmt.Errorf("-%s: unknown option", short)
	}
	return nil
}
