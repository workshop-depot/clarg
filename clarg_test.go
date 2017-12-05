package clarg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this is just one way to do it. the flag back variables can be anywhere.

var cmdDefault struct {
	*Cmd
	count int
	data  string
}

var cmdList struct {
	*Cmd
	age  int
	name string
}

var cmdSend struct {
	*Cmd
	dst     string
	payload string
}

func prepCmd() {
	cmdDefault.Cmd = New("")
	cmdList.Cmd = New("list")
	cmdSend.Cmd = New("send")

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

	err := parse(args, cmdDefault.Cmd, cmdSend.Cmd, cmdList.Cmd)
	assert.NoError(err)
}

func TestName(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	prepCmd()

	err := parse(args, cmdDefault.Cmd, cmdSend.Cmd, cmdList.Cmd)
	assert.NoError(err)
	assert.Equal("", cmdDefault.Name())
	assert.Equal("send", cmdSend.Name())
	assert.Equal("list", cmdList.Name())
}

func TestTopArg(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66"}
	prepCmd()

	err := parse(args, cmdDefault.Cmd, cmdSend.Cmd, cmdList.Cmd)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
}

func TestMultipleCommands01(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-data", "Hi!", "-cnt", "66",
		"send", "-dst", "10", "-p", "QWERTY", "list", "-age", "20", "-name", "Kaveh"}
	prepCmd()

	err := parse(args, cmdDefault.Cmd, cmdSend.Cmd, cmdList.Cmd)
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

	err := parse(args, cmdDefault.Cmd, cmdSend.Cmd, cmdList.Cmd)
	assert.NoError(err)
	assert.Equal("Hi!", cmdDefault.data)
	assert.Equal(66, cmdDefault.count)
	assert.Equal("10", cmdSend.dst)
	assert.Equal("QWERTY", cmdSend.payload)
	assert.Equal(20, cmdList.age)
	assert.Equal("Kaveh", cmdList.name)
}

func ExampleCmd() {
	cmdDefaultFlags := New("")
	cmdListFlags := New("list")
	cmdSendFlags := New("send")

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
