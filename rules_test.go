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
