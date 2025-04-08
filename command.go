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
	LogLevel    string
	Name        string
	Options     Options
	Params      Params
	Runner      func(command *Command, args []string)
	Setup       func(command *Command)
	Subcommands []Command

	FullName string
}

func (c *Command) Error(mesg string, args ...any) {
	if len(args) > 0 {
		mesg = fmt.Sprintf(mesg, args...)
	}
	for _, line := range strings.Split(mesg, "\n") {
		fmt.Fprintf(os.Stderr, "%s: %s\n", c.FullName, line)
	}
}

func (c *Command) Fatal(code int, mesg string, args ...any) {
	c.Error(mesg, args...)
	os.Exit(code)
}

func (c *Command) Ok() {
	os.Exit(0)
}

func (c *Command) Panic(mesg string, args ...any) {
	c.Fatal(254, mesg, args...)
}

func (c *Command) Run() {
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
	usage := &Usage{Command: c}
	usage.write()
	os.Exit(2)
}

func (c *Command) run(args []string) {
	if c.Options == nil {
		c.Options = []*Option{}
	}
	if c.Params == nil {
		c.Params = []*Param{}
	}

	if c.Logging {
		c.Options = append(c.Options, &Option{
			Binding:     &c.LogFile,
			Description: "Set the logging outfile file",
			Name:        "log-file",
		})
		c.Options = append(c.Options, &Option{
			Binding:     &c.LogLevel,
			Description: "Set the logging level",
			Name:        "log-level",
		})
	}

	help := false
	c.Options = append(c.Options, &Option{
		Binding:     &help,
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
		var level slog.Level
		switch c.LogLevel {
		case "", "none":
			enable = false
		case "debug":
			level = slog.LevelDebug
		case "error":
			level = slog.LevelError
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelError
		default:
			c.Error("unknown log level \"%s\"", c.LogLevel)
			c.Error("logging is disabled")
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
			slog.SetLogLoggerLevel(level)
		}
	}

	if c.Setup != nil {
		c.Setup(c)
	}

	if c.Runner != nil {
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
	c.Runner(c, args)
}
