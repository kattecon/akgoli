package utils

import "strings"

// Remove occupancies of the given string from the slice. Mutates original array.
func RemoveStrFromSliceInPlace(slice []string, target string) []string {
	j := 0

	for i, s := range slice {
		if s != target {
			if i != j {
				slice[j] = s
			}
			j += 1
		}
	}

	if len(slice) == j {
		return slice
	} else {
		return slice[:j]
	}
}

// Trim slice.
func TrimAllInPlace(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}

	return slice
}
