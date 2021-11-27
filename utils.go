package validation

import (
	"encoding/json"
	"reflect"
	"regexp"
)

// IsEmpty reports weather val is empty or not.
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

// IsString reports weather val's kind is string or not.
func IsString(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.String
}

// IsInt reports weather val's kind is int or not.
func IsInt(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int
}

// IsInt8 reports weather val's kind is int8 or not.
func IsInt8(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int8
}

// IsInt16 reports weather val's kind is int16 or not.
func IsInt16(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int16
}

// IsInt32 reports weather val's kind is int32 or not.
func IsInt32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int32
}

// IsInt64 reports weather val's kind is int64 or not.
func IsInt64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Int64
}

// IsUint reports weather val's kind is uint or not.
func IsUint(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint
}

// IsUint8 reports weather val's kind is uint8 or not.
func IsUint8(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint8
}

// IsUint16 reports weather val's kind is uint16 or not.
func IsUint16(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint16
}

// IsUint32 reports weather val's kind is uint32 or not.
func IsUint32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint32
}

// IsUint64 reports weather val's kind is uint64 or not.
func IsUint64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Uint64
}

// IsFloat32 reports weather val's kind is float32 or not.
func IsFloat32(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Float32
}

// IsFloat64 reports weather val's kind is float64 or not.
func IsFloat64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Float64
}

// IsComplex64 reports weather val's kind is complex64 or not.
func IsComplex64(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Complex64
}

// IsComplex128 reports weather val's kind is complex128 or not.
func IsComplex128(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Complex128
}

// IsBool reports weather val's kind is bool or not.
func IsBool(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Bool
}

// IsSlice reports weather val's kind is slice or not.
func IsSlice(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Slice
}

// IsArray reports weather val's kind is array or not.
func IsArray(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Array
}

// IsStruct reports weather val's kind is struct or not.
func IsStruct(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Struct
}

// IsMap reports weather val's kind is map or not.
func IsMap(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Map
}

// IsInteger reports weather val is integer or not.
func IsInteger(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Uint, reflect.Int, reflect.Uint8, reflect.Int8, reflect.Uint16, reflect.Int16, reflect.Uint32, reflect.Int32, reflect.Uint64, reflect.Int64:
		return true
	default:
		return false
	}
}

// IsUnsignedInteger reports weather val is unsigned integer or not.
func IsUnsignedInteger(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Int:
		if val.(int) > 0 {
			return true
		}
	case reflect.Int8:
		if val.(int8) > 0 {
			return true
		}
	case reflect.Int16:
		if val.(int16) > 0 {
			return true
		}
	case reflect.Int32:
		if val.(int32) > 0 {
			return true
		}
	case reflect.Int64:
		if val.(int64) > 0 {
			return true
		}
	}
	return false
}

// IsFloat reports weather val is float or not.
func IsFloat(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsUnsignedFloat reports weather val is unsigned float or not.
func IsUnsignedFloat(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Float32, reflect.Float64:
		vString := toString(val)
		if vString[0] != '-' {
			return true
		}
		return false
	default:
		return false
	}
}

// IsComplex reports weather val is complex number or not.
func IsComplex(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

// IsNumeric reports weather val is numeric or not.
func IsNumeric(val interface{}) bool {
	if !IsInteger(val) && !IsFloat(val) && !IsComplex(val) {
		return false
	}
	return true
}

// IsCollection reports weather val's kins is one of (Array, Slice, Map) or not.
func IsCollection(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return true
	default:
		return false
	}
}

// IsEmail reports weather s is a valid email address or not.
func IsEmail(s string) bool {
	r, _ := regexp.Compile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")
	match := r.MatchString(s)
	if match {
		return true
	}
	return false
}

// IsJSON reports weather s is a valid json or not.
func IsJSON(s string) bool {
	var decodedJson map[string]interface{}
	err := json.Unmarshal([]byte(s), &decodedJson)
	if err == nil {
		return true
	}
	return false
}
