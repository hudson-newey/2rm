package rm

import (
	"hudson-newey/2rm/src/cli"
	"hudson-newey/2rm/src/util"
	"os"
	"strings"
)

func extractFilePaths(input []string) []string {
	filePaths := []string{}

	for _, path := range input {
		if !strings.HasPrefix(path, "-") {
			if util.IsDirectory(path) && !strings.HasSuffix(path, "/") {
				filePaths = append(filePaths, path+"/")
			} else {
				filePaths = append(filePaths, path)
			}
		} else {
			isSupportedCla := false
			for _, argument := range cli.SupportedCliArguments {
				cliFormat := "--" + argument
				shortCliFormat := "-" + argument
				if cliFormat == path || shortCliFormat == path {
					isSupportedCla = true
				}
			}

			if isSupportedCla {
				continue
			}

			errorMessage := "unrecognized option '" + path + "'\nTry '2rm --help' for more information."
			if util.PathExists(path) {
				errorMessage = "unrecognized option '" + path + "'\n" +
					"Try '2rm ./" + path + "' to remove the file '" + path + "'.\n" +
					"Try '2rm --help' for more information."
			}

			// I (personally) don't like this behavior of only erroring one file at a
			// time. However, to maintain compatibility with GNU rm I have replicated
			// this behavior of only showing the first error every run.
			cli.PrintError(errorMessage)
			os.Exit(1)
		}
	}

	return filePaths
}

func isTmpPath(absolutePath string) bool {
	return strings.HasPrefix(absolutePath, "/tmp")
}
