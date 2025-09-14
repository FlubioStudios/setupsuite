package utils

func FindIndex(ar []string, x string, start int) int {
	for i := start; i < len(ar); i++ {
		if ar[i] == x {
			return i
		}
	}
	return -1
}
