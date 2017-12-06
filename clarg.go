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
func Parse(top *flag.FlagSet, subs ...*flag.FlagSet) error {
	return parse(os.Args[1:], top, subs...)
}

func parse(args []string, top *flag.FlagSet, subs ...*flag.FlagSet) error {
	if top != nil {
		if err := top.Parse(args); err != nil {
			return err
		}
		args = top.Args()
	}
	if len(args) == 0 {
		return nil
	}
	cmdTable := make(map[string]*flag.FlagSet)
	for _, cmd := range subs {
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
