package validation

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func ValidateField(fieldName string, fieldValue interface{}, fieldRules []string) (err error, validationErrors map[string]string) {
	validationErrors = make(map[string]string)

	for _, rule := range fieldRules {
		if rule == "" {
			continue
		}
		ruleFunc, ruleExists := rules[rule]
		if !ruleExists {
			err = errors.New("unknown validation rule: " + rule)
			return err, validationErrors
		}

		var ruleValue string
		if strings.ContainsRune(rule, ':') {
			ruleValue = strings.Split(rule, ":")[1]
		}

		err, validationError := ruleFunc(fieldName, fieldValue, ruleValue)
		if err != nil {
			return err, validationErrors
		}

		if validationError != "" {
			validationErrors[fieldName] = validationError
			break
		}

	}

	return
}

func ValidateStruct(structData interface{}, validationRules map[string][]string) (err error, validationErrors map[string]string) {
	t := reflect.TypeOf(structData)
	v := reflect.ValueOf(structData)
	if t.Kind() != reflect.Struct {
		err = errors.New("can only proceed `struct` kind")
		return err, validationErrors
	}

	var validateNestedStruct func(t reflect.Type, v reflect.Value, parentName string) (err error, validationErrors map[string]string)
	validateNestedStruct = func(t reflect.Type, v reflect.Value, parentName string) (err error, validationErrors map[string]string) {
		validationErrors = make(map[string]string)

		if parentName != "" {
			parentName += "."
		}

	FIELDS:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValidationTag := field.Tag.Get("validation")
			fieldName := parentName + field.Name
			fieldType := field.Type
			fieldValue := v.Field(i)

			if fieldType.Kind() == reflect.Struct {
				err, nestedStructValidationErrors := validateNestedStruct(fieldType, fieldValue, fieldName)
				for key, val := range nestedStructValidationErrors {
					validationErrors[key] = val
				}

				if err != nil {
					return err, validationErrors
				}
				continue
			}

			if _, ok := validationRules[fieldName]; !ok && fieldValidationTag == "" {
				continue
			}

			fieldExists := true
			if fieldValue.IsZero() {
				fieldExists = false
			}

			var fieldRules []string
			if validationRules[fieldName] != nil {
				fieldRules = validationRules[fieldName]
			} else {
				fieldRules = strings.Split(fieldValidationTag, "|")
			}

			for _, rule := range fieldRules {
				if rule == "" {
					continue
				}
				var ruleValue string
				if strings.ContainsRune(rule, ':') {
					ruleDetailed := strings.Split(rule, ":")
					ruleValue = ruleDetailed[1]
					rule = ruleDetailed[0]
				}

				ruleFunc, ruleExists := rules[rule]
				if !ruleExists {
					err = errors.New("unknown validation rule: " + rule)
					return err, validationErrors
				}

				err, validationError := ruleFunc(fieldName, fieldValue.Interface(), fieldExists, ruleValue)
				if err != nil {
					return err, validationErrors
				}

				if validationError != "" {
					validationErrors[fieldName] = validationError
					continue FIELDS
				}
			}

		}
		return err, validationErrors
	}

	return validateNestedStruct(t, v, "")
}

func ValidateJson(jsonData string, validationRules map[string][]string) (err error, validationErrors map[string]string) {
	var decodedJson map[string]interface{}

	err = json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return err, nil
	}

	return ValidateMap(decodedJson, validationRules)
}

func ValidateMap(mapData map[string]interface{}, validationRules map[string][]string) (err error, validationErrors map[string]string) {
	validationErrors = make(map[string]string)

FIELDS:
	for field, fieldRules := range validationRules {
		fieldVal, fieldExists := mapData[field]

		for _, rule := range fieldRules {
			var ruleVal string
			if strings.ContainsRune(rule, ':') {
				ruleDetailed := strings.Split(rule, ":")
				ruleVal = ruleDetailed[1]
				rule = ruleDetailed[0]
			}

			ruleFunc, ruleExists := rules[rule]
			if !ruleExists {
				err = errors.New("unknown validation rule: " + rule)
				return err, validationErrors
			}

			err, validationError := ruleFunc(field, fieldVal, fieldExists, ruleVal)
			if err != nil {
				return err, validationErrors
			}

			if validationError != "" {
				validationErrors[field] = validationError
				continue FIELDS
			}
		}
	}

	return
}
