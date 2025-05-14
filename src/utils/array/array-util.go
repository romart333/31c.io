package arrayUtil

func FilterItems[T any](items []T, conditions []func(T) bool) []T {
	if len(conditions) == 0 {
		return items
	}
	var result []T
	for _, entity := range items {
		matches := true
		for _, condition := range conditions {
			if !condition(entity) {
				matches = false
				break
			}
		}
		if matches {
			result = append(result, entity)
		}
	}
	return result
}

func FindItem[T any](items []T, conditions []func(T) bool) *T {
	var result *T
	for _, entity := range items {
		meetsAllConditions := true
		for _, condition := range conditions {
			if !condition(entity) {
				meetsAllConditions = false
				break
			}
		}
		if meetsAllConditions {
			result = &entity
			break
		}
	}
	return result
}

func ItemExists[T comparable](array []T, item T) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
