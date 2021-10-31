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

func init() {
	AddRule("required", requiredRule)
	AddRule("string", stringRule)
}
