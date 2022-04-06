package valdn

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type (
	Rules       map[string][]string
	Errors      map[string]string
	fieldsExist map[string]bool
)

var (
	TagName      = "valdn"
	TagSeparator = "|"
)

type validation struct {
	rules       Rules
	errors      Errors
	fieldsExist fieldsExist
}

// createNewValidation copies rules and initialise new validation with it.
// rules are copied in case they will be manipulated later it doesn't affect the original rules.
func createNewValidation(rules Rules) *validation {
	return &validation{
		rules:       copyRules(rules),
		errors:      make(Errors),
		fieldsExist: make(fieldsExist),
	}
}

// Validate validates single value by rules.
// If an error is found it will not check the rest of the rules and return the error.
// It panics if one of the rules is not registered.
func Validate(name string, val interface{}, rules []string) error {
	for _, r := range rules {
		if r == "" {
			continue
		}

		rName, rVal, rFunc, rExist := getRuleInfo(r)
		if !rExist {
			panic("unknown rule: " + rName)
		}

		if err := rFunc(name, val, rVal); err != nil {
			return err
		}
	}
	return nil
}

// ValidateStruct validates struct, and it's nested fields by rules and returns Errors.
// It panics if val is not kind of struct.
// Unexported fields will be ignored.
// If an error is found it will not check the rest of the field's rules and continue to the next field.
// If a parent has error it's nested fields will not be validated.
// It panics if one of the rules is not registered.
func ValidateStruct(val interface{}, rules Rules) Errors {
	if !IsStruct(val) {
		panic("val is not a struct")
	}
	v := createNewValidation(rules)
	v.addTagRules(val, "")

	v.validateStruct(val, "")
	v.validateNonExistRequiredFields()

	return v.errors
}

// ValidateMap validates map, and it's nested fields by rules and returns Errors.
// It panics if val is not kind of map.
// If an error is found it will not check the rest of the field's rules and continue to the next field.
// If a parent has error it's nested fields will not be validated.
// It panics if one of the rules is not registered.
func ValidateMap(val interface{}, rules Rules) Errors {
	if !IsMap(val) {
		panic(fmt.Errorf("ValidateMap: %v is not kind of map", val))
	}
	v := createNewValidation(rules)
	v.addTagRules(val, "")

	v.validateMap(val, "")
	v.validateNonExistRequiredFields()

	return v.errors
}

// ValidateSlice validates slice, and it's nested fields by rules and returns Errors.
// If an error is found it will not check the rest of the field's rules and continue to the next field.
// If a parent has error it's nested fields will not be validated.
// It panics if one of the rules is not registered.
func ValidateSlice(val interface{}, rules Rules) Errors {
	if !IsSlice(val) {
		panic(fmt.Errorf("ValidateSlice: %v is not kind of slice", val))
	}
	v := createNewValidation(rules)
	v.addTagRules(val, "")

	v.validateSlice(val, "")
	v.validateNonExistRequiredFields()

	return v.errors
}

// ValidateJSON transforms JSON string to a map and validates it by rules and returns Errors.
// If an error is found it will not check the rest of the field's rules and continue to the next field.
// If parent has error it's nested fields will not be validated.
// It panics if val is not JSON.
// It panics if one of the rules is not registered.
func ValidateJSON(val string, rules Rules) Errors {
	var jsonMap map[string]interface{}

	if err := json.Unmarshal([]byte(val), &jsonMap); err != nil {
		panic(err)
	}
	return ValidateMap(jsonMap, rules)
}

// ValidateRequest validates request by rules and returns Errors.
// It validates request of content type: multipart/form-data, application/json and application/x-www-form-urlencoded.
// It validates url parameters.
// It panics if body is not compatible with header content type.
// It panics if one of the rules is not registered.
// If an error is found it will not check the rest of the field's rules and continue to the next field.
func ValidateRequest(r *http.Request, rules Rules) Errors {
	m := parseRequest(r, rules)
	return ValidateMap(m, rules)
}

func (v *validation) registerField(name string) {
	v.fieldsExist[name] = true
}

func (v *validation) addError(name string, err error) {
	v.errors[name] = err.Error()
}

func (v *validation) getFieldRules(name string) []string {
	val, ok := v.rules[name]
	if ok {
		return val
	}
	return v.rules[getParentName(name)+".*"]
}

