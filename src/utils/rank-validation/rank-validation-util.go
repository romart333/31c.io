package rankValidationUtil

func IsValidRank(rank uint64) bool {
	return rank > 0 && rank <= 100
}
