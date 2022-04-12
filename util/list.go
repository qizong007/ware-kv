package util

func IsIntInList(list []int, num int) bool {
	for i := range list {
		if list[i] == num {
			return true
		}
	}
	return false
}
