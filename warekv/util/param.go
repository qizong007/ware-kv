package util

func SetIfHitLimit(param int, min int, max int) int {
	if param < min {
		return min
	}
	if param > max {
		return max
	}
	return param
}
