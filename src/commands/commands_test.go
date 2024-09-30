package commands_test

import (
	"hudson-newey/2rm/src/commands"
	"testing"
)

func TestCommandsEcho(t *testing.T) {
	commands.Execute("echo 'Hello World'")
}
