package validation

import (
	"testing"
)

func Test_AddRule(t *testing.T) {
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
			name: "test add rule",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil && !tt.wantErr {
					t.Errorf("AddRule() error: failed to add rule, wantErr: %v, error: %v, args: %v", tt.wantErr, err, tt.args)
				}
			}()
			AddRule(tt.args.ruleName, tt.args.ruleFunc)
			if _, ok := rules["test"]; !ok {
				t.Errorf("AddRule() error: failed to add rule, wantErr: %v, error: %v, args: %v", tt.wantErr, nil, tt.args)
			}
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

func Test_stringRule(t *testing.T) {
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
			name: "test string rule",
			args: args{
				field:      "name",
				fieldValue: "Kyriakos",
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test string rule with non-string value",
			args: args{
				field:      "name",
				fieldValue: 44,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringFunc, stringExists := rules["string"]
			if !stringExists {
				panic("string rule is not exist")
			}
			err, validationError := stringFunc(tt.args.field, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("string rule: err: %v, wantErr: %v, validationError: %v, wantValidationError: %v, args: %v", err, validationError, tt.wantErr, tt.wantValidationErr, tt.args)
			}
			if (validationError != "") != tt.wantValidationErr {
				t.Errorf("string rule: err: %v, wantErr: %v, validationError: %v, wantValidationError: %v, args: %v", err, validationError, tt.wantErr, tt.wantValidationErr, tt.args)
			}
		})
	}
}

func Test_intRule(t *testing.T) {
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
			name: "test int rule",
			args: args{
				fieldName:  "intField",
				fieldValue: 31,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test int rule with non-int value",
			args: args{
				fieldName:  "intField",
				fieldValue: "s",
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := intRule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("intRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("intRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_int8Rule(t *testing.T) {
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
			name: "test int8 rule",
			args: args{
				fieldName:  "int8Field",
				fieldValue: int8(5),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test int8 rule with non-int8 value",
			args: args{
				fieldName:  "int8Field",
				fieldValue: int16(200),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := int8Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("intRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("intRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}
