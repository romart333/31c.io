package nameValidationUtil

func IsValidName(name string) bool {
	return name != "" && len(name) <= 255
}
