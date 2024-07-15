package utils

// Panic upon error
func check(e error) {
	if e != nil {
		panic(e)
	}
}
