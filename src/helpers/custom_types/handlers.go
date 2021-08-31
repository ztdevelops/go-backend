package custom_types

import (
	"reflect"
)

func GetStructName(target interface{}) string {
	return reflect.TypeOf(target).Name()
}