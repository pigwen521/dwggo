package arrmap

import (
	"reflect"
)

func InArrayStr(val string, arrs *[]string) bool {
	for _, v := range *arrs {
		if val == v {
			return true
		}
	}
	return false
}

//数组
func InArrayInt(val int, arrs *[]int) bool {
	for _, v := range *arrs {
		if val == v {
			return true
		}
	}
	return false
}

//是否为数组
func IsArray(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Slice:
		return true
	}
	return false
}

//是否为数组
func IsMap(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		return true
	}
	return false
}

//是否为数组
func IsPointer(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Ptr, reflect.UnsafePointer:
		return true
	}
	return false
}

//删除元素
func DelByIndex_Str(index int, arr *[]string) *[]string {
	*arr = append((*arr)[:index], (*arr)[index+1:]...)
	return arr
}
