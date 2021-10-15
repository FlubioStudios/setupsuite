package utils

func FindIndex(ar []string, x string, start int) int {
	for i, n := range ar[start] {
		if x == string(n) {
			return i
		}
	}
	return len(ar)
}
