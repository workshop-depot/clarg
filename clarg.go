// Package clarg provides simple Sub-Commands for Go using flag package
// - can be used via copy/paste too, because it's small and it's just some functions!
package clarg

import (
	"flag"
	"os"
	"reflect"
)

// Parse parses arguments for list of commands.
// The first command is the default command (top level args) and can be nil.
func Parse(commands ...*flag.FlagSet) error {
	return parse(os.Args[1:], commands...)
}

func parse(args []string, commands ...*flag.FlagSet) error {
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
	cmdTable := make(map[string]*flag.FlagSet)
	for _, cmd := range commands {
		cmdTable[name(cmd)] = cmd
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

func name(c *flag.FlagSet) string {
	return reflect.ValueOf(c).Elem().FieldByName("name").String()
}
