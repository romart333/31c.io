package stringUtil

import (
	"strings"
)

func CaseInsensitiveContains(str string, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
