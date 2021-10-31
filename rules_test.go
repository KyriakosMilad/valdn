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
				t.Errorf("int8Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("int8Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_int16Rule(t *testing.T) {
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
			name: "test int16 rule",
			args: args{
				fieldName:  "int16Field",
				fieldValue: int16(200),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test int16 rule with non-int16 value",
			args: args{
				fieldName:  "int16Field",
				fieldValue: int32(2147483646),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := int16Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("int16Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("int16Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_int32Rule(t *testing.T) {
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
			name: "test int32 rule",
			args: args{
				fieldName:  "int32Field",
				fieldValue: int32(2147483646),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test int32 rule with non-int32 value",
			args: args{
				fieldName:  "int32Field",
				fieldValue: int64(9223372036854775806),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := int32Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("int32Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("int32Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_int64Rule(t *testing.T) {
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
			name: "test int64 rule",
			args: args{
				fieldName:  "int64Field",
				fieldValue: int64(9223372036854775806),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test int64 rule with non-int64 value",
			args: args{
				fieldName:  "int64Field",
				fieldValue: int8(2),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := int64Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("int64Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("int64Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_uintRule(t *testing.T) {
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
			name: "test uint rule",
			args: args{
				fieldName:  "uintField",
				fieldValue: uint(15),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test uint rule with signed int value",
			args: args{
				fieldName:  "uintField",
				fieldValue: -15,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := uintRule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("uintRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("uintRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_uint8Rule(t *testing.T) {
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
			name: "test uint8 rule",
			args: args{
				fieldName:  "uint8Field",
				fieldValue: uint8(200),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test uint8 rule with non-uint8 value",
			args: args{
				fieldName:  "uint8Field",
				fieldValue: uint16(65534),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := uint8Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("uint8Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("uint8Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_uint16Rule(t *testing.T) {
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
			name: "test uint16 rule",
			args: args{
				fieldName:  "uint16Field",
				fieldValue: uint16(65534),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test uint16 rule with non-uint16 value",
			args: args{
				fieldName:  "uint16Field",
				fieldValue: uint32(4294967294),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := uint16Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("uint16Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("uint16Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_uint32Rule(t *testing.T) {
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
			name: "test uint32 rule",
			args: args{
				fieldName:  "uint32Field",
				fieldValue: uint32(4294967294),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test uint32 rule with non-uint32 value",
			args: args{
				fieldName:  "uint32Field",
				fieldValue: uint64(18446744073709551614),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := uint32Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("uint32Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("uint32Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_uint64Rule(t *testing.T) {
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
			name: "test uint64 rule",
			args: args{
				fieldName:  "uint64Field",
				fieldValue: uint64(18446744073709551614),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test uint64 rule with non-uint64 value",
			args: args{
				fieldName:  "uint64Field",
				fieldValue: int64(9223372036854775806),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := uint64Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("uint64Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("uint64Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_float32Rule(t *testing.T) {
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
			name: "test float32 rule",
			args: args{
				fieldName:  "float32Field",
				fieldValue: float32(2.2),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test float32 rule with non-float32 value",
			args: args{
				fieldName:  "float32Field",
				fieldValue: 55,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := float32Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("float32Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("float32Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_float64Rule(t *testing.T) {
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
			name: "test float64 rule",
			args: args{
				fieldName:  "float64Field",
				fieldValue: 2.2,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test float64 rule with non-float64 value",
			args: args{
				fieldName:  "float64Field",
				fieldValue: 55,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := float64Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("float64Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("float64Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_complex64Rule(t *testing.T) {
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
			name: "test complex64 rule",
			args: args{
				fieldName:  "complex64Field",
				fieldValue: complex64(2 + 2i),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test complex64 rule with non-complex64 value",
			args: args{
				fieldName:  "complex64Field",
				fieldValue: 456456 + 456456i,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := complex64Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("complex64Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("complex64Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_complex128Rule(t *testing.T) {
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
			name: "test complex128 rule",
			args: args{
				fieldName:  "complex128Field",
				fieldValue: 2 + 2i,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test complex128 rule with non-complex128 value",
			args: args{
				fieldName:  "complex128Field",
				fieldValue: complex64(456456 + 456456i),
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := complex128Rule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("complex128Rule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("complex128Rule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_boolRule(t *testing.T) {
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
			name: "test bool rule with true",
			args: args{
				fieldName:  "boolField",
				fieldValue: true,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test bool rule with false",
			args: args{
				fieldName:  "boolField",
				fieldValue: false,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test bool rule with non-bool value",
			args: args{
				fieldName:  "boolField",
				fieldValue: 1,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := boolRule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("boolRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("boolRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}

func Test_sliceRule(t *testing.T) {
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
			name: "test slice rule with false",
			args: args{
				fieldName:  "sliceField",
				fieldValue: []int{4, 2},
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test slice rule with non-slice value",
			args: args{
				fieldName:  "sliceField",
				fieldValue: 1,
				ruleValue:  "",
			},
			wantErr:           false,
			wantValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErr := sliceRule(tt.args.fieldName, tt.args.fieldValue, tt.args.ruleValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceRule() got = %v, want %v", err, tt.wantErr)
			}
			if (validationErr != "") != tt.wantValidationErr {
				t.Errorf("sliceRule() got = %v, want %v", validationErr, tt.wantValidationErr)
			}
		})
	}
}
