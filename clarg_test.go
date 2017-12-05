package clarg

import (
	"flag"
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

	err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.NoError(err)
}

func TestName(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	prepCmd()

	err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.NoError(err)
	assert.Equal("", name(cmdDefault.FlagSet))
	assert.Equal("send", name(cmdSend.FlagSet))
	assert.Equal("list", name(cmdList.FlagSet))
}

func TestTopArg(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66"}
	prepCmd()

	err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
}

func TestMultipleCommands01(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66",
		"send", "-dst", "10", "-p", "QWERTY", "list", "-age", "20", "-name", "Kaveh"}
	prepCmd()

	err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
	assert.Equal("10", cmdSend.dst)
	assert.Equal("QWERTY", cmdSend.payload)
	assert.Equal(20, cmdList.age)
	assert.Equal("Kaveh", cmdList.name)
}

func TestMultipleCommands02(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66",
		"list", "-age", "20", "-name", "Kaveh",
		"send", "-dst", "10", "-p", "QWERTY"}
	prepCmd()

	err := parse(args, cmdDefault.FlagSet, cmdSend.FlagSet, cmdList.FlagSet)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
	assert.Equal("10", cmdSend.dst)
	assert.Equal("QWERTY", cmdSend.payload)
	assert.Equal(20, cmdList.age)
	assert.Equal("Kaveh", cmdList.name)
}

func ExampleParse() {
	cmdDefaultFlags := flag.NewFlagSet("", flag.ExitOnError)
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

	cmdDefaultFlags.IntVar(&cmdDefault.count, "cnt", 0, "")
	cmdDefaultFlags.StringVar(&cmdDefault.data, "data", "", "")

	cmdListFlags.IntVar(&cmdList.age, "age", -1, "")
	cmdListFlags.StringVar(&cmdList.name, "name", "", "")

	cmdSendFlags.StringVar(&cmdSend.dst, "dst", "", "")
	cmdSendFlags.StringVar(&cmdSend.payload, "p", "", "")

	if err := Parse(cmdDefaultFlags,
		cmdListFlags,
		cmdSendFlags); err != nil {
		// show/handle error
	}

	// use values of back fields for flags
}
