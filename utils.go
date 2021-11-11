package validation

import (
	"fmt"
	"reflect"
)

func IsEmpty(val interface{}) bool {
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)
	switch t.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
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

func IsUint(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint
}

func IsUint8(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint8
}

func IsUint16(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint16
}

func IsUint32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint32
}

func IsUint64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint64
}

func IsFloat32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Float32
}

func IsFloat64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Float64
}

func IsComplex64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Complex64
}

func IsComplex128(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Complex128
}

func IsBool(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Bool
}

func IsSlice(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Slice
}

func IsArray(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Array
}

func IsStruct(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Struct
}

func IsMap(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Map
}

func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
