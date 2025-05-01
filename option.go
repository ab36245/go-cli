package cli

import "fmt"

type Option struct {
	Binding     any
	Description string
	Flag        bool
	Name        string
	Short       string

	value *value
}

type optionUsage struct {
	description string
	dflt        string
	kind        string
	long        string
	short       string
}

func (o *Option) begin() error {
	if o.Name == "" {
		return fmt.Errorf("option without a name")
	}
	value, err := newValue(o.Binding, o.Flag)
	if err != nil {
		return fmt.Errorf("--\"%s\" option error: %w", o.Name, err)
	}
	o.value = value
	return nil
}

func (o *Option) long(arg string, args *[]string) error {
	var err error
	if arg != "" {
		err = o.value.set(arg)
	} else if o.value.inc() {
		// err = nil
	} else if len(*args) > 0 {
		err = o.value.set(shiftArgs(args))
	} else {
		err = fmt.Errorf("requires a value")
	}
	if err != nil {
		err = fmt.Errorf("option \"--%s\" %s", o.Name, err)
	}
	return err
}

func (o *Option) short(arg *string, args *[]string) error {
	var err error
	if o.value.inc() {
		// err = nil
	} else if *arg != "" {
		err = o.value.set(shiftArg(arg))
	} else if len(*args) > 0 {
		err = o.value.set(shiftArgs(args))
	} else {
		err = fmt.Errorf("requires a value")
	}
	if err != nil {
		err = fmt.Errorf("option \"-%s\" %s", o.Short, err)
	}
	return err
}

func (o *Option) usage() optionUsage {
	u := optionUsage{}
	u.description = o.Description
	u.dflt = o.value.dflt
	u.kind = o.value.kind
	u.long = fmt.Sprintf("--%s", o.Name)
	if o.Short != "" {
		u.short = fmt.Sprintf("-%s", o.Short)
	}
	return u
}
