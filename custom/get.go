package custom


func IsCustom(breed string) bool {
	for _, b := range Breeds() {
		if b == breed {
			return true
		}
	}
	return false
}
