package validation

import "testing"

func TestAddRule(t *testing.T) {
	type args struct {
		ruleName string
		ruleFunc func(field string, fieldValue interface{}, fieldExists bool, ruleValue string) (err error, validationError string)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test add rule",
			args: args{
				ruleName: "test",
				ruleFunc: func(field string, fieldValue interface{}, fieldExists bool, ruleValue string) (err error, validationError string) {
					return
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddRule(tt.args.ruleName, tt.args.ruleFunc)
		})
	}
}

func TestRequiredRule(t *testing.T) {
	type args struct {
		field       string
		fieldValue  interface{}
		fieldExists bool
		ruleValue   string
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
				field:       "name",
				fieldValue:  "Kyriakos",
				fieldExists: true,
				ruleValue:   "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test required rule with empty data",
			args: args{
				field:       "name",
				fieldValue:  "",
				fieldExists: true,
				ruleValue:   "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
		{
			name: "test required rule with non exist field",
			args: args{
				field:       "name",
				fieldValue:  "Kyriakos",
				fieldExists: false,
				ruleValue:   "",
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
			err, validationError := requiredFunc(tt.args.field, tt.args.fieldValue, tt.args.fieldExists, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("required rule: err: %v, wantErr: %v, validationError: %v, wantValidationError: %v, args: %v", err, validationError, tt.wantErr, tt.wantValidationErr, tt.args)
			}
			if (validationError != "") != tt.wantValidationErr {
				t.Errorf("required rule: err: %v, wantErr: %v, validationError: %v, wantValidationError: %v, args: %v", err, validationError, tt.wantErr, tt.wantValidationErr, tt.args)
			}
		})
	}
}
