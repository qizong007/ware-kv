package util

import "strconv"

func Str2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Str2Int(str string) (int, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	return int(num), err
}

func Str2Uint(str string) (uint, error) {
	num, err := strconv.ParseUint(str, 10, 64)
	return uint(num), err
}

func Str2Float64(str string) (float64, error) {
	num, err := strconv.ParseFloat(str, 64)
	return num, err
}
