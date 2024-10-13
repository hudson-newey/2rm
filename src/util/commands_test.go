package util_test

import (
	"hudson-newey/2rm/src/util"
	"testing"
)

func TestCommandsEcho(t *testing.T) {
	util.Execute("echo 'Hello World'")
}
