package strings

func Split(s string, sep string) []string {
	result := make([]string, 0)
	lastIndex := 0
	for i, item := range s {
		if string(item) != sep {
			continue
		}
		curr := s[lastIndex:i]
		if curr != "" {
			result = append(result, curr)
		}
		lastIndex = i + len(sep)
	}
	if lastIndex < len(s) && s[lastIndex:] != "" {
		result = append(result, s[lastIndex:])
	}
	return result
}
