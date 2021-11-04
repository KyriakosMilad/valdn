package validation

import (
	"testing"
)

func Test_CustomRule(t *testing.T) {
	type args struct {
		ruleName string
		ruleFunc func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test add custom rule",
			args: args{
				ruleName: "test",
				ruleFunc: func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string) {
					return
				},
			},
			wantErr: false,
		},
		{
			name: "test add rule already exists",
			args: args{
				ruleName: "test",
				ruleFunc: func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string) {
					return
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CustomRule(tt.args.ruleName, tt.args.ruleFunc)
		})
	}
}

func Test_requiredRule(t *testing.T) {
	type args struct {
		field      string
		fieldValue interface{}
		ruleValue  string
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		wantValidationErr bool
	}{
		{
			name: "test required rule",
			args: args{
				field:      "name",
				fieldValue: "Kyriakos",
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test required rule with zero value",
			args: args{
				field:      "name",
				fieldValue: "",
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requiredFunc, requiredExists := rules["required"]
			if !requiredExists {
				panic("required rule is not exist")
			}
			err, validationError := requiredFunc(tt.args.field, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("required rule: err: %v, wantErr: %v, validationError: %v, wantValidationError: %v, args: %v", err, tt.wantErr, validationError, tt.wantValidationErr, tt.args)
			}
			if (validationError != "") != tt.wantValidationErr {
				t.Errorf("required rule: err: %v, wantErr: %v, validationError: %v, wantValidationError: %v, args: %v", err, tt.wantErr, validationError, tt.wantValidationErr, tt.args)
			}
		})
	}
}

func Test_typeRule(t *testing.T) {
	type user struct {
		name string
	}
	type args struct {
		fieldName  string
		fieldValue interface{}
		ruleValue  string
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		wantValidationErr bool
	}{
		{
			name: "test type rule with string",
			args: args{
				fieldName:  "typeField",
				fieldValue: "string",
				ruleValue:  "string",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with uint",
			args: args{
				fieldName:  "typeField",
				fieldValue: uint(44),
				ruleValue:  "uint",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with int",
			args: args{
				fieldName:  "typeField",
				fieldValue: -44,
				ruleValue:  "int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with float",
			args: args{
				fieldName:  "typeField",
				fieldValue: 44.44,
				ruleValue:  "float64",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with complex number",
			args: args{
				fieldName:  "typeField",
				fieldValue: 44 + 22i,
				ruleValue:  "complex128",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with bool",
			args: args{
				fieldName:  "typeField",
				fieldValue: true,
				ruleValue:  "bool",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with map",
			args: args{
				fieldName:  "typeField",
				fieldValue: map[string]interface{}{"key": 55},
				ruleValue:  "map[string]interface {}",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with struct",
			args: args{
				fieldName:  "typeField",
				fieldValue: user{name: "test"},
				ruleValue:  "user",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with array",
			args: args{
				fieldName:  "typeField",
				fieldValue: [2]int{1, 2},
				ruleValue:  "[2]int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with slice",
			args: args{
				fieldName:  "typeField",
				fieldValue: []int{1, 2},
				ruleValue:  "[]int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := typeRule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("typeRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("typeRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_kindRule(t *testing.T) {
	type args struct {
		fieldName  string
		fieldValue interface{}
		ruleValue  string
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		wantValidationErr bool
	}{
		{
			name: "test type rule with string",
			args: args{
				fieldName:  "kindField",
				fieldValue: "string",
				ruleValue:  "string",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with uint",
			args: args{
				fieldName:  "kindField",
				fieldValue: uint(44),
				ruleValue:  "uint",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with int",
			args: args{
				fieldName:  "kindField",
				fieldValue: -44,
				ruleValue:  "int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with float",
			args: args{
				fieldName:  "kindField",
				fieldValue: 44.44,
				ruleValue:  "float64",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with complex number",
			args: args{
				fieldName:  "kindField",
				fieldValue: 44 + 22i,
				ruleValue:  "complex128",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with bool",
			args: args{
				fieldName:  "kindField",
				fieldValue: true,
				ruleValue:  "bool",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with map",
			args: args{
				fieldName:  "kindField",
				fieldValue: map[string]interface{}{"key": 55},
				ruleValue:  "map",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with struct",
			args: args{
				fieldName:  "kindField",
				fieldValue: struct{}{},
				ruleValue:  "struct",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with array",
			args: args{
				fieldName:  "kindField",
				fieldValue: [2]int{1, 2},
				ruleValue:  "array",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with slice",
			args: args{
				fieldName:  "kindField",
				fieldValue: []int{1, 2},
				ruleValue:  "slice",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := kindRule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("kindRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("kindRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}
