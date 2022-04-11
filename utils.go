package valdn

import (
	"encoding/json"
	"net"
	"reflect"
	"regexp"
)

const (
	emailRegex = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
	macRegex   = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
	urlRegex   = "[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
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

// IsKind reports weather val's kind equals kind.
func IsKind(val interface{}, kind string) bool {
	if k := reflect.TypeOf(val).Kind(); toString(k) == kind {
		return true
	}
	return false
}

// IsKindIn reports weather val's kind is one of kinds.
func IsKindIn(val interface{}, kinds []string) bool {
	kind := toString(reflect.TypeOf(val).Kind())
	for _, k := range kinds {
		if k == kind {
			return true
		}
	}
	return false
}

// IsType reports weather val's type equals typ.
func IsType(val interface{}, typ string) bool {
	var typeInString string
	if t := reflect.TypeOf(val); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	if typeInString == typ {
		return true
	}
	return false
}

// IsTypeIn reports weather val's type is one of types.
func IsTypeIn(val interface{}, types []string) bool {
	var typeInString string
	if t := reflect.TypeOf(val); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	for _, t := range types {
		if t == typeInString {
			return true
		}
	}
	return false
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vString := toString(val)
		if vString[0] != '-' {
			return true
		}
		return false
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

// IsCollection reports weather val's kins is one of (Array, Slice, Map, Struct) or not.
func IsCollection(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return true
	default:
		return false
	}
}

// IsEmail reports weather s is a valid email address or not.
func IsEmail(s string) bool {
	r := regexp.MustCompile(emailRegex)
	return r.MatchString(s)
}

// IsJSON reports weather s is a valid json or not.
func IsJSON(s string) bool {
	var decodedJSON map[string]interface{}
	err := json.Unmarshal([]byte(s), &decodedJSON)
	return err == nil
}

// IsIPv4 reports weather s is a valid IPv4 or not.
func IsIPv4(s string) bool {
	if net.ParseIP(s) == nil {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return true
		}
	}
	return false
}

// IsIPv6 reports weather s is a valid IPv6 or not.
func IsIPv6(s string) bool {
	if net.ParseIP(s) == nil {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return true
		}
	}
	return false
}

// IsIP reports weather s is a valid IP or not.
func IsIP(s string) bool {
	return net.ParseIP(s) != nil
}

// IsMAC reports weather s is a valid MAC address or not.
func IsMAC(s string) bool {
	r := regexp.MustCompile(macRegex)
	return r.MatchString(s)
}

// IsURL reports weather s is a valid URL or not.
func IsURL(s string) bool {
	r := regexp.MustCompile(urlRegex)
	return r.MatchString(s)
}

// IsFile reports weather v is a valid file or not.
func IsFile(v interface{}) bool {
	_, err := getFileSize(v)
	if err != nil {
		return false
	}
	return true
}
