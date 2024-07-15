package utils

import "strconv"

type Filterable interface {
	int | string
}

func Unique[T Filterable](items []T) (list []T) {
	var keys = make(map[T]bool, len(items))

	for _, item := range items {
		if _, exists := keys[item]; !exists {
			keys[item] = true
			list = append(list, item)
		}
	}

	return
}

// Convert a string to an int, without the nasty error-check
func ConvStrToI(s string) (i int) {
	i, err := strconv.Atoi(s)
	check(err)
	return
}
