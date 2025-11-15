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
	for line := range strings.SplitSeq(mesg, "\n") {
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
		c.Options = append(c.Options, &Option{
			Binding:     String(&c.LogFile),
			Description: "Set the logging output file",
			Name:        "log-file",
		})
		c.Options = append(c.Options, &Option{
			Binding: Enum(&c.LogLevel).
				Map("debug", slog.LevelDebug).
				Map("error", slog.LevelError).
				Map("info", slog.LevelInfo).
				Map("none", 999).
				Map("warn", slog.LevelWarn),
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
	dispatch := func(sc *Command) {
		sc.FullName = c.FullName + " " + sc.Name
		sc.run(args)
	}

	if len(args) > 0 {
		name := args[0]
		args = args[1:]

		// Look for exact match
		for _, sc := range c.Subcommands {
			if sc.Name == name {
				dispatch(&sc)
				return
			}
		}

		// Look for longest substring
		commands := []*Command{}
		for _, sc := range c.Subcommands {
			if strings.HasPrefix(sc.Name, name) {
				commands = append(commands, &sc)
			}
		}
		if len(commands) < 1 {
			c.Fatal(2, "unknown command \"%s\"", name)
		} else if len(commands) > 1 {
			c.Error("ambiguous command name \"%s\"", name)
			c.Error("possible names are:")
			for _, sc := range commands {
				c.Error("- %s", sc.Name)
			}
			c.Fatal(2, "enter a longer name to disambiguate")
		} else {
			dispatch(commands[0])
		}
	} else if c.Default != nil {
		dispatch(c.Default)
	} else {
		c.Usage()
	}
}

func (c *Command) runRunner(args []string) {
	c.OnRun(c, args)
}
