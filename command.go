package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Command struct {
	Brief       string
	Default     *Command
	Description string
	Logging     bool
	LogFile     string
	LogLevel    slog.Level
	Name        string
	Options     Options
	Params      Params
	Subcommands []Command

	OnInit  func() error
	OnRun   func(command *Command, args []string)
	OnSetup func(command *Command)
	OnUsage func(command *Command)

	FullName string
}

func (c Command) Error(mesg string, args ...any) {
	if len(args) > 0 {
		mesg = fmt.Sprintf(mesg, args...)
	}
	for _, line := range strings.Split(mesg, "\n") {
		fmt.Fprintf(os.Stderr, "%s: %s\n", c.FullName, line)
	}
}

func (c Command) Fatal(code int, mesg string, args ...any) {
	c.Error(mesg, args...)
	os.Exit(code)
}

func (c Command) Ok() {
	os.Exit(0)
}

func (c Command) Panic(mesg string, args ...any) {
	c.Fatal(254, mesg, args...)
}

func (c Command) Run() {
	c.RunWithArgs(os.Args[1:])
}

func (c *Command) RunWithArgs(args []string) {
	if c.Name == "" {
		full, err := os.Executable()
		if err != nil {
			full = os.Args[0]
		}
		path, name := filepath.Split(full)
		if strings.HasPrefix(name, "__debug_") {
			_, name = filepath.Split(path[:len(path)-1])
		}
		c.Name = name
	}
	c.FullName = c.Name
	c.run(args)
}

func (c *Command) Usage() {
	s := CommandUsage(c)
	fmt.Fprintf(os.Stderr, "%s", s)
	os.Exit(2)
}

func (c *Command) run(args []string) {
	if c.OnInit != nil {
		if err := c.OnInit(); err != nil {
			c.Fatal(1, "%s", err)
		}
	}

	if c.Options == nil {
		c.Options = []*Option{}
	}
	if c.Params == nil {
		c.Params = []*Param{}
	}

	if c.Logging {
		levels := map[string]slog.Level{
			"debug": slog.LevelDebug,
			"error": slog.LevelError,
			"info":  slog.LevelInfo,
			"none":  999,
			"warn":  slog.LevelWarn,
		}
		c.Options = append(c.Options, &Option{
			Binding:     String(&c.LogFile),
			Description: "Set the logging output file",
			Name:        "log-file",
		})
		c.Options = append(c.Options, &Option{
			Binding:     Enum(&c.LogLevel, levels),
			Description: "Set the logging level",
			Name:        "log-level",
		})
	}

	help := false
	c.Options = append(c.Options, &Option{
		Binding:     BoolFlag(&help),
		Description: "Show this usage message",
		Name:        "help",
		Short:       "h",
	})
	c.Options.Init()
	c.Params.Init()

	err := c.Options.Parse(&args)
	if err != nil {
		c.Error("%s", err)
		c.Usage()
	}

	if help {
		c.Usage()
	}

	if c.Params != nil {
		err := c.Params.Parse(&args)
		if err != nil {
			c.Error("%s", err)
			c.Usage()
		}
	}

	if c.Logging {
		var enable = true
		if c.LogLevel == 999 {
			enable = false
		}
		var file *os.File
		var err error
		if enable {
			switch c.LogFile {
			case "":
				enable = false
			case "err", "stderr":
				file = os.Stderr
			case "out", "stdout":
				file = os.Stdout
			default:
				file, err = os.Create(c.LogFile)
				if err != nil {
					c.Error("error opening log file \"%s\": %s", c.LogFile, err)
					c.Error("logging is disabled")
				}
			}
		}
		if enable {
			logger := slog.New(slog.NewJSONHandler(file, nil))
			slog.SetDefault(logger)
			slog.SetLogLoggerLevel(c.LogLevel)
		}
	}

	if c.OnSetup != nil {
		c.OnSetup(c)
	}

	if c.OnRun != nil {
		c.runRunner(args)
	} else if len(c.Subcommands) > 0 {
		c.runCommand(args)
	} else {
		c.Panic("no runner or sub-commands registered")
	}
}

func (c *Command) runCommand(args []string) {
	if len(c.Subcommands) == 0 {
		c.Panic("no sub-commands registered")
	}
	var subcommand *Command
	if len(args) > 0 {
		name := args[0]
		args = args[1:]
		for _, sc := range c.Subcommands {
			if sc.Name == name {
				subcommand = &sc
			}
		}
		if subcommand == nil {
			c.Fatal(2, "unknown command \"%s\"", name)
		}
	} else if c.Default != nil {
		subcommand = c.Default
	} else {
		c.Usage()
	}
	subcommand.FullName = c.FullName + " " + subcommand.Name
	subcommand.run(args)
}

func (c *Command) runRunner(args []string) {
	c.OnRun(c, args)
}

func (c *Command) parseOptions(args *[]string) error {
	for _, o := range c.Options {
		o.defaultValue = o.Binding.String()
	}
	for len(*args) > 0 {
		arg := (*args)[0]
		if arg == "--" {
			*args = (*args)[1:]
			break
		}
		if strings.HasPrefix(arg, "--") {
			*args = (*args)[1:]
			if err := c.parseLongOption(arg[2:], args); err != nil {
				return err
			}
		} else if strings.HasPrefix(arg, "-") {
			*args = (*args)[1:]
			if err := c.parseShortOptions(arg[1:], args); err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (c *Command) parseLongOption(arg string, args *[]string) error {
	name := arg
	value := ""
	reset := false
	if strings.HasPrefix(arg, "no-") {
		name = arg[3:]
		reset = true
	} else {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			name = parts[0]
			value = parts[1]
		}
	}

	for _, o := range c.Options {
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

func (c *Command) parseShortOptions(arg string, args *[]string) error {
	for arg != "" {
		if err := c.parseShortOption(&arg, args); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) parseShortOption(arg *string, args *[]string) error {
	short := string((*arg)[0])
	*arg = (*arg)[1:]
	for _, o := range c.Options {
		if o.Short != short {
			continue
		}
		b := o.Binding
		var err error
		if f, ok := b.(OptionFlag); ok {
			f.Update()
		} else if *arg != "" {
			value := *arg
			*arg = ""
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
