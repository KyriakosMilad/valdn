package validation

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_copyRules(t *testing.T) {
	type args struct {
		r Rules
	}
	tests := []struct {
		name               string
		args               args
		expectedRulesCount int
	}{
		{
			name: "test copy rules",
			args: args{
				r: Rules{"test": {"required", "kind:string"}},
			},
			expectedRulesCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := copyRules(tt.args.r); len(got) != tt.expectedRulesCount {
				t.Errorf("copyRules() = %v, expectedRulesCount %v", got, tt.expectedRulesCount)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type user struct {
		name string
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test to string with string value",
			args: args{val: "string"},
			want: "string",
		},
		{
			name: "test to string with uint value",
			args: args{val: 44},
			want: "44",
		},
		{
			name: "test to string with int value",
			args: args{val: -44},
			want: "-44",
		},
		{
			name: "test to string with float value",
			args: args{val: 44.44},
			want: "44.44",
		},
		{
			name: "test to string with complex number value",
			args: args{val: 44 + 22i},
			want: "(44+22i)",
		},
		{
			name: "test to string with bool value",
			args: args{val: true},
			want: "true",
		},
		{
			name: "test to string with map value",
			args: args{val: map[string]interface{}{"key": 55}},
			want: "map[key:55]",
		},
		{
			name: "test to string with struct value",
			args: args{val: user{name: "test"}},
			want: "{test}",
		},
		{
			name: "test to string with array value",
			args: args{val: [2]int{1, 2}},
			want: "[1 2]",
		},
		{
			name: "test to string with slice value",
			args: args{val: []int{1, 2}},
			want: "[1 2]",
		},
		{
			name: "test to string with Type value",
			args: args{val: reflect.TypeOf(map[string]interface{}{"key": 5})},
			want: "map[string]interface {}",
		},
		{
			name: "test to string with Kind value",
			args: args{val: reflect.String},
			want: "string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.want)
			if got := toString(tt.args.val); got != tt.want {
				t.Errorf("toString() = %v, expectedRulesCount %v", got, tt.want)
			}
		})
	}
}

func Test_splitRuleNameAndRuleValue(t *testing.T) {
	tests := []struct {
		name          string
		rule          string
		nameExpected  string
		valueExpected string
	}{
		{
			name:          "test get rule value from rule does have value",
			rule:          "val:test",
			nameExpected:  "val",
			valueExpected: "test",
		},
		{
			name:          "test get rule value from rule does not have value",
			rule:          "val",
			nameExpected:  "val",
			valueExpected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ruleName, ruleValue := splitRuleNameAndRuleValue(tt.rule); (ruleName != tt.nameExpected) || ruleValue != tt.valueExpected {
				t.Errorf("getRuleValue(): rName = %v, ruleValue = %v, ruleNameExpected = %v, ruleValueExpected = %v", ruleName, ruleValue, tt.nameExpected, tt.valueExpected)
			}
		})
	}
}

func Test_makeParentNameJoinable(t *testing.T) {
	tests := []struct {
		name    string
		parName string
		want    string
	}{
		{
			name:    "test make parent rName joinable",
			parName: "Parent",
			want:    "Parent.",
		},
		{
			name:    "test make parent rName joinable with . at the end",
			parName: "Parent.",
			want:    "Parent.",
		},
		{
			name:    "test make empty parent rName joinable",
			parName: "",
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeParentNameJoinable(tt.parName); got != tt.want {
				t.Errorf("makeParentNameJoinable() = %v, expectedRulesCount %v", got, tt.want)
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
		fNumber int
		pType   reflect.Type
		pValue  reflect.Value
		parName string
	}
	tests := []struct {
		name               string
		args               args
		fName              string
		fType              reflect.Type
		fValue             interface{}
		fieldValidationTag string
	}{
		{
			name: "test get struct field info",
			args: args{
				fNumber: 1,
				pType:   reflect.TypeOf(parentStruct),
				pValue:  reflect.ValueOf(parentStruct),
				parName: "",
			},
			fName:              "Age",
			fType:              reflect.TypeOf(parentStruct.Age),
			fValue:             1,
			fieldValidationTag: "required|int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldName, fieldType, fieldValue := getStructFieldInfo(tt.args.fNumber, tt.args.pType, tt.args.pValue, tt.args.parName)
			if fieldName != tt.fName {
				t.Errorf("getStructFieldInfo() got = %v, expectedRulesCount %v", fieldName, tt.fName)
			}
			if !reflect.DeepEqual(fieldType, tt.fType) {
				t.Errorf("getStructFieldInfo() fType = %v, expectedRulesCount %v", fieldType, tt.fType)
			}
			if !reflect.DeepEqual(fieldValue.Interface(), tt.fValue) {
				t.Errorf("getStructFieldInfo() fValue = %v, expectedRulesCount %v", fieldValue.Interface(), tt.fValue)
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
				t.Errorf("convertInterfaceToMap() type = %v, expectedRulesCount %v", reflect.TypeOf(got), reflect.TypeOf(tt.value))
			}
			if got := convertInterfaceToMap(tt.value); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.value)) {
				t.Errorf("convertInterfaceToMap() length = %v, lengthExpected %v", len(got), tt.lengthExpected)
			}
		})
	}
}
