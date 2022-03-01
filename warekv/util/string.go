package util

import "strconv"

func Str2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Str2Int(str string) (int, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	return int(num), err
}

func IsStrInList(str string, list []string) bool {
	for i := range list {
		if list[i] == str {
			return true
		}
	}
	return false
}
