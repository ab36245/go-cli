package cli

import (
	"fmt"

	"github.com/ab36245/go-writer"
)

var UsageWriter = writer.New()

func CommandUsage(c *Command) string {
	CommandDescription(c)
	CommandSummary(c)
	OptionsUsage(c.Options)
	return UsageWriter.String()
}

func CommandDescription(c *Command) {
	w := UsageWriter
	if c.Description != "" {
		w.End(c.Description)
	} else if c.Brief != "" {
		w.End(c.Brief)
	}
}

func CommandSummary(c *Command) {
	w := UsageWriter
	w.Add("Usage: %s", c.FullName)
	if len(c.Options) > 0 {
		w.Add(" [options]")
	}
	if len(c.Params) > 0 {
		for _, p := range c.Params {
			w.Add(" %s", p.Name)
		}
	} else if len(c.Subcommands) > 0 {
		w.Add(" command [args...]")
	}
	w.End("")
}

func OptionsUsage(os Options) {
	if len(os) == 0 {
		return
	}
	w := UsageWriter
	w.Over("")
	{
		w.Over("Options")
		{
			for _, o := range os {
				OptionUsage(o)
			}
		}
		w.Back("")
	}
	w.Back("")
}

func OptionUsage(o *Option) {
	long := fmt.Sprintf("--%s", o.Name)
	b := o.Binding
	_, flag := b.(OptionFlag)
	if flag {
		long += "["
	}
	long += fmt.Sprintf("=<%s>", b.Type())
	if flag {
		long += "]"
	}

	short := ""
	if o.Short != "" {
		short = fmt.Sprintf("-%s", o.Short)
		if !flag {
			short += fmt.Sprintf("<%s>", b.Type())
		}
	}

	w := UsageWriter
	w.Add("%s", long)
	if short != "" {
		w.Add(", %s", short)
	}
	w.Over("")
	{
		if o.Description != "" {
			w.End("%s", o.Description)
		}
		if o.defaultValue != "" {
			w.End("Default: %s", o.defaultValue)
		}
	}
	w.Back("")
}
