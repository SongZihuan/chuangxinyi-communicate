package utils

func StringIn(dst string, lst []string) bool {
	for _, l := range lst {
		if l == dst {
			return true
		}
	}

	return false
}

func TruncateString(s string, n int64) string {
	if len(s) <= int(n) {
		return s
	}
	return s[:n]
}
