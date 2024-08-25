package main

import (
	"os"

	"hudson-newey/2rm/src/patches"
)

func main() {
	originalArguments := os.Args[1:];
	patches.RmPatch(originalArguments)
}
