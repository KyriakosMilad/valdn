package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const defaultMaxMemory = 32 << 20 // 32 MB

func copyRules(r Rules) Rules {
	newMap := make(Rules)
	for k, v := range r {
		newMap[k] = v
	}
	return newMap
}

func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
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

func getStructFieldInfo(number int, parTyp reflect.Type, parVal reflect.Value, parName string) (string, reflect.Type, reflect.Value) {
	field := parTyp.Field(number)
	name := parName + field.Name
	typ := field.Type
	val := parVal.Field(number)

	return name, typ, val
}

func convertInterfaceToMap(value interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	if val, ok := value.(map[string]interface{}); ok {
		for k, v := range val {
			newMap[k] = v
		}
	}
	return newMap
}

// stringToFloat converts val to float64.
// It returns error if val is not a float or an integer.
func interfaceToFloat(val interface{}) (error, float64) {
	var f64 float64
	if !IsInteger(val) && !IsFloat(val) {
		return errors.New("val must be an integer or a float"), f64
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
	return nil, f64
}

// stringToFloat converts s's to float64.
// It returns error if s's value is not a float or an integer.
func stringToFloat(s string) (error, float64) {
	var f64 float64
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return errors.New("string must contain an integer or a float"), f64
		}
	} else {
		f64 = float64(i)
	}
	return nil, f64
}

// getLen gets v's length.
// It returns error if v is not array, slice, map, string, integer or float.
func getLen(v interface{}) (error, int) {
	switch {
	case IsCollection(v) || IsString(v):
		return nil, reflect.ValueOf(v).Len()
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
		return nil, l
	default:
		return fmt.Errorf("can't get length of kind %v", reflect.TypeOf(v).Kind()), 0
	}
}

func requestToMap(r *http.Request, rules Rules) map[string]interface{} {
	reqMap := make(map[string]interface{})

	// parse request by content type
	contentType := r.Header.Get("Content-Type")
	switch {
	case contentType == "application/json":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &reqMap)
		if err != nil {
			panic(err)
		}
	case strings.Contains(contentType, "multipart/form-data"):
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			panic(err)
		}
		for k := range rules {
			// check if field is a file
			if _, fhs, err := r.FormFile(k); err == nil {
				reqMap[k] = fhs
			} else {
				reqMap[k] = r.Form.Get(k)
			}
		}
	default:
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		for k := range rules {
			reqMap[k] = r.Form.Get(k)
		}
	}
	return reqMap
}
