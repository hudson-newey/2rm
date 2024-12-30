package cli

import (
	"fmt"
	"os"
)

func PrintError(error string) {
	message := "2rm: " + error
	fmt.Fprintln(os.Stderr, message)
}

func PrintErrorValue(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func Pause() {
	fmt.Scanln()
}
