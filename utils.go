package validation

import "reflect"

func IsZero(val interface{}) bool {
	return reflect.ValueOf(val).IsZero()
}

func IsString(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.String
}

func IsInt(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int
}
