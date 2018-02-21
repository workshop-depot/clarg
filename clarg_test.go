package clarg

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this is just one way to do it. the flag back variables can be anywhere.

var cmdDefault struct {
	*flag.FlagSet
	count int
	data  string
}

var cmdList struct {
	*flag.FlagSet
	age  int
	name string
}

var cmdSend struct {
	*flag.FlagSet
	dst     string
	payload string
}

func prepCmd() {
	cmdDefault.FlagSet = flag.NewFlagSet("", flag.ExitOnError)
	cmdList.FlagSet = flag.NewFlagSet("list", flag.ExitOnError)
	cmdSend.FlagSet = flag.NewFlagSet("send", flag.ExitOnError)

	cmdDefault.IntVar(&cmdDefault.count, "cnt", 0, "")
	cmdDefault.StringVar(&cmdDefault.data, "data", "", "")

	cmdList.IntVar(&cmdList.age, "age", -1, "")
	cmdList.StringVar(&cmdList.name, "name", "", "")

	cmdSend.StringVar(&cmdSend.dst, "dst", "", "")
	cmdSend.StringVar(&cmdSend.payload, "p", "", "")
}

func TestNoArg(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	prepCmd()

	n, err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.NoError(err)
	assert.Equal("", n)
}

func TestTopArg(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66"}
	prepCmd()

	n, err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.Equal("", n)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
}

func TestMultipleCommands01(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66",
		"list", "-age", "20", "-name", "Kaveh"}
	prepCmd()

	n, err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.Equal("list", n)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
	assert.Equal(20, cmdList.age)
	assert.Equal("Kaveh", cmdList.name)
}

func TestMultipleCommands02(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66",
		"send", "-dst", "10", "-p", "QWERTY"}
	prepCmd()

	n, err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.Equal("send", n)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
	assert.Equal("10", cmdSend.dst)
	assert.Equal("QWERTY", cmdSend.payload)
}

func TestNonDefined(t *testing.T) {
	assert := assert.New(t)

	args := []string{"hey"}
	prepCmd()

	n, err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.Equal("", n)
	assert.Error(err)
	assert.Contains(err.Error(), "command hey is not defined")
}

func TestNilTop(t *testing.T) {
	assert := assert.New(t)

	args := []string{"list", "-age", "20", "-name", "Kaveh"}
	prepCmd()

	n, err := parse(args, nil, cmdSend.FlagSet, cmdList.FlagSet)
	assert.Equal("list", n)
	assert.NoError(err)
	assert.Equal(20, cmdList.age)
	assert.Equal("Kaveh", cmdList.name)
}

func ExampleParse() {
	topFlags := flag.NewFlagSet("", flag.ExitOnError)
	cmdListFlags := flag.NewFlagSet("list", flag.ExitOnError)
	cmdSendFlags := flag.NewFlagSet("send", flag.ExitOnError)

	var cmdDefault struct {
		count int
		data  string
	}

	var cmdList struct {
		age  int
		name string
	}

	var cmdSend struct {
		dst     string
		payload string
	}

	topFlags.IntVar(&cmdDefault.count, "cnt", 0, "-cnt <count>")
	topFlags.StringVar(&cmdDefault.data, "data", "", "-data <data string>")

	cmdListFlags.IntVar(&cmdList.age, "age", -1, "-age <age>")
	cmdListFlags.StringVar(&cmdList.name, "name", "", "-name <name>")

	cmdSendFlags.StringVar(&cmdSend.dst, "dst", "", "-dst <destination>")
	cmdSendFlags.StringVar(&cmdSend.payload, "p", "", "-p <payload>")

	if name, err := Parse(topFlags,
		cmdListFlags,
		cmdSendFlags); err != nil {
		// show/handle error
	} else {
		_ = name // the name of the command
	}

	// use values of back fields for flags
}

func ExampleParse_env() {
	topFlags := flag.NewFlagSet("", flag.ExitOnError)

	var cmdDefault struct {
		count int
		data  string
	}

	topFlags.IntVar(&cmdDefault.count, "cnt", 0, "-cnt <count>")
	topFlags.StringVar(&cmdDefault.data, "data", "", "-data <data string>")

	if name, err := Parse(topFlags); err != nil {
		// show/handle error
	} else {
		fromenv := func(set *flag.FlagSet, flagName, envName string) error {
			f := set.Lookup(flagName)
			if f == nil {
				return nil
			}
			if f.Value.String() != f.DefValue {
				return nil
			}
			e := os.Getenv(envName)
			if e == "" {
				return nil
			}
			if e == f.DefValue {
				return nil
			}
			return f.Value.Set(e)
		}
		topFlags.VisitAll(func(f *flag.Flag) {
			// sample env var name convention: APPNAME_FLAGNAME
			if err := fromenv(topFlags,
				f.Name,
				strings.ToUpper(filepath.Base(os.Args[0])+"_"+f.Name)); err != nil {
				log.Fatal(err)
			}
		})
		_ = name // the name of the command
	}

	// use values of back fields for flags
}

func ExampleParse_nonArgs() {
	topFlags := flag.NewFlagSet("", flag.ExitOnError)

	var cmdDefault struct {
		count int
		data  string
	}

	topFlags.IntVar(&cmdDefault.count, "cnt", 0, "-cnt <count>")
	topFlags.StringVar(&cmdDefault.data, "data", "", "-data <data string>")

	if name, err := Parse(topFlags); err != nil {
		// show/handle error
	} else {
		nonFlags := topFlags.Args()
		_ = name // the name of the command
		_ = nonFlags
	}

	// use values of back fields for flags
}
