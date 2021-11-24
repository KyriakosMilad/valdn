package validation

import (
	"reflect"
	"testing"
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
