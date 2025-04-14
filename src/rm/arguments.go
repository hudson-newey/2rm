package rm

import "hudson-newey/2rm/src/util"

func removeDangerousArguments(arguments []string) []string {
	// I have excluded the root slash as a forbidden argument just in case
	// you make a typo like rm ./myDirectory /
	// when you were just trying to delete myDirectory
	// If you really have to delete your root directory consider using the GNU
	// rm command
	forbiddenArguments := []string{"/"}
	returnedArguments := []string{}

	for _, arg := range arguments {
		isForbidden := util.InArray(forbiddenArguments, arg)
		if !isForbidden {
			returnedArguments = append(returnedArguments, arg)
		}
	}

	return returnedArguments
}
