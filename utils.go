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
