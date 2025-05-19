package utils

func FilterNotIn[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}

	var result = make([]T, 0)
	for _, item := range a {
		if _, found := set[item]; !found {
			result = append(result, item)
		}
	}

	return result
}