func (v *validation) getParentRules(name string) []string {
	val, ok := v.rules[name]
	if ok {
		return val
	}
	if name != "" {
		return v.rules[getParentName(name)+".*"]
	}
	return []string{}
}

// addTagRules gets rules from struct tag for every field and adds them to field rules if field has no rules.
func (v *validation) addTagRules(val interface{}, parName string) {
	parName = makeParentNameJoinable(parName)

	if IsMap(val) {
		for _, key := range reflect.ValueOf(val).MapKeys() {
			value := reflect.ValueOf(val).MapIndex(key).Interface()
			switch {
			case IsStruct(value), IsMap(value), IsSlice(value):
				v.addTagRules(value, parName+toString(key))
			}
		}
	}

	if IsSlice(val) {
		for i := 0; i < reflect.ValueOf(val).Len(); i++ {
			value := reflect.ValueOf(val).Index(i).Interface()
			switch {
			case IsStruct(value), IsMap(value), IsSlice(value):
				v.addTagRules(value, parName+toString(i))
			}
		}
	}

	if IsStruct(val) {
		t := reflect.TypeOf(val)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			typ := f.Type
			name := parName + f.Name
			tRules := f.Tag.Get(TagName)

			// add tag rules only if field has no rules
			_, ok := v.rules[name]
			if !ok && tRules != "" {
				v.rules[name] = strings.Split(tRules, TagSeparator)
			}

			switch typ.Kind() {
			case reflect.Struct, reflect.Map, reflect.Slice:
				v.addTagRules(f, name)
			}
		}
	}
}

func (v *validation) validateStruct(val interface{}, name string) {
	r := v.getParentRules(name)

	if err := Validate(name, val, r); err != nil {
		v.addError(name, err)
		return
	}

	typ := reflect.TypeOf(val)
	value := reflect.ValueOf(val)
	v.validateStructFields(typ, value, name)
}

func (v *validation) validateMap(val interface{}, name string) {
	r := v.getParentRules(name)
	if err := Validate(name, val, r); err != nil {
		v.addError(name, err)
		return
	}

	v.validateMapFields(convertInterfaceToMap(val), name)
}

func (v *validation) validateSlice(val interface{}, name string) {
	r := v.getParentRules(name)
	if err := Validate(name, val, r); err != nil {
		v.addError(name, err)
		return
	}

	v.validateSliceFields(convertInterfaceToSlice(val), name)
}

func (v *validation) validateByType(name string, t reflect.Type, val interface{}) {
	v.registerField(name)
	rules := v.getFieldRules(name)

	switch t.Kind() {
	case reflect.Struct:
		v.validateStruct(val, name)
	case reflect.Map:
		v.validateMap(val, name)
	case reflect.Slice:
		v.validateSlice(val, name)
	default:
		err := Validate(name, val, rules)
		if err != nil {
			v.addError(name, err)
		}
	}
}

func (v *validation) validateStructFields(parTyp reflect.Type, parVal reflect.Value, parName string) {
	parName = makeParentNameJoinable(parName)
	for i := 0; i < parTyp.NumField(); i++ {
		name, typ, val := getStructFieldInfo(i, parTyp, parVal, parName)
		// ignore unexported field
		if !val.CanInterface() {
			continue
		}
		v.validateByType(name, typ, val.Interface())
	}
}

func (v *validation) validateMapFields(val map[string]interface{}, parName string) {
	parName = makeParentNameJoinable(parName)
	for name, value := range val {
		name = parName + name
		typ := reflect.TypeOf(value)
		v.validateByType(name, typ, value)
	}
}

func (v *validation) validateSliceFields(val []interface{}, parName string) {
	parName = makeParentNameJoinable(parName)
	for idx, value := range val {
		name := parName + toString(idx)
		typ := reflect.TypeOf(value)
		v.validateByType(name, typ, value)
	}
}

func (v *validation) validateNonExistRequiredFields() {
	for name, rules := range v.rules {
		if name == "*" {
			continue
		}
		for _, r := range rules {
			rName, rVal := splitRuleNameAndRuleValue(r)
			if rName == "required" {
				_, ok := v.fieldsExist[name]
				if !ok {
					v.addError(name, errors.New(GetErrMsg("required", rVal, name, "")))
				}
			}
		}
	}
}
