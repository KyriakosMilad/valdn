package validation

import "reflect"

var rules = make(map[string]func(fieldName string, fieldValue interface{}, ruleValue string) (error, string))

func AddRule(ruleName string, ruleFunc func(fieldName string, fieldValue interface{}, ruleValue string) (error, string)) {
	_, ruleExists := rules[ruleName]
	if ruleExists {
		panic("rule already registered")
	}
	rules[ruleName] = ruleFunc
}

func requiredRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if IsZero(fieldValue) {
		validationError := fieldName + " is required"
		return nil, validationError
	}
	return nil, ""
}

func stringRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsString(fieldValue) {
		validationError := fieldName + " must be a string"
		return nil, validationError
	}
	return nil, ""
}

func intRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsInt(fieldValue) {
		validationError := fieldName + " must be an integer"
		return nil, validationError
	}
	return nil, ""
}

func int8Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsInt8(fieldValue) {
		validationError := fieldName + " must be type of int8"
		return nil, validationError
	}
	return nil, ""
}

func int16Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsInt16(fieldValue) {
		validationError := fieldName + " must be type of int16"
		return nil, validationError
	}
	return nil, ""
}

func int32Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsInt32(fieldValue) {
		validationError := fieldName + " must be type of int32"
		return nil, validationError
	}
	return nil, ""
}

func int64Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsInt64(fieldValue) {
		validationError := fieldName + " must be type of int64"
		return nil, validationError
	}
	return nil, ""
}

func uintRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsUint(fieldValue) {
		validationError := fieldName + " must be type of uint"
		return nil, validationError
	}
	return nil, ""
}

func uint8Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsUint8(fieldValue) {
		validationError := fieldName + " must be type of uint8"
		return nil, validationError
	}
	return nil, ""
}

func uint16Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsUint16(fieldValue) {
		validationError := fieldName + " must be type of uint16"
		return nil, validationError
	}
	return nil, ""
}

func uint32Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsUint32(fieldValue) {
		validationError := fieldName + " must be type of uint32"
		return nil, validationError
	}
	return nil, ""
}

func uint64Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsUint64(fieldValue) {
		validationError := fieldName + " must be type of uint64"
		return nil, validationError
	}
	return nil, ""
}

func float32Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsFloat32(fieldValue) {
		validationError := fieldName + " must be type of float32"
		return nil, validationError
	}
	return nil, ""
}

func float64Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsFloat64(fieldValue) {
		validationError := fieldName + " must be type of float64"
		return nil, validationError
	}
	return nil, ""
}

func complex64Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsComplex64(fieldValue) {
		validationError := fieldName + " must be type of complex64"
		return nil, validationError
	}
	return nil, ""
}

func complex128Rule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsComplex128(fieldValue) {
		validationError := fieldName + " must be type of complex128"
		return nil, validationError
	}
	return nil, ""
}

func boolRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsBool(fieldValue) {
		validationError := fieldName + " must be type of bool"
		return nil, validationError
	}
	return nil, ""
}

func sliceRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsSlice(fieldValue) {
		validationError := fieldName + " must be type of slice"
		return nil, validationError
	}
	return nil, ""
}

func arrayRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsArray(fieldValue) {
		validationError := fieldName + " must be type of array"
		return nil, validationError
	}
	return nil, ""
}

func structRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsStruct(fieldValue) {
		validationError := fieldName + " must be type of struct"
		return nil, validationError
	}
	return nil, ""
}

func mapRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if !IsMap(fieldValue) {
		validationError := fieldName + " must be type of map"
		return nil, validationError
	}
	return nil, ""
}

func typeRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	var typeInString string
	if t := reflect.TypeOf(fieldValue); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	if typeInString != ruleValue {
		validationError := fieldName + " must be type of " + ruleValue
		return nil, validationError
	}
	return nil, ""
}

func init() {
	AddRule("required", requiredRule)
	AddRule("string", stringRule)
	AddRule("int", intRule)
	AddRule("int8", int8Rule)
	AddRule("int16", int16Rule)
	AddRule("int32", int32Rule)
	AddRule("int64", int64Rule)
	AddRule("uint", uintRule)
	AddRule("uint8", uint8Rule)
	AddRule("uint16", uint16Rule)
	AddRule("uint32", uint32Rule)
	AddRule("uint64", uint64Rule)
	AddRule("float32", float32Rule)
	AddRule("float64", float64Rule)
	AddRule("complex64", complex64Rule)
	AddRule("complex128", complex128Rule)
	AddRule("bool", boolRule)
	AddRule("slice", sliceRule)
	AddRule("array", arrayRule)
	AddRule("struct", structRule)
	AddRule("map", mapRule)
	AddRule("type", typeRule)
}
