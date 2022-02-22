package util

func IsIntInList(num int, list []int) bool {
	if list == nil {
		return false
	}
	for i := range list {
		if list[i] == num {
			return true
		}
	}
	return false
}
