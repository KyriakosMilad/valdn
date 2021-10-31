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

func IsInt8(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int8
}

func IsInt16(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int16
}

func IsInt32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int32
}

func IsInt64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int64
}
