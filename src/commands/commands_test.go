package commands

import "testing"

func TestCommandsEcho(t *testing.T) {
	Execute("echo 'Hello World'")
}
