package validation

import (
	"reflect"
	"testing"
)

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
				t.Errorf("getStructFieldInfo() got = %v, want %v", fieldName, tt.fName)
			}
			if !reflect.DeepEqual(fieldType, tt.fType) {
				t.Errorf("getStructFieldInfo() fType = %v, want %v", fieldType, tt.fType)
			}
			if !reflect.DeepEqual(fieldValue.Interface(), tt.fValue) {
				t.Errorf("getStructFieldInfo() fValue = %v, want %v", fieldValue.Interface(), tt.fValue)
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
