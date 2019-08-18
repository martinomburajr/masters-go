package utils

import (
	"fmt"
	"math"
	"reflect"
)

func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func ConvertToFloat(unk interface{}) (float32, error) {
	switch i := unk.(type) {
	case float32:
		return i, nil
	case float64:
		return float32(i), nil
	case int64:
		return float32(i), nil
	case int32:
		return float32(i), nil
	case int:
		return float32(i), nil
	case uint64:
		return float32(i), nil
	case uint32:
		return float32(i), nil
	case uint:
		return float32(i), nil
	default:
		return float32(math.NaN()), fmt.Errorf("cannot convert to float32")
	}
}
