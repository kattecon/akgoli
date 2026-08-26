package utils

import "strings"

// RemoveStrFromSliceInPlace removes all occurrences of target from slice. It
// compacts matching elements in place within the existing backing array and
// returns a shorter slice header. Elements beyond the returned length are
// still present in the backing array. The caller must use the returned slice,
// not the original, to see the correct result. The relative order of the
// remaining elements is preserved.
func RemoveStrFromSliceInPlace(slice []string, target string) []string {
	j := 0

	for i, s := range slice {
		if s != target {
			// Skip the write when i == j (no elements removed yet).
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

// TrimAllInPlace applies strings.TrimSpace to every element of slice,
// replacing each element in place. It returns the same slice header with
// the same length. Callers that hold references to the same backing array
// will see the trimmed values.
func TrimAllInPlace(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}

	return slice
}
