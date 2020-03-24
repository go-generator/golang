package utils

func IsContainedInStrings(s string, input []string) bool {
	for _, v := range input {
		if v == s {
			return true
		}
	}
	return false
}
