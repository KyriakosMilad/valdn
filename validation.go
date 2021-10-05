package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func ValidateField(fieldName string, fieldValue interface{}, fieldRules []string) (error, map[string]string) {
	validationErrors := make(map[string]string)

	for _, rule := range fieldRules {
		if rule == "" {
			continue
		}

		var ruleValue string
		if strings.ContainsRune(rule, ':') {
			ruleValue = strings.Split(rule, ":")[1]
		}

		ruleFunc, ruleExist := rules[rule]
		if !ruleExist {
			err := errors.New("unknown validation rule: " + rule)
			return err, nil
		}

		err, validationError := ruleFunc(fieldName, fieldValue, ruleValue)
		if err != nil {
			return err, nil
		}

		if validationError != "" {
			validationErrors[fieldName] = validationError
			break
		}

	}

	return nil, validationErrors
}

func ValidateStruct(structData interface{}, validationRules map[string][]string) (err error, validationErrors map[string]string) {
	t := reflect.TypeOf(structData)
	v := reflect.ValueOf(structData)
	if t.Kind() != reflect.Struct {
		err = errors.New("can only proceed `struct` kind")
		return err, nil
	}

	var validateNestedStruct func(t reflect.Type, v reflect.Value, parentName string) (err error, validationErrors map[string]string)
	validateNestedStruct = func(t reflect.Type, v reflect.Value, parentName string) (err error, validationErrors map[string]string) {
		validationErrors = make(map[string]string)

		if parentName != "" {
			parentName += "."
		}

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
					return err, nil
				}
				continue
			}

			if _, ok := validationRules[fieldName]; !ok && fieldValidationTag == "" {
				continue
			}

			var fieldRules []string
			if validationRules[fieldName] != nil {
				fieldRules = validationRules[fieldName]
			} else {
				fieldRules = strings.Split(fieldValidationTag, "|")
			}

			err, fieldValidationErrors := ValidateField(fieldName, fieldValue.Interface(), fieldRules)

			if len(fieldValidationErrors) > 0 {
				for k, v := range fieldValidationErrors {
					validationErrors[k] = v
				}
			}

			if err != nil {
				return err, nil
			}
		}
		return err, validationErrors
	}

	return validateNestedStruct(t, v, "")
}

func ValidateJson(jsonData string, validationRules map[string][]string) (error, map[string]string) {
	var decodedJson map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return err, nil
	}

	return ValidateMap(decodedJson, validationRules, "")
}

func ValidateMap(mapData map[string]interface{}, validationRules map[string][]string, parentName string) (error, map[string]string) {
	validationErrors := make(map[string]string)
	var fieldValidationErrors map[string]string
	var err error
	fieldsExist := make(map[string]bool)

	for fieldName, fieldValue := range mapData {
		var fullName string
		if parentName == "" {
			fullName = fieldName
		} else {
			fullName = parentName + "." + fieldName
		}

		fieldsExist[fullName] = true

		fieldRules, fieldHasRules := validationRules[fullName]

		if !fieldHasRules {
			continue
		}

		switch fieldValue.(type) {
		// uncomment after refactor ValidateStruct()
		//case struct{}:
		//	structRules := make(map[string][]string)
		//	for k, v := range validationRules {
		//		if strings.Contains(k, fullName + ".") {
		//			structRules[k] = v
		//		}
		//	}
		//	err, fieldValidationErrors = ValidateStruct(fieldValue, structRules)
		case map[string]interface{}:
			mapRules := make(map[string][]string)
			for k, v := range validationRules {
				if strings.Contains(k, fullName+".") {
					mapRules[k] = v
				}
			}
			mapFields := make(map[string]interface{})
			if val, ok := fieldValue.(map[string]interface{}); ok {
				for k, v := range val {
					mapFields[k] = v
					fieldsExist[fullName+"."+k] = true
				}
			}
			err, fieldValidationErrors = ValidateMap(mapFields, mapRules, fullName)
		default:
			if reflect.TypeOf(fieldValue).Kind() == reflect.Map {
				return fmt.Errorf("error validating %v: type %v is not supported", fullName, reflect.TypeOf(fieldValue)), nil
			}
			err, fieldValidationErrors = ValidateField(fieldName, fieldValue, fieldRules)
		}

		if err != nil {
			return err, nil
		}

		if len(fieldValidationErrors) > 0 {
			for k, v := range fieldValidationErrors {
				validationErrors[k] = v
			}
		}
	}

	for k, v := range validationRules {
		for _, val := range v {
			if val == "required" {
				_, ok := fieldsExist[k]
				if !ok {
					validationErrors[k] = k + " is required"
				}
			}
		}
	}

	return err, validationErrors
}
