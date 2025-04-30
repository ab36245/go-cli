package cli

import (
	"fmt"
	"os"
	"strings"
)

type Usage struct {
	Command *Command
	Builder strings.Builder
}

func (u *Usage) write() {
	u.description()
	u.header()
	u.options()
	u.params()
	u.subcommands()
	fmt.Fprintf(os.Stderr, "%s", u.Builder.String())
}

func (u *Usage) description() {
	c := u.Command
	if c.Description != "" {
		u.add("%s", c.Description)
		u.add("\n")
	} else if c.Brief != "" {
		u.add("%s", c.Brief)
		u.add("\n")
	} else {
		return
	}
	u.add("\n")
}

func (u *Usage) header() {
	c := u.Command
	u.add("Usage: %s", c.FullName)
	if len(c.Options) > 0 {
		u.add(" [options]")
	}
	if len(c.Params) > 0 {
		for _, p := range c.Params {
			u.add(" %s", p.Name)
		}
	} else if len(c.Subcommands) > 0 {
		u.add(" command [args...]")
	}
	u.add("\n")
}

func (u *Usage) options() {
	os := &u.Command.Options
	if len(*os) == 0 {
		return
	}

	u.add("  Options\n")
	for _, o := range *os {
		ou := o.usage()
		u.add("    %s", ou.long)
		if ou.short != "" {
			u.add(", %s", ou.short)
		}
		u.add(" <%s>", ou.kind)
		if ou.dflt != "" {
			u.add(" [default %s]", ou.dflt)
		}
		u.add("\n")
		if ou.description != "" {
			u.add("      %s\n", ou.description)
		}
	}
}

func (u *Usage) params() {
	ps := &u.Command.Params
	if len(*ps) == 0 {
		return
	}

	u.add("  Params\n")
	for _, p := range *ps {
		pu := p.usage()
		u.add("    %s", pu.name)
		u.add(" <%s>", pu.kind)
		if pu.dflt != "" {
			u.add(" [default %s]", pu.dflt)
		}
		u.add("\n")
		if pu.description != "" {
			u.add("      %s\n", pu.description)
		}
	}
}

func (u *Usage) subcommands() {
	cs := &u.Command.Subcommands
	if len(*cs) == 0 {
		return
	}

	u.add("  Sub-commands\n")
	maxSize := 0
	for _, c := range *cs {
		size := len([]rune(c.Name))
		if maxSize < size {
			maxSize = size
		}
	}
	for _, c := range *cs {
		u.add("    %-*s", maxSize, c.Name)
		if c.Brief != "" {
			u.add("  %s", c.Brief)
		}
		if u.Command.Default != nil && c.Name == u.Command.Default.Name {
			u.add(" (default)")
		}
		u.add("\n")
	}
}

func (u *Usage) add(mesg string, args ...any) {
	u.Builder.WriteString(fmt.Sprintf(mesg, args...))
}
