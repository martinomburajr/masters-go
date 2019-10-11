package utils

import (
	"fmt"
	"math"
	"reflect"
)

func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func ConvertToFloat64(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float32:
		return float64(i), nil
	case float64:
		return i, nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	default:
		return float64(math.NaN()), fmt.Errorf("cannot convert to float64")
	}
}
