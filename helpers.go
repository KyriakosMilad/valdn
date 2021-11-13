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

func getRuleInfo(rule string) (string, string, RuleFunc, bool) {
	ruleName, ruleValue := splitRuleNameAndRuleValue(rule)
	ruleFunc, ruleExists := rules[ruleName]
	return ruleName, ruleValue, ruleFunc, ruleExists
}

func makeParentNameJoinable(parentName string) string {
	if parentName != "" && parentName[len(parentName)-1] != '.' {
		return parentName + "."
	}
	return parentName
}

func getStructFieldInfo(fieldNumber int, parentType reflect.Type, parentValue reflect.Value, parentName string) (string, reflect.Type, reflect.Value) {
	field := parentType.Field(fieldNumber)
	fieldName := parentName + field.Name
	fieldType := field.Type
	fieldValue := parentValue.Field(fieldNumber)

	return fieldName, fieldType, fieldValue
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

func copyRules(rules Rules) Rules {
	newMap := make(Rules)
	for k, v := range rules {
		newMap[k] = v
	}
	return newMap
}
