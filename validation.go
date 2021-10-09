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

func ValidateStruct(structData interface{}, validationRules map[string][]string, parentName string) (error, map[string]string) {
	var err error
	validationErrors := make(map[string]string)
	fieldsExist := make(map[string]bool)

	baseT := reflect.TypeOf(structData)
	baseV := reflect.ValueOf(structData)
	if baseT.Kind() != reflect.Struct {
		err = errors.New("can only proceed `struct` kind")
		return err, nil
	}

	var validateNestedStruct func(t reflect.Type, v reflect.Value, parentName string) (error, map[string]string)
	validateNestedStruct = func(t reflect.Type, v reflect.Value, parentName string) (error, map[string]string) {
		var fieldValidationErrors map[string]string

		if parentName != "" {
			parentName += "."
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValidationTag := field.Tag.Get("validation")
			fieldName := parentName + field.Name
			fieldType := field.Type
			fieldValue := v.Field(i)
			if !fieldValue.IsZero() {
				fieldsExist[fieldName] = true
			}

			fieldRules, fieldHasRules := validationRules[fieldName]
			if !fieldHasRules && fieldValidationTag == "" && fieldType.Kind() != reflect.Struct {
				continue
			}

			if !fieldHasRules {
				fieldRules = strings.Split(fieldValidationTag, "|")
				fieldRequired := false
				for _, v := range fieldRules {
					if v == "required" {
						fieldRequired = true
						break
					}
				}
				if fieldRequired && parentName != "" {
					parentName = strings.Split(parentName, ".")[0]
					parentRequired := false
					for j := 0; j < baseT.NumField(); j++ {
						f := baseT.Field(j)
						if parentName == f.Name {
							parentRules, parentHasRules := validationRules[parentName]
							if !parentHasRules {
								fTag := f.Tag.Get("validation")
								if fTag != "" {
									parentRules = strings.Split(fTag, "|")
								}
							}
						ParentRulesLoop:
							for _, v := range parentRules {
								if v == "required" {
									parentRequired = true
									break ParentRulesLoop
								}
							}
							if !parentRequired {
							FieldRulesLoop:
								for i, v := range fieldRules {
									if v == "required" {
										fieldRules = append(fieldRules[:i], fieldRules[i+1:]...)
										break FieldRulesLoop
									}
								}
							}
						}
					}
				}
			}

			switch {
			case fieldType.Kind() == reflect.Struct:
				err, structFieldErrors := ValidateField(fieldName, fieldValue.Interface(), fieldRules)
				if err != nil {
					return err, nil
				}
				if len(structFieldErrors) > 0 {
					return nil, structFieldErrors
				}
				for i := 0; i < fieldValue.NumField(); i++ {
					fieldsExist[fieldName+"."+reflect.TypeOf(fieldValue).Field(i).Name] = true
				}
				err, fieldValidationErrors = validateNestedStruct(fieldType, fieldValue, fieldName)
			case fieldType == reflect.TypeOf(map[string]interface{}{}):
				err, mapFieldErrors := ValidateField(fieldName, fieldValue.Interface(), fieldRules)
				if err != nil {
					return err, nil
				}
				if len(mapFieldErrors) > 0 {
					return nil, mapFieldErrors
				}
				mapRules := make(map[string][]string)
				for k, v := range validationRules {
					if strings.Contains(k, fieldName+".") {
						mapRules[k] = v
					}
				}
				if len(mapRules) == 0 {
					continue
				}
				mapFields := make(map[string]interface{})
				if val, ok := fieldValue.Interface().(map[string]interface{}); ok {
					for k, v := range val {
						mapFields[k] = v
						fieldsExist[fieldName+"."+k] = true
					}
				}
				err, fieldValidationErrors = ValidateMap(mapFields, mapRules, fieldName)
			default:
				if fieldType.Kind() == reflect.Map {
					return fmt.Errorf("error validating %v: type %v is not supported", fieldName, fieldType), nil
				}
				err, fieldValidationErrors = ValidateField(fieldName, fieldValue.Interface(), fieldRules)
			}

			if len(fieldValidationErrors) > 0 {
				for k, v := range fieldValidationErrors {
					validationErrors[k] = v
				}
			}

			if err != nil {
				return err, nil
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

	return validateNestedStruct(baseT, baseV, parentName)
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

		switch {
		case reflect.TypeOf(fieldValue).Kind() == reflect.Struct:
			var structFieldErrors map[string]string
			err, structFieldErrors = ValidateField(fieldName, fieldValue, fieldRules)
			if err != nil {
				return err, nil
			}
			if len(structFieldErrors) > 0 {
				return nil, structFieldErrors
			}
			structRules := make(map[string][]string)
			for k, v := range validationRules {
				if strings.Contains(k, fullName+".") {
					structRules[k] = v
				}
			}
			for i := 0; i < reflect.TypeOf(fieldValue).NumField(); i++ {
				fieldsExist[fullName+"."+reflect.TypeOf(fieldValue).Field(i).Name] = true
			}
			err, fieldValidationErrors = ValidateStruct(fieldValue, structRules, fullName)
		case reflect.TypeOf(fieldValue) == reflect.TypeOf(map[string]interface{}{}):
			var mapFieldErrors map[string]string
			err, mapFieldErrors = ValidateField(fieldName, fieldValue, fieldRules)
			if err != nil {
				return err, nil
			}
			if len(mapFieldErrors) > 0 {
				return err, mapFieldErrors
			}
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
