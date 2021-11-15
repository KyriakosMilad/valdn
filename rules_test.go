package validation

import (
	"reflect"
	"testing"
)

func Test_AddRule(t *testing.T) {
	type args struct {
		name string
		f    RuleFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test add rule",
			args: args{
				name: "test",
				f: func(name string, fVal interface{}, rVal string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "test add rule already exist",
			args: args{
				name: "test",
				f: func(name string, fVal interface{}, rVal string) error {
					return nil
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil && !tt.wantErr {
					t.Errorf("AddRule() error: failed to add rule, wantPanic: %v, error: %v, args: %v", tt.wantErr, err, tt.args)
				}
			}()
			AddRule(tt.args.name, tt.args.f)
			if _, ok := registeredRules["test"]; !ok {
				t.Errorf("AddRule() error: failed to add rule, wantPanic: %v, error: %v, args: %v", tt.wantErr, nil, tt.args)
			}
		})
	}
}

func Test_OverwriteRule(t *testing.T) {
	type args struct {
		name string
		f    RuleFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test overwrite rule",
			args: args{
				name: "test0",
				f: func(field string, fVal interface{}, rVal string) error {
					return nil
				},
			},
		},
		{
			name: "test overwrite rule already exist",
			args: args{
				name: "test0",
				f: func(field string, fVal interface{}, rVal string) error {
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OverwriteRule(tt.args.name, tt.args.f)
			if _, ok := registeredRules["test0"]; !ok {
				t.Errorf("OverwriteRule() error: failed to overwrite rule, args: %v", tt.args)
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
		rName     string
		rVal      string
		f         RuleFunc
		ruleExist bool
	}{
		{
			name: "test get rule info",
			args: args{
				rule: "kind:string",
			},
			rName:     "kind",
			rVal:      "string",
			f:         registeredRules["kind"],
			ruleExist: true,
		},
		{
			name: "test get info of non exist rule",
			args: args{
				rule: "string",
			},
			rName:     "string",
			rVal:      "",
			f:         nil,
			ruleExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, rVal, f, ruleExist := getRuleInfo(tt.args.rule)
			if name != tt.rName {
				t.Errorf("getRuleInfo() rName = %v, want %v", name, tt.rName)
			}
			if rVal != tt.rVal {
				t.Errorf("getRuleInfo() rVal = %v, want %v", rVal, tt.rVal)
			}
			if !reflect.DeepEqual(toString(f), toString(tt.f)) {
				t.Errorf("getRuleInfo() f = %v, want %v", toString(f), toString(tt.f))
			}
			if ruleExist != tt.ruleExist {
				t.Errorf("getRuleInfo() ruleExist = %v, want %v", ruleExist, tt.ruleExist)
			}
		})
	}
}

func Test_requiredRule(t *testing.T) {
	type args struct {
		field string
		fVal  interface{}
		rVal  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test required rule",
			args: args{
				field: "rName",
				fVal:  "Kyriakos",
				rVal:  "",
			},
			wantErr: false,
		},
		{
			name: "test required rule with zero value",
			args: args{
				field: "rName",
				fVal:  "",
				rVal:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requiredFunc, requiredExist := registeredRules["required"]
			if !requiredExist {
				panic("required rule is not exist")
			}
			err := requiredFunc(tt.args.field, tt.args.fVal, tt.args.rVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("required rule: err: %v, wantErr: %v, args: %v", err, tt.wantErr, tt.args)
			}
		})
	}
}

func Test_typeRule(t *testing.T) {
	type user struct {
		name string
	}
	type args struct {
		fieldName string
		fVal      interface{}
		rVal      string
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
				fieldName: "typeField",
				fVal:      "string",
				rVal:      "string",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with uint",
			args: args{
				fieldName: "typeField",
				fVal:      uint(44),
				rVal:      "uint",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with int",
			args: args{
				fieldName: "typeField",
				fVal:      -44,
				rVal:      "int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with float",
			args: args{
				fieldName: "typeField",
				fVal:      44.44,
				rVal:      "float64",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with complex number",
			args: args{
				fieldName: "typeField",
				fVal:      44 + 22i,
				rVal:      "complex128",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with bool",
			args: args{
				fieldName: "typeField",
				fVal:      true,
				rVal:      "bool",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with map",
			args: args{
				fieldName: "typeField",
				fVal:      map[string]interface{}{"key": 55},
				rVal:      "map[string]interface {}",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with struct",
			args: args{
				fieldName: "typeField",
				fVal:      user{name: "test"},
				rVal:      "user",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with array",
			args: args{
				fieldName: "typeField",
				fVal:      [2]int{1, 2},
				rVal:      "[2]int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with slice",
			args: args{
				fieldName: "typeField",
				fVal:      []int{1, 2},
				rVal:      "[]int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := typeRule(tt.args.fieldName, tt.args.fVal, tt.args.rVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("typeRule() got = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kindRule(t *testing.T) {
	type args struct {
		fName string
		fVal  interface{}
		rVal  string
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
				fName: "kindField",
				fVal:  "string",
				rVal:  "string",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with uint",
			args: args{
				fName: "kindField",
				fVal:  uint(44),
				rVal:  "uint",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with int",
			args: args{
				fName: "kindField",
				fVal:  -44,
				rVal:  "int",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with float",
			args: args{
				fName: "kindField",
				fVal:  44.44,
				rVal:  "float64",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with complex number",
			args: args{
				fName: "kindField",
				fVal:  44 + 22i,
				rVal:  "complex128",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with bool",
			args: args{
				fName: "kindField",
				fVal:  true,
				rVal:  "bool",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with map",
			args: args{
				fName: "kindField",
				fVal:  map[string]interface{}{"key": 55},
				rVal:  "map",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with struct",
			args: args{
				fName: "kindField",
				fVal:  struct{}{},
				rVal:  "struct",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with array",
			args: args{
				fName: "kindField",
				fVal:  [2]int{1, 2},
				rVal:  "array",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
		{
			name: "test type rule with slice",
			args: args{
				fName: "kindField",
				fVal:  []int{1, 2},
				rVal:  "slice",
			},
			wantErr:           false,
			wantValidationErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := kindRule(tt.args.fName, tt.args.fVal, tt.args.rVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("kindRule() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
