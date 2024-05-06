package api

func CommonElems(first, second []string) bool {
	set := make(map[string]bool)
	for _, e := range first {
		set[e] = true
	}
	for _, e := range second {
		_, ok := set[e]
		if ok {
			return true
		}
	}
	return false
}
