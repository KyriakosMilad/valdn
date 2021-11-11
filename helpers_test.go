package validation

import (
	"reflect"
	"testing"
)

func Test_splitRuleNameAndRuleValue(t *testing.T) {
	tests := []struct {
		name              string
		rule              string
		ruleNameExpected  string
		ruleValueExpected string
	}{
		{
			name:              "test get rule value from rule does have value",
			rule:              "val:test",
			ruleNameExpected:  "val",
			ruleValueExpected: "test",
		},
		{
			name:              "test get rule value from rule does not have value",
			rule:              "val",
			ruleNameExpected:  "val",
			ruleValueExpected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ruleName, ruleValue := splitRuleNameAndRuleValue(tt.rule); (ruleName != tt.ruleNameExpected) || ruleValue != tt.ruleValueExpected {
				t.Errorf("getRuleValue(): ruleName = %v, ruleValue = %v, ruleNameExpected = %v, ruleValueExpected = %v", ruleName, ruleValue, tt.ruleNameExpected, tt.ruleValueExpected)
			}
		})
	}
}

func Test_getRuleInfo(t *testing.T) {
	type args struct {
		rule string
	}
	tests := []struct {
		name      string
		args      args
		ruleName  string
		ruleValue string
		ruleFunc  RuleFunc
		ruleExist bool
	}{
		{
			name: "test get rule info",
			args: args{
				rule: "kind:string",
			},
			ruleName:  "kind",
			ruleValue: "string",
			ruleFunc:  rules["kind"],
			ruleExist: true,
		},
		{
			name: "test get info of rule does not exist",
			args: args{
				rule: "string",
			},
			ruleName:  "string",
			ruleValue: "",
			ruleFunc:  nil,
			ruleExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruleName, ruleValue, ruleFunc, ruleExist := getRuleInfo(tt.args.rule)
			if ruleName != tt.ruleName {
				t.Errorf("getRuleInfo() ruleName = %v, want %v", ruleName, tt.ruleName)
			}
			if ruleValue != tt.ruleValue {
				t.Errorf("getRuleInfo() ruleValue = %v, want %v", ruleValue, tt.ruleValue)
			}
			if !reflect.DeepEqual(toString(ruleFunc), toString(tt.ruleFunc)) {
				t.Errorf("getRuleInfo() ruleFunc = %v, want %v", toString(ruleFunc), toString(tt.ruleFunc))
			}
			if ruleExist != tt.ruleExist {
				t.Errorf("getRuleInfo() ruleExist = %v, want %v", ruleExist, tt.ruleExist)
			}
		})
	}
}

func Test_makeParentNameJoinable(t *testing.T) {
	tests := []struct {
		name       string
		parentName string
		want       string
	}{
		{
			name:       "test make parent name joinable",
			parentName: "Parent",
			want:       "Parent.",
		},
		{
			name:       "test make parent name joinable with . at the end",
			parentName: "Parent.",
			want:       "Parent.",
		},
		{
			name:       "test make empty parent name joinable",
			parentName: "",
			want:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeParentNameJoinable(tt.parentName); got != tt.want {
				t.Errorf("makeParentNameJoinable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addValidationTagRules(t *testing.T) {
	type Parent struct {
		Name string `validation:"required|string"`
		Age  int    `validation:"required|string"`
	}
	type args struct {
		t               reflect.Type
		validationRules map[string][]string
		parentName      string
	}
	tests := []struct {
		name           string
		args           args
		lengthExpected int
	}{
		{
			name: "test add validation tag rules",
			args: args{
				t:               reflect.TypeOf(Parent{}),
				validationRules: map[string][]string{"test": {"required"}},
				parentName:      "",
			},
			lengthExpected: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialize(tt.args.validationRules)
			addValidationTagRules(tt.args.t, tt.args.parentName)
			if len(validationRules) != tt.lengthExpected {
				t.Errorf("addValidationTagRules() validationRules = %v len = %v, want %v", tt.args.validationRules, len(tt.args.validationRules), tt.lengthExpected)
			}
		})
	}
}

func Test_getStructFieldInfo(t *testing.T) {
	type Parent struct {
		Name string `validation:"required|string"`
		Age  int    `validation:"required|int"`
	}
	parentStruct := Parent{
		Name: "test_struct",
		Age:  1,
	}
	type args struct {
		fieldNumber int
		parentType  reflect.Type
		parentValue reflect.Value
		parentName  string
	}
	tests := []struct {
		name               string
		args               args
		fieldName          string
		fieldType          reflect.Type
		fieldValue         interface{}
		fieldValidationTag string
	}{
		{
			name: "test get struct field info",
			args: args{
				fieldNumber: 1,
				parentType:  reflect.TypeOf(parentStruct),
				parentValue: reflect.ValueOf(parentStruct),
				parentName:  "",
			},
			fieldName:          "Age",
			fieldType:          reflect.TypeOf(parentStruct.Age),
			fieldValue:         1,
			fieldValidationTag: "required|int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldName, fieldType, fieldValue := getStructFieldInfo(tt.args.fieldNumber, tt.args.parentType, tt.args.parentValue, tt.args.parentName)
			if fieldName != tt.fieldName {
				t.Errorf("getStructFieldInfo() got = %v, want %v", fieldName, tt.fieldName)
			}
			if !reflect.DeepEqual(fieldType, tt.fieldType) {
				t.Errorf("getStructFieldInfo() fieldType = %v, want %v", fieldType, tt.fieldType)
			}
			if !reflect.DeepEqual(fieldValue.Interface(), tt.fieldValue) {
				t.Errorf("getStructFieldInfo() fieldValue = %v, want %v", fieldValue.Interface(), tt.fieldValue)
			}
		})
	}
}

func Test_convertInterfaceToMap(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		lengthExpected int
	}{
		{
			name:           "test convert interface to map",
			value:          map[string]interface{}{"test": 1},
			lengthExpected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertInterfaceToMap(tt.value); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.value)) {
				t.Errorf("convertInterfaceToMap() type = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.value))
			}
			if got := convertInterfaceToMap(tt.value); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.value)) {
				t.Errorf("convertInterfaceToMap() length = %v, lengthExpected %v", len(got), tt.lengthExpected)
			}
		})
	}
}

