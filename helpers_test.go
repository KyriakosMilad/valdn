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

func Test_getParentName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test get parent name",
			args: args{name: "parent.child.grandchild"},
			want: "parent.child",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getParentName(tt.args.name); got != tt.want {
				t.Errorf("getParentName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRuleExists(t *testing.T) {
	type args struct {
		rules []string
		rule  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test check if exist-rule exists",
			args: args{
				rules: []string{"rule2", "rule1"},
				rule:  "rule1",
			},
			want: true,
		},
		{
			name: "test check if non-exist-rule exists",
			args: args{
				rules: []string{"rule2", "rule1"},
				rule:  "rule1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRuleExists(tt.args.rules, tt.args.rule); got != tt.want {
				t.Errorf("isRuleExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isParentRequired(t *testing.T) {
	type args struct {
		fieldName       string
		validationRules map[string][]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test check if required-parent required",
			args: args{
				fieldName:       "parent.child",
				validationRules: map[string][]string{"parent": {"required"}},
			},
			want: true,
		},
		{
			name: "test check if non-required-parent required",
			args: args{
				fieldName:       "parent.child",
				validationRules: map[string][]string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isParentRequired(tt.args.fieldName, tt.args.validationRules); got != tt.want {
				t.Errorf("isParentRequired() = %v, want %v", got, tt.want)
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
			addValidationTagRules(tt.args.t, tt.args.validationRules, tt.args.parentName)
			if len(tt.args.validationRules) != tt.lengthExpected {
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

func Test_addValidationErrors(t *testing.T) {
	type args struct {
		validationErrors    map[string]string
		newValidationErrors map[string]string
	}
	tests := []struct {
		name           string
		args           args
		lengthExpected int
	}{
		{
			name: "test add validation errors",
			args: args{
				validationErrors:    map[string]string{"1": "one"},
				newValidationErrors: map[string]string{"2": "two"},
			},
			lengthExpected: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addValidationErrors(tt.args.validationErrors, tt.args.newValidationErrors)
			if len(tt.args.validationErrors) != tt.lengthExpected {
				t.Errorf("addValidationErrors() validationErrors = %v len = %v, want %v", tt.args.validationErrors, len(tt.args.validationErrors), tt.lengthExpected)
			}
		})
	}
}

func Test_getNestedRules(t *testing.T) {
	type args struct {
		validationRules map[string][]string
		structName      string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "test get nested rules",
			args: args{
				validationRules: map[string][]string{"Parent.child.name": {"required"}},
				structName:      "Parent",
			},
			want: map[string][]string{"Parent.child.name": {"required"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNestedRules(tt.args.validationRules, tt.args.structName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStructRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateNestedStruct(t *testing.T) {
	type Parent struct {
		Name string
	}
	type args struct {
		fieldName       string
		fieldValue      interface{}
		validationRules map[string][]string
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate nested struct",
			args: args{
				fieldName:       "Parent",
				fieldValue:      Parent{Name: "Beshoy"},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate nested struct with unsuitable data",
			args: args{
				fieldName:       "Parent",
				fieldValue:      Parent{},
				validationRules: map[string][]string{"Parent.Child": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := validateNestedStruct(tt.args.fieldName, tt.args.fieldValue, tt.args.validationRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNestedStruct() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateNestedStruct() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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

func Test_validateNestedMap(t *testing.T) {
	type args struct {
		fieldName       string
		fieldValue      interface{}
		validationRules map[string][]string
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate nested map",
			args: args{
				fieldName:       "Parent",
				fieldValue:      map[string]interface{}{"name": "Beshay"},
				validationRules: map[string][]string{"Parent.name": {"kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				fieldName:       "Parent",
				fieldValue:      map[string]interface{}{},
				validationRules: map[string][]string{"Parent.child": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := validateNestedMap(tt.args.fieldName, tt.args.fieldValue, tt.args.validationRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNestedMap() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateNestedMap() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
				fieldValue:      map[string]interface{}{},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
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
			err, validationErrors := validateByType(tt.args.fieldName, tt.args.fieldType, tt.args.fieldValue, tt.args.validationRules)
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
			err := validateStructFields(tt.args.t, tt.args.v, tt.args.parentName, tt.args.validationRules, tt.args.validationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStructFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, tt.args.validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(tt.args.validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateStructFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, tt.args.validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
			err := validateMapFields(tt.args.mapData, tt.args.parentName, tt.args.validationRules, tt.args.validationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMapFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, tt.args.validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(tt.args.validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateMapFields() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, tt.args.validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func Test_registerStructFields(t *testing.T) {
	type Parent struct {
		Name string
	}
	type args struct {
		structData   interface{}
		parentName   string
		fieldsExists map[string]bool
	}
	tests := []struct {
		name                      string
		args                      args
		expectedExistsFieldsCount int
	}{
		{
			name: "test register struct fields",
			args: args{
				structData:   Parent{Name: "Youlitta"},
				parentName:   "Parent",
				fieldsExists: map[string]bool{},
			},
			expectedExistsFieldsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerStructFields(tt.args.structData, tt.args.parentName, tt.args.fieldsExists)
			if len(tt.args.fieldsExists) != tt.expectedExistsFieldsCount {
				t.Errorf("registerStructFields(): fieldsExists %v, expectedValidationErrorsCount %v, args %v", tt.args.fieldsExists, tt.expectedExistsFieldsCount, tt.args)
			}
		})
	}
}

func Test_registerMapFields(t *testing.T) {
	type args struct {
		mapData      interface{}
		parentName   string
		fieldsExists map[string]bool
	}
	tests := []struct {
		name                      string
		args                      args
		expectedExistsFieldsCount int
	}{
		{
			name: "test register map fields",
			args: args{
				mapData:      map[string]interface{}{"age": 4},
				parentName:   "Parent",
				fieldsExists: map[string]bool{},
			},
			expectedExistsFieldsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerMapFields(tt.args.mapData, tt.args.parentName, tt.args.fieldsExists)
			if len(tt.args.fieldsExists) != tt.expectedExistsFieldsCount {
				t.Errorf("registerMapFields(): fieldsExists %v, expectedValidationErrorsCount %v, args %v", tt.args.fieldsExists, tt.expectedExistsFieldsCount, tt.args)
			}
		})
	}
}

func Test_registerNestedFieldsByType(t *testing.T) {
	type Parent struct {
		Name string
	}
	type args struct {
		fieldType    reflect.Type
		fieldValue   interface{}
		fieldName    string
		fieldsExists map[string]bool
	}
	tests := []struct {
		name                      string
		args                      args
		expectedExistsFieldsCount int
	}{
		{
			name: "test register nested fields by type - struct",
			args: args{
				fieldType:    reflect.TypeOf(Parent{Name: "Nefrtari"}),
				fieldValue:   Parent{Name: "Nefrtari"},
				fieldName:    "Parent",
				fieldsExists: map[string]bool{},
			},
			expectedExistsFieldsCount: 1,
		},
		{
			name: "test register nested fields by type - map",
			args: args{
				fieldType:    reflect.TypeOf(map[string]interface{}{"Age": 1}),
				fieldValue:   map[string]interface{}{"Age": 1},
				fieldName:    "Parent",
				fieldsExists: map[string]bool{},
			},
			expectedExistsFieldsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerNestedFieldsByType(tt.args.fieldType, tt.args.fieldValue, tt.args.fieldName, tt.args.fieldsExists)
			if len(tt.args.fieldsExists) != tt.expectedExistsFieldsCount {
				t.Errorf("registerNestedFieldsByType(): fieldsExists %v, expectedValidationErrorsCount %v, args %v", tt.args.fieldsExists, tt.expectedExistsFieldsCount, tt.args)
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
			validateNonExistRequiredFields(tt.args.validationRules, tt.args.fieldsExists, tt.args.validationErrors)
			if len(tt.args.validationErrors) != tt.validationErrorsCountExpected {
				t.Errorf("validateNonExistRequiredFields() validationErrors = %v len = %v, want %v", tt.args.validationErrors, len(tt.args.validationErrors), tt.validationErrorsCountExpected)
			}
		})
	}
}
