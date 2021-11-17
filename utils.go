package validation

import (
	"fmt"
	"reflect"
)

// IsEmpty reports weather val is empty or not
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

// IsString reports weather val's kind is string or not
func IsString(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.String
}

// IsInt reports weather val's kind is int or not
func IsInt(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int
}

// IsInt8 reports weather val's kind is int8 or not
func IsInt8(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int8
}

// IsInt16 reports weather val's kind is int16 or not
func IsInt16(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int16
}

// IsInt32 reports weather val's kind is int32 or not
func IsInt32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int32
}

// IsInt64 reports weather val's kind is int64 or not
func IsInt64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int64
}

// IsUint reports weather val's kind is uint or not
func IsUint(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint
}

// IsUint8 reports weather val's kind is uint8 or not
func IsUint8(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint8
}

// IsUint16 reports weather val's kind is uint16 or not
func IsUint16(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint16
}

// IsUint32 reports weather val's kind is uint32 or not
func IsUint32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint32
}

// IsUint64 reports weather val's kind is uint64 or not
func IsUint64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint64
}

// IsFloat32 reports weather val's kind is float32 or not
func IsFloat32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Float32
}

// IsFloat64 reports weather val's kind is float64 or not
func IsFloat64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Float64
}

// IsComplex64 reports weather val's kind is complex64 or not
func IsComplex64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Complex64
}

// IsComplex128 reports weather val's kind is complex128 or not
func IsComplex128(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Complex128
}

// IsBool reports weather val's kind is bool or not
func IsBool(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Bool
}

// IsSlice reports weather val's kind is slice or not
func IsSlice(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Slice
}

// IsArray reports weather val's kind is array or not
func IsArray(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Array
}

// IsStruct reports weather val's kind is struct or not
func IsStruct(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Struct
}

// IsMap reports weather val's kind is map or not
func IsMap(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Map
}

// toString converts any kind to string
func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
