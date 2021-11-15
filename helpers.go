package validation

import (
	"reflect"
	"strings"
)

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

func copyRules(r Rules) Rules {
	newMap := make(Rules)
	for k, v := range r {
		newMap[k] = v
	}
	return newMap
}
