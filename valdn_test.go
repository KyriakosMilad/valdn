package valdn

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_createNewValidation(t *testing.T) {
	type args struct {
		rules Rules
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test create new validation",
			args: args{
				rules: Rules{"test": {"required"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := &validation{rules: Rules{"test": {"required"}}, errors: make(Errors), fieldsExist: make(fieldsExist)}
			if got := createNewValidation(tt.args.rules); !reflect.DeepEqual(want, got) {
				t.Errorf("createNewValidation() = %v, want = %v", got, want)
			}
		})
	}
}

func Test_Validate(t *testing.T) {
	type args struct {
		fieldName  string
		fieldValue interface{}
		fieldRules []string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "test validate",
			args: args{
				fieldName:  "Name",
				fieldValue: "Kyria",
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test validate with unsuitable data",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "test validate with non exist rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"bla:bla"},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test validate with empty rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{},
			},
			wantErr:   false,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("Validate() error = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			err := Validate(tt.args.fieldName, tt.args.fieldValue, tt.args.fieldRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v, args %v", err, tt.wantErr, tt.args)
			}
		})
	}
}

func Test_ValidateCollection(t *testing.T) {
	type User struct {
		ID          int64     `json:"id" db:"id"`
		Name        string    `json:"name" db:"name" valdn:"required|minLen:5:maxLen:30"`
		Email       string    `json:"email" db:"email" valdn:"required|email"`
		Phone       string    `json:"phone" db:"phone" valdn:"required|minLen:5:maxLen:20"`
		CountryCode string    `json:"country_code" db:"country_code" valdn:"required|len:2"`
		CreatedAt   time.Time `json:"created_at" db:"created_at" valdn:"skip"`
		UpdatedAt   time.Time `json:"updated_at" db:"updated_at" valdn:"skip"`
	}
	type args struct {
		val   interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate collection with unsuitable kind",
			args: args{
				val:   "string kind",
				rules: nil,
			},
			want:      nil,
			wantPanic: true,
		},
		{
			name: "test validate collection with struct",
			args: args{
				val: User{
					Name:        "kyrikos",
					Email:       "test@test.test",
					Phone:       "15125125125",
					CountryCode: "eg",
				},
				rules: Rules{},
			},
			want:      Errors{},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateCollection() panicErr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			if got := ValidateCollection(tt.args.val, tt.args.rules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateStruct(t *testing.T) {
	type Child struct {
		Name string `valdn:"required|kind:string"`
		Age  int
	}
	type args struct {
		val   interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate nested struct",
			args: args{
				val:   Child{Name: "Mina"},
				rules: Rules{},
			},
			want:      Errors{},
			wantPanic: false,
		},
		{
			name: "test validate nested struct with unsuitable data",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"required"}},
			},
			want: Errors{
				"Name": GetErrMsg("required", "", "Name", ""),
				"Age":  GetErrMsg("required", "", "Age", 0),
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateStruct() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			got := ValidateCollection(tt.args.val, tt.args.rules)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateMap(t *testing.T) {
	type args struct {
		val   map[interface{}]interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate nested map",
			args: args{
				val:   map[interface{}]interface{}{"Age": 44},
				rules: Rules{"Age": {"required"}},
			},
			want:      Errors{},
			wantPanic: false,
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				val:   map[interface{}]interface{}{"Age": 44},
				rules: Rules{"Name": {"required"}},
			},
			want: Errors{
				"Name": GetErrMsg("required", "", "Name", ""),
			},
			wantPanic: false,
		},
		{
			name: "test validate nested map with non string key",
			args: args{
				val:   map[interface{}]interface{}{14: 44},
				rules: Rules{"Name": {"required"}},
			},
			want: Errors{
				"Name": GetErrMsg("required", "", "Name", ""),
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateMap() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			got := ValidateCollection(tt.args.val, tt.args.rules)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateSlice(t *testing.T) {
	type args struct {
		val   interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate nested slice",
			args: args{
				val:   []interface{}{44},
				rules: Rules{"0": {"kind:int"}},
			},
			want:      Errors{},
			wantPanic: false,
		},
		{
			name: "test validate nested slice with unsuitable data",
			args: args{
				val:   []int{44},
				rules: Rules{"0": {"kind:string"}},
			},
			want: Errors{
				"0": GetErrMsg("kind", "string", "0", ""),
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateSlice() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			got := ValidateCollection(tt.args.val, tt.args.rules)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateArray(t *testing.T) {
	type args struct {
		val   interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate nested array",
			args: args{
				val:   [1]interface{}{44},
				rules: Rules{"0": {"kind:int"}},
			},
			want:      Errors{},
			wantPanic: false,
		},
		{
			name: "test validate nested array with unsuitable data",
			args: args{
				val:   [1]int{44},
				rules: Rules{"0": {"kind:string"}},
			},
			want: Errors{
				"0": GetErrMsg("kind", "string", "0", ""),
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateSlice() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			got := ValidateCollection(tt.args.val, tt.args.rules)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateJSON(t *testing.T) {
	type args struct {
		val   string
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: `test validate json`,
			args: args{
				val:   `{"rName":"Ramses", "city":"Tiba"}`,
				rules: Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantPanic: false,
			want:      Errors{},
		},
		{
			name: `test validate json with unsuitable data`,
			args: args{
				val:   `{"rName":"Ramses", "age":90}`,
				rules: Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantPanic: false,
			want: Errors{
				"city": GetErrMsg("required", "", "city", 0),
			}},
		{
			name: `test validate json with non-json-string`,
			args: args{
				val:   `{"rName"}`,
				rules: Rules{},
			},
			wantPanic: true,
			want:      Errors{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateJson() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			got := ValidateJSON(tt.args.val, tt.args.rules)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateRequest(t *testing.T) {
	type args struct {
		r     *http.Request
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test ValidateRequest with multipart/form-data",
			args: args{
				r:     formDataRequest(),
				rules: Rules{"field1": {"required"}, "field2": {"required"}, "field3": {"required"}, "file": {"required"}},
			},
			wantPanic: false,
			want:      Errors{"field3": GetErrMsg("required", "", "field3", "")},
		},
		{
			name: "test ValidateRequest with application/x-www-form-urlencoded",
			args: args{
				r:     urlencodedRequest(),
				rules: Rules{"lang": {"required", "kind:int"}},
			},
			wantPanic: false,
			want:      Errors{"lang": GetErrMsg("kind", "int", "lang", "go")},
		},
		{
			name: "test ValidateRequest with application/json",
			args: args{
				r:     jsonRequest(),
				rules: Rules{"lang": {"required", "kind:int"}},
			},
			wantPanic: false,
			want:      Errors{"lang": GetErrMsg("kind", "int", "lang", "go")},
		},
		{
			name: "test ValidateRequest with url params",
			args: args{
				r:     paramsRequest(),
				rules: Rules{"lang": {"required", "kind:int"}},
			},
			wantPanic: false,
			want:      Errors{"lang": GetErrMsg("kind", "int", "lang", "go")},
		},
		{
			name: "test ValidateRequest with empty json",
			args: args{
				r:     emptyJSONRequest(),
				rules: Rules{},
			},
			wantPanic: true,
			want:      Errors{},
		},
		{
			name: "test ValidateRequest with unsuitable data",
			args: args{
				r:     formDataRequest(),
				rules: Rules{"lang": {"required"}},
			},
			wantPanic: false,
			want:      Errors{"lang": GetErrMsg("required", "", "lang", "")},
		},
		{
			name: "test ValidateRequest with rule does not exist",
			args: args{
				r:     formDataRequest(),
				rules: Rules{},
			},
			wantPanic: true,
			want:      Errors{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateRequest() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			got := ValidateRequest(tt.args.r, tt.args.rules)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validation_registerField(t *testing.T) {
	type args struct {
		fieldName string
	}
	tests := []struct {
		name string
		args args
		want fieldsExist
	}{
		{
			name: "test register field",
			args: args{fieldName: "test"},
			want: fieldsExist{"test": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(Rules{})
			v.registerField(tt.args.fieldName)
			if !reflect.DeepEqual(v.fieldsExist, tt.want) {
				t.Errorf("registerField() = %v, want %v", v.fieldsExist, tt.want)
			}
		})
	}
}

func Test_validation_addError(t *testing.T) {
	type args struct {
		name string
		err  error
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "test add error",
			args: args{
				name: "test",
				err:  errors.New("just test err"),
			},
			want: Errors{"test": "just test err"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(Rules{})
			v.addError(tt.args.name, tt.args.err)
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("addError() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_getFieldRules(t *testing.T) {
	type fields struct {
		rules Rules
	}
	type args struct {
		fieldName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "test get field rules",
			fields: fields{rules: Rules{"test": {"required", "king:string"}}},
			args:   args{fieldName: "test"},
			want:   []string{"required", "king:string"},
		},
		{
			name:   "test get field rules using .*",
			fields: fields{rules: Rules{"parent.*": {"required", "king:string"}}},
			args:   args{fieldName: "parent.test"},
			want:   []string{"required", "king:string"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.fields.rules)
			got := v.getFieldRules(tt.args.fieldName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validation_getParentRules(t *testing.T) {
	type fields struct {
		rules Rules
	}
	type args struct {
		fieldName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "test get parent rules",
			fields: fields{rules: Rules{"test": {"required", "king:string"}}},
			args:   args{fieldName: "test"},
			want:   []string{"required", "king:string"},
		},
		{
			name:   "test get parent rules using .*",
			fields: fields{rules: Rules{"parent.*": {"required", "king:string"}}},
			args:   args{fieldName: "parent.test"},
			want:   []string{"required", "king:string"},
		},
		{
			name:   "test get parent rules with empty name",
			fields: fields{rules: Rules{"parent.*": {"required", "king:string"}}},
			args:   args{fieldName: ""},
			want:   []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.fields.rules)
			got := v.getParentRules(tt.args.fieldName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getParentRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validation_addTagRules(t *testing.T) {
	type Parent struct {
		Name string `valdn:"required|string"`
		Age  int    `valdn:"required|int"`
	}
	type args struct {
		val     interface{}
		rules   Rules
		parName string
	}
	tests := []struct {
		name string
		args args
		want Rules
	}{
		{
			name: "test add tag rules to struct",
			args: args{
				val:     Parent{},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "Name": {"required", "string"}, "Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct field that already has rules",
			args: args{
				val:     Parent{},
				rules:   Rules{"test": {"required"}, "Name": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "Name": {"required"}, "Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct inside map",
			args: args{
				val:     map[string]interface{}{"parent": Parent{}},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "parent.Name": {"required", "string"}, "parent.Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct inside slice",
			args: args{
				val:     []interface{}{Parent{}},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "0.Name": {"required", "string"}, "0.Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct inside array",
			args: args{
				val:     [1]interface{}{Parent{}},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "0.Name": {"required", "string"}, "0.Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct inside slice inside map",
			args: args{
				val:     []interface{}{map[string]interface{}{"parent": Parent{}}},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "0.parent.Name": {"required", "string"}, "0.parent.Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct inside map inside slice",
			args: args{
				val:     map[string]interface{}{"parent": []interface{}{Parent{}}},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "parent.0.Name": {"required", "string"}, "parent.0.Age": {"required", "int"}},
		},
		{
			name: "test add tag rules to struct inside map inside array",
			args: args{
				val:     map[string]interface{}{"parent": [1]interface{}{Parent{}}},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			want: Rules{"test": {"required"}, "parent.0.Name": {"required", "string"}, "parent.0.Age": {"required", "int"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.addTagRules(tt.args.val, tt.args.parName)
			if !reflect.DeepEqual(v.rules, tt.want) {
				t.Errorf("addTagRules() = %v, want %v", v.rules, tt.want)
			}
		})
	}
}

func Test_validation_validateStruct(t *testing.T) {
	type Child struct {
		Name string
	}
	type Parent struct {
		Name           string
		Age            int
		StringKeyMap   map[string]interface{}
		InterfaceSlice []interface{}
		Child
	}
	type unexported struct {
		name string
	}
	type args struct {
		val   interface{}
		rules Rules
		name  string
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "validate struct",
			args: args{
				val:   Parent{Name: "Mina", Age: 26},
				rules: Rules{"parent": {"kind:struct"}},
			},
			want: Errors{},
		},
		{
			name: "validate nested struct",
			args: args{
				val: Parent{
					Name:  "Ikhnaton",
					Child: Child{Name: "Tut"},
				},
				rules: Rules{"Name": {"required", "kind:string"}, "Child.Name": {"required", "kind:string"}},
			},
			want: Errors{},
		},
		{
			name: "validate nested struct with unsuitable data",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
				},
				rules: Rules{"Name": {"required", "kind:string"}, "Child.Name": {"required"}},
			},
			want: Errors{"Child.Name": GetErrMsg("required", "", "Child.Name", "")},
		},
		{
			name: "validate struct contains map with unsuitable data",
			args: args{
				val: Parent{
					Name:         "Ikhnaton",
					Age:          2,
					StringKeyMap: map[string]interface{}{"name": "Tia"},
				},
				rules: Rules{"StringKeyMap.name": {"kind:int"}},
			},
			want: Errors{"StringKeyMap.name": GetErrMsg("kind", "int", "StringKeyMap.name", "Tia")},
		},
		{
			name: "validate struct contains slice with unsuitable data",
			args: args{
				val: Parent{
					Name:           "Ikhnaton",
					Age:            2,
					InterfaceSlice: []interface{}{"Tia"},
				},
				rules: Rules{"InterfaceSlice.0": {"kind:int"}},
			},
			want: Errors{"InterfaceSlice.0": GetErrMsg("kind", "int", "InterfaceSlice.0", "Tia")},
		},
		{
			name: "validate struct has unexported fields",
			args: args{
				val: unexported{
					name: "unexported field",
				},
				rules: Rules{"name": {"required", "kind:string"}},
			},
			want: Errors{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateStruct(tt.args.val, tt.args.name)
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateStruct() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateMap(t *testing.T) {
	type User struct {
		Name int
	}
	type args struct {
		val   map[string]interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate map",
			args: args{
				val:   map[string]interface{}{"Name": "Ramses", "city": "Tiba"},
				rules: Rules{"Name": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			want: Errors{},
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"Name": "Ramses", "age": 90},
				rules: Rules{"Name": {"required", "kind:int"}},
			},
			want: Errors{"Name": GetErrMsg("kind", "int", "Name", "Ramses")},
		},
		{
			name: "test validate nested map",
			args: args{
				val:   map[string]interface{}{"user": map[string]interface{}{"rName": "Kyriakos M.", "country": "Egypt"}},
				rules: Rules{"user": {"required"}, "user.rName": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			want: Errors{},
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"user": map[string]interface{}{"Name": 1}},
				rules: Rules{"user.Name": {"required", "kind:string"}},
			},
			want: Errors{"user.Name": GetErrMsg("kind", "string", "user.Name", 1)},
		},
		{
			name: "test validate nested map with non-string key",
			args: args{
				val:   map[string]interface{}{"user": map[int]interface{}{1: 2}},
				rules: Rules{"user": {"required"}},
			},
			want:      Errors{},
			wantPanic: true,
		},
		{
			name: "test validate map containing struct with unsuitable data",
			args: args{
				val:   map[string]interface{}{"user": User{Name: 5}},
				rules: Rules{"user.Name": {"required", "kind:string"}},
			},
			want: Errors{"user.Name": GetErrMsg("kind", "string", "user.Name", 1)},
		},
		{
			name: "validate struct contains slice with unsuitable data",
			args: args{
				val:   map[string]interface{}{"slice": []interface{}{"Tia"}},
				rules: Rules{"slice.0": {"kind:int"}},
			},
			want: Errors{"slice.0": GetErrMsg("kind", "int", "slice.0", "Tia")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("validateMap() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			v := createNewValidation(tt.args.rules)
			v.validateMap(tt.args.val, "")
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateMap() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateSlice(t *testing.T) {
	type User struct {
		Name int
	}
	type args struct {
		val   []interface{}
		rules Rules
	}
	tests := []struct {
		name      string
		args      args
		want      Errors
		wantPanic bool
	}{
		{
			name: "test validate slice",
			args: args{
				val:   []interface{}{"Ramses"},
				rules: Rules{".*": {"required", "kind:string"}},
			},
			want: Errors{},
		},
		{
			name: "test validate slice with unsuitable data",
			args: args{
				val:   []interface{}{"Ramses"},
				rules: Rules{"0": {"required", "kind:int"}},
			},
			want: Errors{"0": GetErrMsg("kind", "int", "0", "Ramses")},
		},
		{
			name: "test validate nested slice",
			args: args{
				val:   []interface{}{[]interface{}{"Egypt"}},
				rules: Rules{".*": {"kind:slice"}, "0.*": {"kind:string"}},
			},
			want: Errors{},
		},
		{
			name: "test validate nested slice with unsuitable data",
			args: args{
				val:   []interface{}{[]interface{}{1973}},
				rules: Rules{".*": {"kind:slice"}, "0.*": {"kind:string"}},
			},
			want: Errors{"0.0": GetErrMsg("kind", "string", "0.0", 1973)},
		},
		{
			name: "test validate nested slice with unsuitable type",
			args: args{
				val:   []interface{}{[]string{"Mina"}},
				rules: Rules{"*": {"required"}},
			},
			want:      Errors{},
			wantPanic: true,
		},
		{
			name: "test validate slice containing map",
			args: args{
				val:   []interface{}{map[string]interface{}{"user": 1973}},
				rules: Rules{".*": {"required"}, "0.user": {"required", "kind:string"}},
			},
			want: Errors{"0.user": GetErrMsg("kind", "string", "0.user", 1973)},
		},
		{
			name: "test validate slice containing struct",
			args: args{
				val:   []interface{}{User{Name: 1973}},
				rules: Rules{".*": {"required"}, "0.Name": {"required", "kind:string"}},
			},
			want: Errors{"0.Name": GetErrMsg("kind", "string", "0.Name", 1973)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("validateSlice() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			v := createNewValidation(tt.args.rules)
			v.validateSlice(tt.args.val, "")
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateSlice() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateByType(t *testing.T) {
	type Parent struct {
		Name string
	}
	type args struct {
		name        string
		typ         reflect.Type
		val         interface{}
		rules       Rules
		fieldsExist map[string]bool
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "test validate by type - struct with unsuitable data",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf(Parent{}),
				val:         Parent{},
				rules:       Rules{"Parent.Name": {"required"}},
				fieldsExist: map[string]bool{},
			},
			want: Errors{"Parent.Name": GetErrMsg("required", "", "Parent.Name", "")},
		},
		{
			name: "test validate by type - map with unsuitable data",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf(map[string]interface{}{}),
				val:         map[string]interface{}{"Name": 22},
				rules:       Rules{"Parent.Name": {"kind:string"}},
				fieldsExist: map[string]bool{},
			},
			want: Errors{"Parent.Name": GetErrMsg("kind", "string", "Parent.Name", 22)},
		},
		{
			name: "test validate by type - slice with unsuitable data",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf([]interface{}{}),
				val:         []interface{}{"Pola"},
				rules:       Rules{"Parent.0": {"kind:int"}},
				fieldsExist: map[string]bool{"Parent.0": true},
			},
			want: Errors{"Parent.0": GetErrMsg("kind", "int", "Parent.0", "")},
		},
		{
			name: "test validate by type - array with unsuitable data",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf([]interface{}{}),
				val:         [1]interface{}{"Pola"},
				rules:       Rules{"Parent.0": {"kind:int"}},
				fieldsExist: map[string]bool{"Parent.0": true},
			},
			want: Errors{"Parent.0": GetErrMsg("kind", "int", "Parent.0", "")},
		},
		{
			name: "test validate by type - field with unsuitable data",
			args: args{
				name:        "age",
				typ:         reflect.TypeOf(1),
				val:         0,
				rules:       Rules{"age": {"required"}},
				fieldsExist: map[string]bool{},
			},
			want: Errors{"age": GetErrMsg("required", "", "age", "")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateByType(tt.args.name, tt.args.typ, tt.args.val)
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateByType() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateStructFields(t *testing.T) {
	type Parent struct {
		Name string `valdn:"required|string"`
	}
	type unexported struct {
		name string
	}
	type args struct {
		t       reflect.Type
		v       reflect.Value
		parName string
		rules   Rules
		errors  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "test validate struct fields",
			args: args{
				t:       reflect.TypeOf(Parent{Name: "Sherry"}),
				v:       reflect.ValueOf(Parent{Name: "Sherry"}),
				parName: "",
				rules:   Rules{"Name": {"required"}},
				errors:  map[string]string{},
			},
			want: Errors{},
		},
		{
			name: "test validate struct fields with unsuitable data",
			args: args{
				t:       reflect.TypeOf(Parent{}),
				v:       reflect.ValueOf(Parent{}),
				parName: "",
				rules:   Rules{"Name": {"required"}},
				errors:  map[string]string{},
			},
			want: Errors{"Name": GetErrMsg("required", "", "Name", "")},
		},
		{
			name: "test validate unexported field",
			args: args{
				t:       reflect.TypeOf(unexported{name: "unexported field"}),
				v:       reflect.ValueOf(unexported{name: "unexported field"}),
				parName: "",
				rules:   Rules{"name": {"required", "kind:string"}},
			},
			want: Errors{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateStructFields(tt.args.t, tt.args.v, tt.args.parName)
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateStructFields() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateMapFields(t *testing.T) {
	type args struct {
		val     map[string]interface{}
		parName string
		rules   Rules
		errors  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "test validate map fields",
			args: args{
				val:     map[string]interface{}{"Name": "Pola"},
				parName: "",
				rules:   Rules{"Parent.Name": {"required"}},
				errors:  map[string]string{},
			},
			want: Errors{},
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				val:     map[string]interface{}{"Name": ""},
				parName: "Parent",
				rules:   Rules{"Parent.Name": {"required"}},
				errors:  map[string]string{},
			},
			want: Errors{"Parent.Name": GetErrMsg("required", "", "Parent.Name", "")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateMapFields(tt.args.val, tt.args.parName)
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateMapFields() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateSliceFields(t *testing.T) {
	type args struct {
		val     []interface{}
		parName string
		rules   Rules
		errors  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "test validate map fields",
			args: args{
				val:     []interface{}{"Pola"},
				parName: "",
				rules:   Rules{"Parent.0": {"required"}},
				errors:  map[string]string{},
			},
			want: Errors{},
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				val:     []interface{}{"Pola"},
				parName: "Parent",
				rules:   Rules{"Parent.0": {"kind:int"}},
				errors:  map[string]string{},
			},
			want: Errors{"Parent.0": GetErrMsg("kind", "int", "Parent.0", "Pola")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateSliceFields(tt.args.val, tt.args.parName)
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateSliceFields() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}

func Test_validation_validateNonExistRequiredFields(t *testing.T) {
	type args struct {
		rules       Rules
		fieldsExist map[string]bool
		errors      map[string]string
	}
	tests := []struct {
		name string
		args args
		want Errors
	}{
		{
			name: "test validate non exist required fields",
			args: args{
				rules:       Rules{"field1": {"required", "string"}, "field3": {"string", "required"}},
				fieldsExist: map[string]bool{"field1": true, "field2": true},
				errors:      map[string]string{},
			},
			want: Errors{"field3": GetErrMsg("required", "", "field3", "")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.fieldsExist = tt.args.fieldsExist
			v.validateNonExistRequiredFields()
			if !reflect.DeepEqual(v.errors, tt.want) {
				t.Errorf("validateNonExistRequiredFields() = %v, want %v", v.errors, tt.want)
			}
		})
	}
}
