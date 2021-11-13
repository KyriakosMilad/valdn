package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	Rules        map[string][]string
	Errors       map[string]string
	fieldsExists map[string]bool
)

type validation struct {
	rules        Rules
	errors       Errors
	fieldsExists fieldsExists
}

func createNewValidation(rules Rules) *validation {
	v := validation{
		rules:        copyRules(rules),
		errors:       make(Errors),
		fieldsExists: make(fieldsExists),
	}
	return &v
}

func Validate(fieldName string, fieldValue interface{}, fieldRules []string) (error, string) {
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

func ValidateNested(val interface{}, rules Rules) (error, Errors) {
	var err error
	v := createNewValidation(rules)

	t := reflect.TypeOf(val)
	switch t.Kind() {
	case reflect.Struct:
		v.addValidationTagRules(t, "")
		err = v.validateStruct(val, "")
	case reflect.Map:
		err = v.validateMap(convertInterfaceToMap(val), "")
	default:
		err = errors.New("ValidateNested() can only validate struct and map ")
	}
	if err != nil {
		return err, nil
	}

	v.validateNonExistRequiredFields()

	return nil, v.errors
}

func ValidateJson(jsonData string, validationRules Rules) (error, Errors) {
	var decodedJson map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return errors.New(toString(err)), nil
	}

	return ValidateNested(decodedJson, validationRules)
}

func (v *validation) registerField(fieldName string) {
	v.fieldsExists[fieldName] = true
}

func (v *validation) addError(fieldName string, err string) {
	v.errors[fieldName] = err
}

func (v *validation) getFieldRules(fieldName string) []string {
	return v.rules[fieldName]
}

func (v *validation) addValidationTagRules(t reflect.Type, parentName string) {
	parentName = makeParentNameJoinable(parentName)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldName := parentName + field.Name
		fieldValidationTag := field.Tag.Get("validation")

		_, fieldHasRules := v.rules[fieldName]
		if !fieldHasRules && fieldValidationTag != "" {
			var fieldRules []string
			for _, v := range strings.Split(fieldValidationTag, "|") {
				fieldRules = append(fieldRules, v)
			}
			v.rules[fieldName] = fieldRules
		}

		if fieldType.Kind() == reflect.Struct {
			v.addValidationTagRules(fieldType, fieldName)
		}
	}
}

func (v *validation) validateStruct(structData interface{}, name string) error {
	fieldRules := v.getFieldRules(name)
	err, structErr := Validate(name, structData, fieldRules)
	if err != nil {
		return err
	}
	if structErr != "" {
		v.errors[name] = structErr
	}

	typ := reflect.TypeOf(structData)
	val := reflect.ValueOf(structData)
	return v.validateStructFields(typ, val, name)
}

func (v *validation) validateMap(mapData interface{}, name string) error {
	if !reflect.DeepEqual(reflect.TypeOf(mapData), reflect.TypeOf(map[string]interface{}{})) {
		return fmt.Errorf("error validating %v: type %v is not supported", name, reflect.TypeOf(mapData))
	}

	fieldRules := v.getFieldRules(name)
	err, mapErr := Validate(name, mapData, fieldRules)
	if err != nil {
		return err
	}
	if mapErr != "" {
		v.errors[name] = mapErr
	}

	return v.validateMapFields(convertInterfaceToMap(mapData), name)
}

func (v *validation) validateByType(fieldName string, fieldType reflect.Type, fieldValue interface{}) error {
	var err error
	v.registerField(fieldName)
	fieldRules := v.getFieldRules(fieldName)

	switch fieldType.Kind() {
	case reflect.Struct:
		err = v.validateStruct(fieldValue, fieldName)
	case reflect.Map:
		err = v.validateMap(fieldValue, fieldName)
	default:
		var fieldValidationError string
		err, fieldValidationError = Validate(fieldName, fieldValue, fieldRules)
		if fieldValidationError != "" {
			v.errors[fieldName] = fieldValidationError
		}
	}
	return err
}

func (v *validation) validateStructFields(typ reflect.Type, val reflect.Value, parentName string) error {
	parentName = makeParentNameJoinable(parentName)
	for i := 0; i < typ.NumField(); i++ {
		fieldName, fieldType, fieldValue := getStructFieldInfo(i, typ, val, parentName)
		err := v.validateByType(fieldName, fieldType, fieldValue.Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validation) validateMapFields(mapData map[string]interface{}, parentName string) error {
	parentName = makeParentNameJoinable(parentName)
	for fieldName, fieldValue := range mapData {
		fieldName = parentName + fieldName
		fieldType := reflect.TypeOf(fieldValue)
		err := v.validateByType(fieldName, fieldType, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validation) validateNonExistRequiredFields() {
	for key, rules := range v.rules {
		for _, val := range rules {
			if val == "required" {
				_, ok := v.fieldsExists[key]
				if !ok {
					v.errors[key] = key + " is required"
				}
			}
		}
	}
}
