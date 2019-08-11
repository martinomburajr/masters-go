package utils

import "reflect"

func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}