package utils

func IsContainedInStrings(s string, input []string) bool {
	for _, v := range input {
		if v == s {
			return true
		}
	}
	return false
}

func CopyMap(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}
