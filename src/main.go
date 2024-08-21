package main

import (
	"fmt"
	"os"
	"strings"

	"hudson-newey/2rm/src/commands"
)

func main() {
	args := strings.Join(os.Args[1:], " ")
	args = strings.ReplaceAll(args, " --no-preserve-root", "")

	fmt.Println(args)

	gitCommand := "rm " + args
	commands.Execute(gitCommand)
}
