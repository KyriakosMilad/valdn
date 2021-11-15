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

func Validate(name string, val interface{}, rlz []string) (error, string) {
	for _, r := range rlz {
		if r == "" {
			continue
		}

		rName, rVal, rFunc, rExist := getRuleInfo(r)
		if !rExist {
			err := errors.New("unknown validation rule: " + rName)
			return err, ""
		}
		err, vErr := rFunc(name, val, rVal)
		if err != nil {
			return err, ""
		}
		if vErr != "" {
			return nil, vErr
		}
	}
	return nil, ""
}

func ValidateNested(val interface{}, r Rules) (error, Errors) {
	var err error
	v := createNewValidation(r)

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

func ValidateJson(val string, r Rules) (error, Errors) {
	var decodedJson map[string]interface{}

	err := json.Unmarshal([]byte(val), &decodedJson)
	if err != nil {
		return errors.New(toString(err)), nil
	}

	return ValidateNested(decodedJson, r)
}

func (v *validation) registerField(name string) {
	v.fieldsExist[name] = true
}

func (v *validation) addError(name string, err string) {
	v.errors[name] = err
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
		tRules := f.Tag.Get("validation")

		_, ok := v.rules[name]
		if !ok && tRules != "" {
			var rlz []string
			for _, r := range strings.Split(tRules, "|") {
				rlz = append(rlz, r)
			}
			v.rules[name] = rlz
		}

		if typ.Kind() == reflect.Struct {
			v.addValidationTagRules(typ, name)
		}
	}
}

func (v *validation) validateStruct(val interface{}, name string) error {
	rlz := v.getFieldRules(name)
	err, structErr := Validate(name, val, rlz)
	if err != nil {
		return err
	}
	if structErr != "" {
		v.addError(name, structErr)
	}

	typ := reflect.TypeOf(val)
	value := reflect.ValueOf(val)
	return v.validateStructFields(typ, value, name)
}

func (v *validation) validateMap(val interface{}, name string) error {
	if !reflect.DeepEqual(reflect.TypeOf(val), reflect.TypeOf(map[string]interface{}{})) {
		return fmt.Errorf("error validating %v: type %v is not supported", name, reflect.TypeOf(val))
	}

	r := v.getFieldRules(name)
	err, mapErr := Validate(name, val, r)
	if err != nil {
		return err
	}
	if mapErr != "" {
		v.addError(name, mapErr)
	}

	return v.validateMapFields(convertInterfaceToMap(val), name)
}

func (v *validation) validateByType(name string, t reflect.Type, val interface{}) error {
	var err error
	v.registerField(name)
	r := v.getFieldRules(name)

	switch t.Kind() {
	case reflect.Struct:
		err = v.validateStruct(val, name)
	case reflect.Map:
		err = v.validateMap(val, name)
	default:
		var vErr string
		err, vErr = Validate(name, val, r)
		if vErr != "" {
			v.addError(name, vErr)
		}
	}
	return err
}

func (v *validation) validateStructFields(parTyp reflect.Type, parVal reflect.Value, parName string) error {
	parName = makeParentNameJoinable(parName)
	for i := 0; i < parTyp.NumField(); i++ {
		name, typ, value := getStructFieldInfo(i, parTyp, parVal, parName)
		err := v.validateByType(name, typ, value.Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validation) validateMapFields(val map[string]interface{}, parName string) error {
	parName = makeParentNameJoinable(parName)
	for name, value := range val {
		name = parName + name
		typ := reflect.TypeOf(value)
		err := v.validateByType(name, typ, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validation) validateNonExistRequiredFields() {
	for name, rlz := range v.rules {
		for _, r := range rlz {
			if r == "required" {
				_, ok := v.fieldsExist[name]
				if !ok {
					v.addError(name, name+" is required")
				}
			}
		}
	}
}
