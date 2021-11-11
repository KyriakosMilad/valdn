package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type (
	Rules  map[string][]string
	Errors map[string]string
	fields map[string]bool
)

var (
	validationErrors Errors
	validationRules  Rules
	fieldsExists     fields
)

func initialize(rules Rules) {
	validationRules = copyRules(rules)
	validationErrors = make(Errors)
	fieldsExists = make(fields)
}

func Validate(val interface{}, rules Rules) (error, Errors) {
	var err error
	initialize(rules)

	t := reflect.TypeOf(val)
	switch t.Kind() {
	case reflect.Struct:
		addValidationTagRules(t, "")
		err = validateStruct(val, "")
	case reflect.Map:
		err = validateMap(convertInterfaceToMap(val), "")
	default:
		err = errors.New("Validate() can only validate struct and map ")
	}
	if err != nil {
		return err, nil
	}

	validateNonExistRequiredFields()

	return err, validationErrors
}

func ValidateField(fieldName string, fieldValue interface{}, fieldRules []string) (error, string) {
	for _, rule := range fieldRules {
		if rule == "" {
			continue
		}

		ruleName, ruleValue, ruleFunc, ruleExists := getRuleInfo(rule)
		if !ruleExists {
			err := errors.New("unknown validation rule: " + ruleName)
			return err, ""
		}
		err, validationErr := ruleFunc(fieldName, fieldValue, ruleValue)
		if err != nil {
			return err, ""
		}
		if validationErr != "" {
			return nil, validationErr
		}
	}
	return nil, ""
}

func ValidateJson(jsonData string, validationRules Rules) (error, Errors) {
	initialize(validationRules)
	var decodedJson map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return errors.New(toString(err)), nil
	}

	err = validateMap(decodedJson, "")
	if err != nil {
		return err, nil
	}

	validateNonExistRequiredFields()
	return nil, validationErrors
}

func validateMap(mapData interface{}, name string) error {
	if !reflect.DeepEqual(reflect.TypeOf(mapData), reflect.TypeOf(map[string]interface{}{})) {
		return fmt.Errorf("error validating %v: type %v is not supported", name, reflect.TypeOf(mapData))
	}

	fieldsExists[name] = true
	fieldRules := validationRules[name]

	err, mapErr := ValidateField(name, mapData, fieldRules)
	if err != nil {
		return err
	}
	if mapErr != "" {
		validationErrors[name] = mapErr
	}

	return validateMapFields(convertInterfaceToMap(mapData), name)
}

func validateStruct(structData interface{}, name string) error {
	fieldsExists[name] = true
	fieldRules := validationRules[name]

	err, structErr := ValidateField(name, structData, fieldRules)
	if err != nil {
		return err
	}
	if structErr != "" {
		validationErrors[name] = structErr
	}

	t := reflect.TypeOf(structData)
	v := reflect.ValueOf(structData)
	return validateStructFields(t, v, name)
}
