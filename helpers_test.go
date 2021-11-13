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
