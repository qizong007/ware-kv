package util

import (
	"fmt"
	"reflect"
)

func GetRealSizeOf(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func sizeof(v reflect.Value) int {
	switch v.Kind() {

	case reflect.Map:
		sum := 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapKey := keys[i]
			s := sizeof(mapKey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapKey))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Slice, reflect.Array:
		n := v.Len()
		if n == 0 {
			return 0
		}
		return sizeof(v.Index(0)) * n

	case reflect.String:
		return v.Len()

	case reflect.Ptr, reflect.Interface, reflect.Uintptr, reflect.UnsafePointer:
		if v.IsNil() {
			return 0
		}
		return sizeof(v.Elem())

	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int, reflect.Uint, reflect.Bool:
		return int(v.Type().Size())

	default:
		fmt.Println("t.Kind() no found:", v.Kind())
	}

	return -1
}
