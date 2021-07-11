package util

import "reflect"

type reflectR struct{}

var Reflect reflectR

func (*reflectR) IsEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.String:
		return v.String() == ""
	default:
		panic("未设定值")
	}
	return false
}