func Test_validateByType(t *testing.T) {
	type Parent struct {
		Name string
	}
	type args struct {
		fieldName       string
		fieldType       reflect.Type
		fieldValue      interface{}
		validationRules map[string][]string
		fieldsExists    map[string]bool
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate by type - struct",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(Parent{}),
				fieldValue:      Parent{Name: "Barssom"},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
				fieldsExists:    map[string]bool{"Parent.Name": true},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate by type - struct with unsuitable data",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(Parent{}),
				fieldValue:      Parent{},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
				fieldsExists:    map[string]bool{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate by type - map",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(map[string]interface{}{}),
				fieldValue:      map[string]interface{}{"Name": "Pola"},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
				fieldsExists:    map[string]bool{"Parent.Name": true},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate by type - map with unsuitable data",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(map[string]interface{}{}),
				fieldValue:      map[string]interface{}{"Name": 22},
				validationRules: map[string][]string{"Parent.Name": {"kind:string"}},
				fieldsExists:    map[string]bool{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate by type - field",
			args: args{
				fieldName:       "age",
				fieldType:       reflect.TypeOf(1),
				fieldValue:      1,
				validationRules: map[string][]string{"age": {"required"}},
				fieldsExists:    map[string]bool{"age": true},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate by type - field with unsuitable data",
			args: args{
				fieldName:       "age",
				fieldType:       reflect.TypeOf(1),
				fieldValue:      0,
				validationRules: map[string][]string{"age": {"required"}},
				fieldsExists:    map[string]bool{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialize(tt.args.validationRules)
			err := validateByType(tt.args.fieldName, tt.args.fieldType, tt.args.fieldValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateByType() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateByType() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func Test_validateStructFields(t *testing.T) {
	type Parent struct {
		Name string `validation:"required|string"`
	}
	type args struct {
		t                reflect.Type
		v                reflect.Value
		parentName       string
		validationRules  map[string][]string
		validationErrors map[string]string
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate struct fields",
			args: args{
				t:                reflect.TypeOf(Parent{Name: "Sherry"}),
				v:                reflect.ValueOf(Parent{Name: "Sherry"}),
				parentName:       "",
				validationRules:  map[string][]string{"Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate struct fields with unsuitable data",
			args: args{
				t:                reflect.TypeOf(Parent{}),
				v:                reflect.ValueOf(Parent{}),
				parentName:       "",
				validationRules:  map[string][]string{"Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialize(tt.args.validationRules)
			err := validateStructFields(tt.args.t, tt.args.v, tt.args.parentName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStructFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateStructFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func Test_validateMapFields(t *testing.T) {
	type args struct {
		mapData          map[string]interface{}
		parentName       string
		validationRules  map[string][]string
		validationErrors map[string]string
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate map fields",
			args: args{
				mapData:          map[string]interface{}{"Name": "Pola"},
				parentName:       "",
				validationRules:  map[string][]string{"Parent.Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				mapData:          map[string]interface{}{"Name": ""},
				parentName:       "Parent",
				validationRules:  map[string][]string{"Parent.Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialize(tt.args.validationRules)
			err := validateMapFields(tt.args.mapData, tt.args.parentName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMapFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateMapFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func Test_validateNonExistRequiredFields(t *testing.T) {
	type args struct {
		validationRules  map[string][]string
		fieldsExists     map[string]bool
		validationErrors map[string]string
	}
	tests := []struct {
		name                          string
		args                          args
		validationErrorsCountExpected int
	}{
		{
			name: "test validate non exist required fields",
			args: args{
				validationRules:  map[string][]string{"field1": {"required", "string"}, "field3": {"string", "required"}},
				fieldsExists:     map[string]bool{"field1": true, "field2": true},
				validationErrors: map[string]string{},
			},
			validationErrorsCountExpected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialize(tt.args.validationRules)
			fieldsExists = tt.args.fieldsExists
			validateNonExistRequiredFields()
			if len(validationErrors) != tt.validationErrorsCountExpected {
				t.Errorf("validateNonExistRequiredFields() validationErrors = %v len = %v, want %v", validationErrors, len(validationErrors), tt.validationErrorsCountExpected)
			}
		})
	}
}
