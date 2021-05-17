package validation

var rules = make(map[string]func(field string, fieldValue interface{}, fieldExists bool, ruleValue string) (err error, validationError string))

func AddRule(ruleName string, ruleFunc func(field string, fieldValue interface{}, fieldExists bool, ruleValue string) (err error, validationError string)) {
	_, ruleExists := rules[ruleName]
	if ruleExists {
		panic("rule already registered")
	}

	rules[ruleName] = ruleFunc

	return
}
