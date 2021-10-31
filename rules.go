package validation

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

func init() {
	AddRule("required", requiredRule)
	AddRule("string", stringRule)
	AddRule("int", intRule)
	AddRule("int8", int8Rule)
	AddRule("int16", int16Rule)
	AddRule("int32", int32Rule)
	AddRule("int64", int64Rule)
	AddRule("uint", uintRule)
}
