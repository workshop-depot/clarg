// Package clarg provides simple Sub-Commands for Go using flag package
// - can be used via copy/paste too, because it's small and it's just some functions!
package clarg

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

// Parse parses arguments for list of commands.
// The first command is the default command (top level args) and can be nil.
// Non-Flag args are available via matched FlagSet's Args() method.
func Parse(top *flag.FlagSet, subs ...*flag.FlagSet) (string, error) {
	return parse(os.Args[1:], top, subs...)
}

func parse(args []string, top *flag.FlagSet, subs ...*flag.FlagSet) (string, error) {
	if top == nil {
		top = flag.NewFlagSet("", flag.ExitOnError)
	}
	top.Usage = func() {
		cnt := 0
		top.VisitAll(func(*flag.Flag) { cnt++ })
		if cnt > 0 {
			fmt.Fprintf(os.Stderr, "Usage:\n")
			top.PrintDefaults()
		}
		for _, cmd := range subs {
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", name(cmd))
			cmd.PrintDefaults()
		}
	}
	if err := top.Parse(args); err != nil {
		return "", err
	}
	args = top.Args()
	if len(args) == 0 {
		return "", nil
	}
	cmdTable := make(map[string]*flag.FlagSet)
	for _, cmd := range subs {
		cmdTable[name(cmd)] = cmd
	}
	if len(args) == 0 {
		return "", nil
	}
	var lastName string
	for name, cmd := range cmdTable {
		if args[0] != name {
			continue
		}
		if err := cmd.Parse(args[1:]); err != nil {
			return "", err
		}
		args = cmd.Args()
		lastName = name
		break
	}
	if lastName == "" {
		return "", fmt.Errorf("command %v is not defined", args[0])
	}

	return lastName, nil
}

func name(c *flag.FlagSet) string {
	return reflect.ValueOf(c).Elem().FieldByName("name").String()
}
