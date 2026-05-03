package cli

import (
	"fmt"

	"github.com/ab36245/go-writer"
)

func CommandUsage(c *Command) string {
	w := writer.New()
	CommandDescription(w, c)
	CommandSummary(w, c)
	OptionsUsage(w, c.Options)
	ParamsUsage(w, c.Params)
	SubcommandsUsage(w, c.Subcommands)
	return w.String()
}

func CommandDescription(w *writer.Writer, c *Command) {
	if c.Description != "" {
		w.End("%s", c.Description)
	} else if c.Brief != "" {
		w.End("%s", c.Brief)
	}
}

func CommandSummary(w *writer.Writer, c *Command) {
	w.Add("Usage: %s", c.FullName)
	if len(c.Options) > 0 {
		w.Add(" [options]")
	}
	if len(c.Params) > 0 {
		for _, p := range c.Params {
			w.Add(" ")
			if p.Min == 0 {
				w.Add("[")
			}
			w.Add("%s", p.Name)
			if p.Min == 0 {
				w.Add("]")
			}
		}
	} else if len(c.Subcommands) > 0 {
		w.Add(" command [args...]")
	}
	w.End("")
}

func OptionsUsage(w *writer.Writer, os Options) {
	if len(os) == 0 {
		return
	}
	w.Over("")
	{
		w.Over("Options")
		{
			for _, o := range os {
				OptionUsage(w, o)
			}
		}
		w.Back("")
	}
	w.Back("")
}

func OptionUsage(w *writer.Writer, o *Option) {
	long := fmt.Sprintf("--%s", o.Name)
	b := o.Binding
	if b.IsFlag() {
		long += "["
	}
	long += fmt.Sprintf("=<%s>", b.Type())
	if b.IsFlag() {
		long += "]"
	}

	short := ""
	if o.Short != "" {
		short = fmt.Sprintf("-%s", o.Short)
		if !b.IsFlag() {
			short += fmt.Sprintf("<%s>", b.Type())
		}
	}

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

func ParamsUsage(w *writer.Writer, ps Params) {
	if len(ps) == 0 {
		return
	}
	w.Over("")
	{
		w.Over("Params")
		{
			width := 0
			for _, p := range ps {
				if width < len(p.Name) {
					width = len(p.Name)
				}
			}
			for _, p := range ps {
				ParamUsage(w, p, width)
			}
		}
		w.Back("")
	}
	w.Back("")
}

func ParamUsage(w *writer.Writer, p *Param, width int) {
	brief := p.Brief
	if brief == "" {
		brief = p.Description
	}
	w.End("%-*s %s", width, p.Name, brief)
}

func SubcommandsUsage(w *writer.Writer, cs []Command) {
	if len(cs) == 0 {
		return
	}
	w.Over("")
	{
		w.Over("Sub-commands")
		{
			width := 0
			for _, c := range cs {
				if width < len(c.Name) {
					width = len(c.Name)
				}
			}
			for _, c := range cs {
				SubcommandUsage(w, c, width)
			}
		}
		w.Back("")
	}
	w.Back("")
}

func SubcommandUsage(w *writer.Writer, c Command, width int) {
	brief := c.Brief
	if brief == "" {
		brief = c.Description
	}
	w.End("%-*s %s", width, c.Name, brief)
}
