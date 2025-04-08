package cli

import (
	"fmt"
	"strings"
)

type Options []*Option

func (o *Options) Init() {
	for _, option := range *o {
		if err := option.begin(); err != nil {
			panic(err)
		}
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
	pos := strings.IndexRune(arg, '=')
	var name string
	if pos < 0 {
		name = arg
		arg = ""
	} else {
		name = arg[:pos]
		arg = arg[pos+1:]
	}
	for _, option := range *o {
		if option.Name == name {
			return option.long(arg, args)
		}
	}
	return fmt.Errorf("unknown option \"--%s\"", name)
}

func (o *Options) parseShort(arg string, args *[]string) error {
	for arg != "" {
		short := arg[0:1]
		arg = arg[1:]

		found := false
		var err error
		for _, option := range *o {
			if option.Short == short {
				found = true
				err = option.short(&arg, args)
				break
			}
		}
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("unknown option \"-%s\"", short)
		}
	}
	return nil
}
