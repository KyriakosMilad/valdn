package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	Rules       map[string][]string
	Errors      map[string]string
	fieldsExist map[string]bool
)

var (
	TagName                = "validation"
	RulesInStringSeparator = "|"
)

type validation struct {
	rules       Rules
	errors      Errors
	fieldsExist fieldsExist
}

// createNewValidation copies rules and initialise new validation with it.
// rules are copied in case they will be manipulated later it doesn't affect the original rules.
func createNewValidation(rules Rules) *validation {
	v := validation{
		rules:       copyRules(rules),
		errors:      make(Errors),
		fieldsExist: make(fieldsExist),
	}
	return &v
}

// Validate validates val with rules.
// If error found it will not check the rest of the rules and return the error.
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

		err := rFunc(name, val, rVal)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateNested validates val and it's nested fields with rules and returns Errors.
// If error found it will not check the rest of field rules and move to the next field.
// If struct or map has error it's nested fields will not be validated.
// It panics if val's kind is not map or struct.
// It panics if one of the rules is not registered.
// It panics if one of the fields is a map and it's type is not map[string]interface{}.
func ValidateNested(val interface{}, rules Rules) Errors {
	v := createNewValidation(rules)

	t := reflect.TypeOf(val)
	switch t.Kind() {
	case reflect.Struct:
		v.addValidationTagRules(t, "")
		v.validateStruct(val, "")
	case reflect.Map:
		v.validateMap(convertInterfaceToMap(val), "")
	default:
		panic("ValidateNested() can only validate struct and map ")
	}

	v.validateNonExistRequiredFields()

	return v.errors
}

// ValidateJson transforms json string to a map and validates it with rules and returns Errors.
// If error found it will not check the rest of field rules and move to the next field.
// If map has error it's nested fields will not be validated.
// It panics if val is not json.
// It panics if one of the rules is not registered.
func ValidateJson(val string, rules Rules) Errors {
	var jsonMap map[string]interface{}

	err := json.Unmarshal([]byte(val), &jsonMap)
	if err != nil {
		panic(err)
	}

	return ValidateNested(jsonMap, rules)
}

func (v *validation) registerField(name string) {
	v.fieldsExist[name] = true
}

func (v *validation) addError(name string, err error) {
	v.errors[name] = err.Error()
}

func (v *validation) getFieldRules(name string) []string {
	return v.rules[name]
}

// addValidationTagRules gets rules from struct tag for every field and adds them to field rules if field has no rules.
func (v *validation) addValidationTagRules(t reflect.Type, parName string) {
	parName = makeParentNameJoinable(parName)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typ := f.Type
		name := parName + f.Name
		tRules := f.Tag.Get(TagName)

		// add tag rules only if field has no rules
		_, ok := v.rules[name]
		if !ok && tRules != "" {
			var rules []string
			for _, r := range strings.Split(tRules, RulesInStringSeparator) {
				rules = append(rules, r)
			}
			v.rules[name] = rules
		}

		if typ.Kind() == reflect.Struct {
			v.addValidationTagRules(typ, name)
		}
	}
}

func (v *validation) validateStruct(val interface{}, name string) {
	rules := v.getFieldRules(name)

	err := Validate(name, val, rules)
	if err != nil {
		v.addError(name, err)
		return
	}

	typ := reflect.TypeOf(val)
	value := reflect.ValueOf(val)
	v.validateStructFields(typ, value, name)
}

func (v *validation) validateMap(val interface{}, name string) {
	if !reflect.DeepEqual(reflect.TypeOf(val), reflect.TypeOf(map[string]interface{}{})) {
		panic(fmt.Errorf("error validating %v: type %v is not supported", name, reflect.TypeOf(val)))
	}

	r := v.getFieldRules(name)
	err := Validate(name, val, r)
	if err != nil {
		v.addError(name, err)
		return
	}

	v.validateMapFields(convertInterfaceToMap(val), name)
}

func (v *validation) validateByType(name string, t reflect.Type, val interface{}) {
	v.registerField(name)
	rules := v.getFieldRules(name)

	switch t.Kind() {
	case reflect.Struct:
		v.validateStruct(val, name)
	case reflect.Map:
		v.validateMap(val, name)
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

func (v *validation) validateNonExistRequiredFields() {
	for name, rules := range v.rules {
		for _, r := range rules {
			if r == "required" {
				_, ok := v.fieldsExist[name]
				if !ok {
					v.addError(name, errors.New(name+" is required"))
				}
			}
		}
	}
}
