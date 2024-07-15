package utils

// Convert a string to an int, without the nasty error-check
func Abs(a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}
