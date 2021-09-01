package custom_types

import (
	"reflect"
)

// GetStructName returns the name of the struct, with the help of the reflect package.
func GetStructName(target interface{}) string {
	return reflect.TypeOf(target).Name()
}