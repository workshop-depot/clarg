// Package clarg provides simple Sub-Commands for Go using flag package
// - can be used via copy/paste too, because it's small.
package clarg

import (
	"flag"
	"os"
)

// Parse parses arguments for list of commands.
// The first command is the default command (top level args) and can be nil.
func Parse(commands ...*Cmd) error {
	return parse(os.Args[1:], commands...)
}

func parse(args []string, commands ...*Cmd) error {
	if len(commands) > 0 && commands[0] != nil {
		if err := commands[0].Parse(args); err != nil {
			return err
		}
		args = commands[0].Args()
	}
	if len(args) == 0 {
		return nil
	}
	commands = commands[1:]
	cmdTable := make(map[string]*Cmd)
	for _, cmd := range commands {
		cmdTable[cmd.cmdName] = cmd
	}
	for len(cmdTable) > 0 {
		if len(args) == 0 {
			return nil
		}
		for name, cmd := range cmdTable {
			if args[0] != name {
				continue
			}
			if err := cmd.Parse(args[1:]); err != nil {
				return err
			}
			args = cmd.Args()
			delete(cmdTable, name)
			break
		}
	}

	return nil
}

// New creates a *Cmd, default value for errorHandling is flag.ExitOnError,
// for top level flags name must be "".
func New(name string, errorHandling ...flag.ErrorHandling) *Cmd {
	eh := flag.ExitOnError
	if len(errorHandling) > 0 {
		eh = errorHandling[0]
	}
	return &Cmd{FlagSet: flag.NewFlagSet(name, eh), cmdName: name}
}

// Cmd represents a sub command
type Cmd struct {
	*flag.FlagSet
	cmdName string // if *flag.FlagSet could give us the name, there would be no need for this field.
}

// Name returns the name
func (c *Cmd) Name() string { return c.cmdName }
