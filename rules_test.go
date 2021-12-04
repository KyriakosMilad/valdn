package validation

import (
	"mime/multipart"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_AddRule(t *testing.T) {
	type args struct {
		name   string
		f      RuleFunc
		errMsg string
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
			AddRule(tt.args.name, tt.args.f, tt.args.errMsg)
			if _, ok := registeredRules["test"]; !ok {
				t.Errorf("AddRule() error: failed to add rule, wantPanic: %v, error: %v, args: %v", tt.wantErr, nil, tt.args)
			}
		})
	}
}

func Test_OverwriteRule(t *testing.T) {
	type args struct {
		name   string
		f      RuleFunc
		errMsg string
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
			OverwriteRule(tt.args.name, tt.args.f, tt.args.errMsg)
			if _, ok := registeredRules["test0"]; !ok {
				t.Errorf("OverwriteRule() error: failed to overwrite rule, args: %v", tt.args)
			}
		})
	}
}

func Test_SetErrMsg(t *testing.T) {
	type args struct {
		ruleName string
		errMsg   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test set error message",
			args: args{
				ruleName: "test_add_err_msg",
				errMsg:   "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rFunc := func(fieldName string, fieldValue interface{}, ruleValue string) error {
				return nil
			}
			AddRule(tt.args.ruleName, rFunc, tt.args.errMsg)
			SetErrMsg(tt.args.ruleName, tt.args.errMsg)
			if registeredRules[tt.args.ruleName].errMsg != tt.args.errMsg {
				t.Errorf("SetErrMsg() can't set err msg, ruleName= %v, errMsg %v", tt.args.ruleName, tt.args.errMsg)
			}
		})
	}
}

