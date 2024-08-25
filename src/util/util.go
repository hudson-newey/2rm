package util

func InArray[T comparable](arr []T, query T) bool {
	for _, value := range(arr) {
		if (value == query) {
			return true
		}
	}

	return false
}

