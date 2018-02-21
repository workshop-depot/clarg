// Package clarg provides simple Sub-Commands for Go using flag package
// - can be used via copy/paste too, because it's small and it's just some functions!
package clarg

import (
	"flag"
	"fmt"
	"os"
)

// Parse parses arguments for list of commands.
// The first command is the default command (top level args) and can be nil.
// Non-Flag args are available via matched FlagSet's Args() method.
func Parse(top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
	return parse(os.Args[1:], top, subs...)
}

func parse(args []string, top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
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
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", cmd.Name())
			cmd.PrintDefaults()
		}
	}
	if err := top.Parse(args); err != nil {
		return nil, err
	}
	args = top.Args()
	if len(args) == 0 {
		return top, nil
	}
	cmdTable := make(map[string]*flag.FlagSet)
	for _, cmd := range subs {
		cmdTable[cmd.Name()] = cmd
	}
	cmd, found := cmdTable[args[0]]
	if !found {
		return nil, fmt.Errorf("command %v is not defined", args[0])
	}
	if err := cmd.Parse(args[1:]); err != nil {
		return nil, err
	}
	return cmd, nil
}
