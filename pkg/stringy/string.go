package stringy

// Concatenate to concatenate string
func Concatenate(strs ...string) string {
	var result string

	for _, str := range strs {
		result += str
	}

	return result
}