func Test_getErrMsg(t *testing.T) {
	type args struct {
		ruleName string
		ruleVal  string
		name     string
		val      interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test get error message",
			args: args{
				ruleName: "kind",
				ruleVal:  "map",
				name:     "title",
				val:      44,
			},
			want: "title must be kind of map",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getErrMsg(tt.args.ruleName, tt.args.ruleVal, tt.args.name, tt.args.val); got != tt.want {
				t.Errorf("getErrMsg() = %v, want: %v", got, tt.want)
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
			f:         registeredRules["kind"].fn,
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
			err := requiredRule(tt.args.field, tt.args.fVal, tt.args.rVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("required rule: err: %v, wantErr: %v, args: %v", err, tt.wantErr, tt.args)
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
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test type rule with string",
			args: args{
				fName: "kindField",
				fVal:  "string",
				rVal:  "string",
			},
			wantErr: false,
		},
		{
			name: "test type rule with uint",
			args: args{
				fName: "kindField",
				fVal:  uint(44),
				rVal:  "uint",
			},
			wantErr: false,
		},
		{
			name: "test type rule with int",
			args: args{
				fName: "kindField",
				fVal:  -44,
				rVal:  "int",
			},
			wantErr: false,
		},
		{
			name: "test type rule with float",
			args: args{
				fName: "kindField",
				fVal:  44.44,
				rVal:  "float64",
			},
			wantErr: false,
		},
		{
			name: "test type rule with complex number",
			args: args{
				fName: "kindField",
				fVal:  44 + 22i,
				rVal:  "complex128",
			},
			wantErr: false,
		},
		{
			name: "test type rule with bool",
			args: args{
				fName: "kindField",
				fVal:  true,
				rVal:  "bool",
			},
			wantErr: false,
		},
		{
			name: "test type rule with map",
			args: args{
				fName: "kindField",
				fVal:  map[string]interface{}{"key": 55},
				rVal:  "map",
			},
			wantErr: false,
		},
		{
			name: "test type rule with struct",
			args: args{
				fName: "kindField",
				fVal:  struct{}{},
				rVal:  "struct",
			},
			wantErr: false,
		},
		{
			name: "test type rule with array",
			args: args{
				fName: "kindField",
				fVal:  [2]int{1, 2},
				rVal:  "array",
			},
			wantErr: false,
		},
		{
			name: "test type rule with slice",
			args: args{
				fName: "kindField",
				fVal:  []int{1, 2},
				rVal:  "slice",
			},
			wantErr: false,
		},
		{
			name: "test type rule with unsuitable data",
			args: args{
				fName: "kindField",
				fVal:  []int{1, 2},
				rVal:  "array",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := kindRule(tt.args.fName, tt.args.fVal, tt.args.rVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("kindRule() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kindInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test kindInRule",
			args: args{
				name:    "test",
				val:     8,
				ruleVal: "int,uint,int8",
			},
			wantErr: false,
		},
		{
			name: "test kindInRule with unsuitable data",
			args: args{
				name:    "test",
				val:     8,
				ruleVal: "string,slice,array",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kindInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("kindInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kindNotInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test kindNotInRule",
			args: args{
				name:    "test",
				val:     8,
				ruleVal: "string,slice,array",
			},
			wantErr: false,
		},
		{
			name: "test kindNotInRule with unsuitable data",
			args: args{
				name:    "test",
				val:     8,
				ruleVal: "int,uint,int8",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kindNotInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("kindNotInRule() error = %v, wantErr %v", err, tt.wantErr)
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
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test type rule with string",
			args: args{
				fieldName: "typeField",
				fVal:      "string",
				rVal:      "string",
			},
			wantErr: false,
		},
		{
			name: "test type rule with uint",
			args: args{
				fieldName: "typeField",
				fVal:      uint(44),
				rVal:      "uint",
			},
			wantErr: false,
		},
		{
			name: "test type rule with int",
			args: args{
				fieldName: "typeField",
				fVal:      -44,
				rVal:      "int",
			},
			wantErr: false,
		},
		{
			name: "test type rule with float",
			args: args{
				fieldName: "typeField",
				fVal:      44.44,
				rVal:      "float64",
			},
			wantErr: false,
		},
		{
			name: "test type rule with complex number",
			args: args{
				fieldName: "typeField",
				fVal:      44 + 22i,
				rVal:      "complex128",
			},
			wantErr: false,
		},
		{
			name: "test type rule with bool",
			args: args{
				fieldName: "typeField",
				fVal:      true,
				rVal:      "bool",
			},
			wantErr: false,
		},
		{
			name: "test type rule with map",
			args: args{
				fieldName: "typeField",
				fVal:      map[string]interface{}{"key": 55},
				rVal:      "map[string]interface {}",
			},
			wantErr: false,
		},
		{
			name: "test type rule with struct",
			args: args{
				fieldName: "typeField",
				fVal:      user{name: "test"},
				rVal:      "user",
			},
			wantErr: false,
		},
		{
			name: "test type rule with array",
			args: args{
				fieldName: "typeField",
				fVal:      [2]int{1, 2},
				rVal:      "[2]int",
			},
			wantErr: false,
		},
		{
			name: "test type rule with slice",
			args: args{
				fieldName: "typeField",
				fVal:      []int{1, 2},
				rVal:      "[]int",
			},
			wantErr: false,
		},
		{
			name: "test type rule with unsuitable data",
			args: args{
				fieldName: "typeField",
				fVal:      []int{1, 2},
				rVal:      "[2]int",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := typeRule(tt.args.fieldName, tt.args.fVal, tt.args.rVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("typeRule() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_typeInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test typeInRule",
			args: args{
				name:    "test",
				val:     [2]int{1, 2},
				ruleVal: "[2]int,[]int,[1]int",
			},
			wantErr: false,
		},
		{
			name: "test typeInRule with unsuitable data",
			args: args{
				name:    "test",
				val:     [2]int{1, 2},
				ruleVal: "[3]int,[]int,[1]int",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := typeInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("typeInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_typeNotInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test typeNotInRule",
			args: args{
				name:    "test",
				val:     [2]int{1, 2},
				ruleVal: "[3]int,[]int,[1]int",
			},
			wantErr: false,
		},
		{
			name: "test typeNotInRule with unsuitable data",
			args: args{
				name:    "test",
				val:     [2]int{1, 2},
				ruleVal: "[2]int,[]int,[1]int",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := typeNotInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("typeNotInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_equalRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test equal rule",
			args: args{
				name:    "age",
				val:     22,
				ruleVal: "22",
			},
			wantErr: false,
		},
		{
			name: "test equal rule with unsuitable value",
			args: args{
				name:    "age",
				val:     22,
				ruleVal: "23",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := equalRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("equalRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_intRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test int rule",
			args: args{
				name:    "age",
				val:     -15,
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test int rule with unsuitable data",
			args: args{
				name:    "age",
				val:     1.5,
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := intRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("intRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uintRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test uint rule",
			args: args{
				name:    "age",
				val:     uint(1),
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test uint rule with unsuitable data",
			args: args{
				name:    "age",
				val:     -1,
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uintRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("uintRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_complexRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test complex rule",
			args: args{
				name:    "complex",
				val:     2 + 2i,
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test complex rule with unsuitable data",
			args: args{
				name:    "age",
				val:     -1,
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := complexRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("complexRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_floatRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test float rule",
			args: args{
				name:    "float",
				val:     2.2,
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test float rule with unsuitable data",
			args: args{
				name:    "age",
				val:     5,
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := floatRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("floatRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ufloatRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test ufloat rule",
			args: args{
				name:    "price",
				val:     2.5,
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test ufloat rule with unsuitable data",
			args: args{
				name:    "price",
				val:     -1.1,
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ufloatRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("ufloatRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_numericRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test numeric rule",
			args: args{
				name:    "age",
				val:     15,
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test numeric rule with unsuitable data",
			args: args{
				name:    "age",
				val:     "t",
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := numericRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("numericRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_betweenRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test between rule with float val and integer rule val",
			args: args{
				name:    "price",
				val:     5.5,
				ruleVal: "3,6",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test between rule with integer val and integer rule val",
			args: args{
				name:    "price",
				val:     5,
				ruleVal: "3,6",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test between rule with float val and float rule val",
			args: args{
				name:    "price",
				val:     5.5,
				ruleVal: "3,6.5",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test between rule with string",
			args: args{
				name:    "price",
				val:     "55",
				ruleVal: "3,6",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test between rule with empty rule val",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test between rule with non-numeric rule val",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "bla,bb",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test between rule with unsuitable data",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "5,6",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("betweenRule() error = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			if err := betweenRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("betweenRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_minRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test min rule with float val and integer rule val",
			args: args{
				name:    "price",
				val:     5.5,
				ruleVal: "3",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test min rule with integer val and integer rule val",
			args: args{
				name:    "price",
				val:     5,
				ruleVal: "3.6",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test min rule with string",
			args: args{
				name:    "price",
				val:     "55",
				ruleVal: "3",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test min rule with empty rule val",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test min rule with non-numeric rule val",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "bla",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test min rule with unsuitable data",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "5",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("minRule() error = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			if err := minRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("minRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_maxRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test max rule with float val and integer rule val",
			args: args{
				name:    "price",
				val:     2.5,
				ruleVal: "3",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test max rule with integer val and integer rule val",
			args: args{
				name:    "price",
				val:     3,
				ruleVal: "3.6",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test max rule with string",
			args: args{
				name:    "price",
				val:     "55",
				ruleVal: "3",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test max rule with empty rule val",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test max rule with non-numeric rule val",
			args: args{
				name:    "price",
				val:     4,
				ruleVal: "bla",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test max rule with unsuitable data",
			args: args{
				name:    "price",
				val:     6,
				ruleVal: "5",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("maxRule() error = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			if err := maxRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("maxRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test inRule",
			args: args{
				name:    "country",
				val:     "EGYPT",
				ruleVal: "GREECE,EGYPT,CYPRUS",
			},
			wantErr: false,
		},
		{
			name: "test inRule with unsuitable data",
			args: args{
				name:    "country",
				val:     "EGYPTtttt",
				ruleVal: "GREECE,EGYPT,CYPRUS",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := inRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("inRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_notInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test notInRule",
			args: args{
				name:    "country",
				val:     "EGYPTtttt",
				ruleVal: "GREECE,EGYPT,CYPRUS",
			},
			wantErr: false,
		},
		{
			name: "test notInRule with unsuitable data",
			args: args{
				name:    "country",
				val:     "EGYPT",
				ruleVal: "GREECE,EGYPT,CYPRUS",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := notInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("notInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lenRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test lenRule with slice",
			args: args{
				name:    "test",
				ruleVal: "3",
				val:     []int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with array",
			args: args{
				name:    "test",
				ruleVal: "3",
				val:     [3]int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with array",
			args: args{
				name:    "test",
				ruleVal: "1",
				val:     map[int]string{1: "test"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with string",
			args: args{
				name:    "test",
				ruleVal: "4",
				val:     "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with integer",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     -55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with unsigned integer",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with float",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     -555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with unsigned float",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenRule with unsuitable ruleVal",
			args: args{
				name:    "test",
				ruleVal: "bla",
				val:     44,
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenRule with struct",
			args: args{
				name:    "test",
				ruleVal: "0",
				val:     struct{}{},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenRule with unsuitable data",
			args: args{
				name:    "test",
				ruleVal: "2",
				val:     333,
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("lenRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := lenRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("lenRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_minLenRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test minLenRule with slice",
			args: args{
				name:    "test",
				ruleVal: "3",
				val:     []int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with array",
			args: args{
				name:    "test",
				ruleVal: "3",
				val:     [3]int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with array",
			args: args{
				name:    "test",
				ruleVal: "1",
				val:     map[int]string{1: "test"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with string",
			args: args{
				name:    "test",
				ruleVal: "4",
				val:     "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with integer",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     -55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with unsigned integer",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with float",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     -555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with unsigned float",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test minLenRule with unsuitable ruleVal",
			args: args{
				name:    "test",
				ruleVal: "bla",
				val:     44,
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test minLenRule with struct",
			args: args{
				name:    "test",
				ruleVal: "0",
				val:     struct{}{},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test minLenRule with unsuitable data",
			args: args{
				name:    "test",
				ruleVal: "4",
				val:     333,
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("minLenRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := minLenRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("minLenRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_maxLenRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test maxLenRule with slice",
			args: args{
				name:    "test",
				ruleVal: "3",
				val:     []int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with array",
			args: args{
				name:    "test",
				ruleVal: "3",
				val:     [3]int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with array",
			args: args{
				name:    "test",
				ruleVal: "1",
				val:     map[int]string{1: "test"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with string",
			args: args{
				name:    "test",
				ruleVal: "4",
				val:     "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with integer",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     -55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with unsigned integer",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with float",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     -555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with unsigned float",
			args: args{
				name:    "test",
				ruleVal: "5",
				val:     555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test maxLenRule with unsuitable ruleVal",
			args: args{
				name:    "test",
				ruleVal: "bla",
				val:     44,
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test maxLenRule with struct",
			args: args{
				name:    "test",
				ruleVal: "0",
				val:     struct{}{},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test maxLenRule with unsuitable data",
			args: args{
				name:    "test",
				ruleVal: "2",
				val:     333,
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("maxLenRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := maxLenRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("maxLenRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lenBetweenRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test lenBetweenRule with slice",
			args: args{
				name:    "test",
				ruleVal: "3,6",
				val:     []int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with array",
			args: args{
				name:    "test",
				ruleVal: "3,3",
				val:     [3]int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with array",
			args: args{
				name:    "test",
				ruleVal: "1,2",
				val:     map[int]string{1: "test"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with string",
			args: args{
				name:    "test",
				ruleVal: "4,10",
				val:     "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with integer",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     -55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with unsigned integer",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with float",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     -555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with unsigned float",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenBetweenRule with unsuitable ruleVal",
			args: args{
				name:    "test",
				ruleVal: "4,",
				val:     44,
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenBetweenRule with struct",
			args: args{
				name:    "test",
				ruleVal: "0,0",
				val:     struct{}{},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenBetweenRule with unsuitable data",
			args: args{
				name:    "test",
				ruleVal: "1,2",
				val:     333,
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("lenBetweenRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := lenBetweenRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("lenBetweenRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lenInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test lenInRule with slice",
			args: args{
				name:    "test",
				ruleVal: "3,6",
				val:     []int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with array",
			args: args{
				name:    "test",
				ruleVal: "3,3",
				val:     [3]int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with array",
			args: args{
				name:    "test",
				ruleVal: "1,2",
				val:     map[int]string{1: "test"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with string",
			args: args{
				name:    "test",
				ruleVal: "4,10",
				val:     "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with integer",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     -55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with unsigned integer",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with float",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     -555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with unsigned float",
			args: args{
				name:    "test",
				ruleVal: "5,7",
				val:     555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenInRule with unsuitable ruleVal",
			args: args{
				name:    "test",
				ruleVal: "4,",
				val:     44,
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenInRule with struct",
			args: args{
				name:    "test",
				ruleVal: "0,0",
				val:     struct{}{},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenInRule with unsuitable data",
			args: args{
				name:    "test",
				ruleVal: "1,2",
				val:     333,
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("lenInRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := lenInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("lenInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lenNotInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test lenNotInRule with slice",
			args: args{
				name:    "test",
				ruleVal: "4,6",
				val:     []int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with array",
			args: args{
				name:    "test",
				ruleVal: "5,4",
				val:     [3]int{1, 2, 3},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with array",
			args: args{
				name:    "test",
				ruleVal: "2",
				val:     map[int]string{1: "test"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with string",
			args: args{
				name:    "test",
				ruleVal: "5,10",
				val:     "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with integer",
			args: args{
				name:    "test",
				ruleVal: "6,7",
				val:     -55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with unsigned integer",
			args: args{
				name:    "test",
				ruleVal: "6,7",
				val:     55555,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with float",
			args: args{
				name:    "test",
				ruleVal: "6,7",
				val:     -555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with unsigned float",
			args: args{
				name:    "test",
				ruleVal: "6,7",
				val:     555.55,
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test lenNotInRule with unsuitable ruleVal",
			args: args{
				name:    "test",
				ruleVal: "4,",
				val:     44,
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenNotInRule with struct",
			args: args{
				name:    "test",
				ruleVal: "0,0",
				val:     struct{}{},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test lenNotInRule with unsuitable data",
			args: args{
				name:    "test",
				ruleVal: "1,2,3",
				val:     333,
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("lenNotInRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := lenNotInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("lenNotInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_regexRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test regexRule",
			args: args{
				name:    "test",
				val:     "email@email.com",
				ruleVal: "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test regexRule with non-string value",
			args: args{
				name:    "test",
				val:     55,
				ruleVal: "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test regexRule with invalid regex",
			args: args{
				name:    "test",
				val:     "55",
				ruleVal: "[",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test regexRule with unsuitable data",
			args: args{
				name:    "test",
				val:     "emailemail.com",
				ruleVal: "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("regexRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := regexRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("regexRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_notRegexRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test notRegexRule",
			args: args{
				name:    "test",
				val:     "emailemail.com",
				ruleVal: "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test notRegexRule with non-string value",
			args: args{
				name:    "test",
				val:     55,
				ruleVal: "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test notRegexRule with invalid regex",
			args: args{
				name:    "test",
				val:     "55",
				ruleVal: "[",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test notRegexRule with unsuitable data",
			args: args{
				name:    "test",
				val:     "email@email.com",
				ruleVal: "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("notRegexRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := notRegexRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("notRegexRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_emailRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test emailRule",
			args: args{
				name:    "test",
				val:     "email@email.com",
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test emailRule with non-string value",
			args: args{
				name:    "test",
				val:     55,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test emailRule with unsuitable data",
			args: args{
				name:    "test",
				val:     "emailemail.com",
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("emailRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := emailRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("emailRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_jsonRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test jsonRule",
			args: args{
				name:    "test",
				val:     `{"name":"Ramses", "city":"Tiba"}`,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test jsonRule with non-string value",
			args: args{
				name:    "test",
				val:     1973,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test jsonRule with unsuitable data",
			args: args{
				name:    "test",
				val:     `name:"Ramses`,
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("jsonRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := jsonRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("jsonRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ipv4Rule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test ipv4Rule",
			args: args{
				name:    "test",
				val:     "255.255.255.255",
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test ipv4Rule with non-string value",
			args: args{
				name:    "test",
				val:     1973,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test ipv4Rule with unsuitable data",
			args: args{
				name:    "test",
				val:     "1.1.1.1.",
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ipv4Rule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := ipv4Rule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("ipv4Rule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ipv6Rule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test ipv6Rule",
			args: args{
				name:    "test",
				val:     "2001:db8:3333:4444:5555:6666:1.2.3.4",
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test ipv6Rule with non-string value",
			args: args{
				name:    "test",
				val:     1973,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test ipv6Rule with unsuitable data",
			args: args{
				name:    "test",
				val:     "1.1.1.1.",
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ipv6Rule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := ipv6Rule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("ipv6Rule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ipRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test ipRule",
			args: args{
				name:    "test",
				val:     "2001:db8:3333:4444:5555:6666:1.2.3.4",
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test ipRule with non-string value",
			args: args{
				name:    "test",
				val:     1973,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test ipRule with unsuitable data",
			args: args{
				name:    "test",
				val:     "bla",
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ipRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := ipRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("ipRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_macRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test macRule",
			args: args{
				name:    "test",
				val:     "3D:F2:C9:A6:B3:4F",
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test macRule with non-string value",
			args: args{
				name:    "test",
				val:     1973,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test macRule with unsuitable data",
			args: args{
				name:    "test",
				val:     "bla",
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("macRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := macRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("macRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_urlRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test urlRule",
			args: args{
				name:    "test",
				val:     "presidency.eg",
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test urlRule with non-string value",
			args: args{
				name:    "test",
				val:     1973,
				ruleVal: "",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test urlRule with unsuitable data",
			args: args{
				name:    "test",
				val:     "presidency_eg",
				ruleVal: "",
			},
			wantErr:   true,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("urlRule() error = %v, wantPanic %v", e, tt.wantErr)
				}
			}()
			if err := urlRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("urlRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_timeRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test timeRule",
			args: args{
				name:    "createdAt",
				val:     time.Now(),
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test timeRule with unsuitable data",
			args: args{
				name:    "createdAt",
				val:     "6/10/1973",
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := timeRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("timeRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_timeFormatRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test timeFormatRule",
			args: args{
				name:    "updatedAt",
				val:     "06/10/1973",
				ruleVal: "02/01/2006",
			},
			wantErr: false,
		},
		{
			name: "test timeFormatRule with unsuitable data",
			args: args{
				name:    "updatedAt",
				val:     "06/13/1973",
				ruleVal: "02/01/2006",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := timeFormatRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("timeFormatRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_timeFormatInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test timeFormatInRule",
			args: args{
				name:    "deletedAt",
				val:     "06/13/1973",
				ruleVal: "01/02/2006,02/01/2006",
			},
			wantErr: false,
		},
		{
			name: "test timeFormatInRule with unsuitable data",
			args: args{
				name:    "deletedAt",
				val:     "06/10/1973",
				ruleVal: "02/01,Jan 2006",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := timeFormatInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("timeFormatInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_timeFormatNotInRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test timeFormatNotInRule",
			args: args{
				name:    "deletedAt",
				val:     "06/10/1973",
				ruleVal: "02/01,Jan 2006",
			},
			wantErr: false,
		},
		{
			name: "test timeFormatNotInRule with unsuitable data",
			args: args{
				name:    "deletedAt",
				val:     "06/13/1973",
				ruleVal: "01/02/2006,02/01/2006",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := timeFormatNotInRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("timeFormatNotInRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileRule(t *testing.T) {
	f, err := os.Open("example.json")
	if err != nil {
		panic(err)
	}
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test fileRule with multipart.FileHeaders",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test fileRule with os.File",
			args: args{
				name:    "file",
				val:     f,
				ruleVal: "",
			},
			wantErr: false,
		},
		{
			name: "test fileRule with unsuitable data",
			args: args{
				name:    "file",
				val:     "bal bla",
				ruleVal: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fileRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("fileRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sizeRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test sizeRule",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "44",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test sizeRule with unsuitable data",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "43",
			},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "test sizeRule non-integer ruleVal",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "forty",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test sizeRule non-file val",
			args: args{
				name:    "file",
				val:     "bla bla",
				ruleVal: "forty",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("sizeRule() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			if err := sizeRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("sizeRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sizeMinRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
		wantErr   bool
	}{
		{
			name: "test sizeMinRule",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "44",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test sizeMinRule with unsuitable data",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "45",
			},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "test sizeMinRule non-integer ruleVal",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "forty",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test sizeMinRule non-file val",
			args: args{
				name:    "file",
				val:     "bla bla",
				ruleVal: "forty",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("sizeMinRule() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			if err := sizeMinRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("sizeMinRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sizeMaxRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test sizeMaxRule",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "44",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test sizeMaxRule with unsuitable data",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "43",
			},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "test sizeMaxRule non-integer ruleVal",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "forty",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test sizeMaxRule non-file val",
			args: args{
				name:    "file",
				val:     "bla bla",
				ruleVal: "forty",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("sizeMaxRule() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			if err := sizeMaxRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("sizeMaxRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sizeBetweenRule(t *testing.T) {
	type args struct {
		name    string
		val     interface{}
		ruleVal string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
		wantErr   bool
	}{
		{
			name: "test sizeBetweenRule",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "44,45",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test sizeBetweenRule with unsuitable data",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "1,40",
			},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "test sizeBetweenRule with unsuitable ruleVal",
			args: args{
				name:    "file",
				val:     multipart.FileHeader{Size: 44},
				ruleVal: "66",
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test sizeBetweenRule with non-file val",
			args: args{
				name:    "file",
				val:     "bla bla",
				ruleVal: "40,55",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("sizeBetweenRule() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			if err := sizeBetweenRule(tt.args.name, tt.args.val, tt.args.ruleVal); (err != nil) != tt.wantErr {
				t.Errorf("sizeBetweenRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
