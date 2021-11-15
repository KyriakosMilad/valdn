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

func createNewValidation(r Rules) *validation {
	v := validation{
		rules:       copyRules(r),
		errors:      make(Errors),
		fieldsExist: make(fieldsExist),
	}
	return &v
}

func Validate(name string, val interface{}, rlz []string) error {
	for _, r := range rlz {
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

func ValidateNested(val interface{}, r Rules) Errors {
	v := createNewValidation(r)

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

func ValidateJson(val string, r Rules) Errors {
	var jsonMap map[string]interface{}

	err := json.Unmarshal([]byte(val), &jsonMap)
	if err != nil {
		panic(err)
	}

	return ValidateNested(jsonMap, r)
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

func (v *validation) addValidationTagRules(t reflect.Type, parName string) {
	parName = makeParentNameJoinable(parName)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typ := f.Type
		name := parName + f.Name
		tRules := f.Tag.Get(TagName)

		_, ok := v.rules[name]
		if !ok && tRules != "" {
			var rlz []string
			for _, r := range strings.Split(tRules, RulesInStringSeparator) {
				rlz = append(rlz, r)
			}
			v.rules[name] = rlz
		}

		if typ.Kind() == reflect.Struct {
			v.addValidationTagRules(typ, name)
		}
	}
}

func (v *validation) validateStruct(val interface{}, name string) {
	rlz := v.getFieldRules(name)

	err := Validate(name, val, rlz)
	if err != nil {
		v.addError(name, err)
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
	}

	v.validateMapFields(convertInterfaceToMap(val), name)
}

func (v *validation) validateByType(name string, t reflect.Type, val interface{}) {
	v.registerField(name)
	r := v.getFieldRules(name)

	switch t.Kind() {
	case reflect.Struct:
		v.validateStruct(val, name)
	case reflect.Map:
		v.validateMap(val, name)
	default:
		err := Validate(name, val, r)
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
	for name, rlz := range v.rules {
		for _, r := range rlz {
			if r == "required" {
				_, ok := v.fieldsExist[name]
				if !ok {
					v.addError(name, errors.New(name+" is required"))
				}
			}
		}
	}
}
