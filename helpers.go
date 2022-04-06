package valdn

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func copyRules(r Rules) Rules {
	newMap := make(Rules, len(r))
	for k, v := range r {
		newMap[k] = v
	}
	return newMap
}

func toString(val interface{}) string {
	return fmt.Sprint(val)
}

func splitRuleNameAndRuleValue(rule string) (string, string) {
	if strings.ContainsRune(rule, ':') {
		ruleSpliced := strings.Split(rule, ":")
		return ruleSpliced[0], ruleSpliced[1]
	}
	return rule, ""
}

func makeParentNameJoinable(name string) string {
	if name != "" && name[len(name)-1] != '.' {
		return name + "."
	}
	return name
}

func getParentName(name string) string {
	nameSpliced := strings.Split(name, ".")
	if len(nameSpliced) > 1 {
		return strings.Join(nameSpliced[:len(nameSpliced)-1], ".")
	}
	return ""
}

func getStructFieldInfo(number int, parTyp reflect.Type, parVal reflect.Value, parName string) (string, reflect.Type, reflect.Value) {
	field := parTyp.Field(number)
	name := parName + field.Name
	typ := field.Type
	val := parVal.Field(number)

	return name, typ, val
}

func convertInterfaceToMap(value interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	val := reflect.ValueOf(value)
	for _, key := range val.MapKeys() {
		newMap[toString(key.Interface())] = val.MapIndex(key).Interface()
	}
	return newMap
}

func convertInterfaceToSlice(value interface{}) []interface{} {
	var s []interface{}
	v := reflect.ValueOf(value)
	for i := 0; i < v.Len(); i++ {
		s = append(s, v.Index(i).Interface())
	}
	return s
}

// stringToFloat converts val to float64.
// It returns error if val is not a float or an integer.
func interfaceToFloat(val interface{}) (float64, error) {
	var f64 float64
	if !IsInteger(val) && !IsFloat(val) {
		return f64, errors.New("val must be an integer or a float")
	}
	if v, ok := val.(float64); ok {
		f64 = v
	}
	if v, ok := val.(float32); ok {
		f64 = float64(v)
	}
	if v, ok := val.(int); ok {
		f64 = float64(v)
	}
	if v, ok := val.(uint); ok {
		f64 = float64(v)
	}
	if v, ok := val.(int8); ok {
		f64 = float64(v)
	}
	if v, ok := val.(uint8); ok {
		f64 = float64(v)
	}
	if v, ok := val.(int16); ok {
		f64 = float64(v)
	}
	if v, ok := val.(uint16); ok {
		f64 = float64(v)
	}
	if v, ok := val.(int32); ok {
		f64 = float64(v)
	}
	if v, ok := val.(uint32); ok {
		f64 = float64(v)
	}
	if v, ok := val.(int64); ok {
		f64 = float64(v)
	}
	if v, ok := val.(uint64); ok {
		f64 = float64(v)
	}
	return f64, nil
}

// stringToFloat converts s to float64.
// It returns error if s is not a float or an integer.
func stringToFloat(s string) (float64, error) {
	var f64 float64
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return f64, errors.New("string must contain an integer or a float")
		}
	} else {
		f64 = float64(i)
	}
	return f64, nil
}

// getLen gets v's length.
// It returns error if v is not array, slice, map, string, integer or float.
func getLen(v interface{}) (int, error) {
	switch {
	case IsCollection(v) || IsString(v):
		return reflect.ValueOf(v).Len(), nil
	case IsInteger(v) || IsFloat(v):
		l := 0
		stringVal := toString(v)
		if stringVal[0] == '-' {
			l -= 1
		}
		if IsFloat(v) {
			l -= 1
		}
		l += len(stringVal)
		return l, nil
	default:
		return 0, fmt.Errorf("can't get length of kind %v", reflect.TypeOf(v).Kind())
	}
}

func getFileSize(v interface{}) (int64, error) {
	if f, ok := v.(*os.File); ok {
		if (os.File{}) == *f {
			return 0, errors.New("can't get size from empty os.File")
		}
		s, err := f.Stat()
		if err != nil {
			return 0, err
		}
		return s.Size(), nil
	}
	if f, ok := v.(*multipart.FileHeader); ok {
		return f.Size, nil
	}
	return 0, fmt.Errorf("%v is not type of *os.File or *multipart.FileHeader", v)
}

func getFileExt(v interface{}) (string, error) {
	if f, ok := v.(*os.File); ok {
		if (os.File{}) == *f {
			return "", errors.New("can't get extension from empty os.File")
		}
		return filepath.Ext(f.Name()), nil
	}
	if f, ok := v.(*multipart.FileHeader); ok {
		return filepath.Ext(f.Filename), nil
	}
	return "", fmt.Errorf("%v is not type of *os.File or *multipart.FileHeader", v)
}
