package utils

func Difference(a, b []string) []string {
	mb := make(map[string]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}

	var diff []string
	for _, x := range a {
		if !mb[x] {
			diff = append(diff, x)
		}
	}

	return diff
}

func Remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
