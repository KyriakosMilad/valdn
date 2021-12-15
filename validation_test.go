package validation

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
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
			if v := createNewValidation(tt.args.rules); reflect.TypeOf(v) != reflect.TypeOf(&validation{}) {
				t.Errorf("createNewValidation() error can't create new validation = %v", v)
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
			name: "test validate field",
			args: args{
				fieldName:  "Name",
				fieldValue: "Kyria",
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "test validate field with unsuitable data",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "test validate field with non exist rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"bla:bla"},
			},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name: "test validate field with empty rule",
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

func Test_ValidateNested(t *testing.T) {
	type Child struct {
		Name string `validation:"required|kind:string"`
		Age  int
	}
	type args struct {
		val   interface{}
		rules Rules
	}
	tests := []struct {
		name                string
		args                args
		wantPanic           bool
		expectedErrorsCount int
	}{
		{
			name: "test validate struct with unsuitable data",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 2,
		},
		{
			name: "test validate struct with non exist rule",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"bla"}},
			},
			wantPanic:           true,
			expectedErrorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"Name": ""},
				rules: Rules{"Name": {"required"}, "Age": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 2,
		},
		{
			name: "test validate map with non exist rule",
			args: args{
				val:   map[string]interface{}{"Name": ""},
				rules: Rules{"Name": {"bla"}},
			},
			wantPanic:           true,
			expectedErrorsCount: 0,
		},
		{
			name: "test validate with non-struct and non-map value",
			args: args{
				val:   "test",
				rules: Rules{"Name": {"bla"}},
			},
			wantPanic:           true,
			expectedErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateNested() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			vErrors := ValidateNested(tt.args.val, tt.args.rules)
			if !reflect.DeepEqual(len(vErrors), tt.expectedErrorsCount) {
				t.Errorf("ValidateNested() vErrors = %v, want %v", vErrors, tt.expectedErrorsCount)
			}
		})
	}
}

func Test_ValidateJson(t *testing.T) {
	type args struct {
		jsonData string
		rules    Rules
	}
	tests := []struct {
		name                string
		args                args
		wantPanic           bool
		expectedErrorsCount int
	}{
		{
			name: `test validate json`,
			args: args{
				jsonData: `{"rName":"Ramses", "city":"Tiba"}`,
				rules:    Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 0,
		},
		{
			name: `test validate json with unsuitable data`,
			args: args{
				jsonData: `{"rName":"Ramses", "age":90}`,
				rules:    Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 1,
		},
		{
			name: `test validate json with non-json-string`,
			args: args{
				jsonData: `{"rName:"}`,
				rules:    Rules{},
			},
			wantPanic:           true,
			expectedErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateJson() panicEr = %v, wantPanic %v, args %v", e, tt.wantPanic, tt.args)
				}
			}()
			vErrors := ValidateJson(tt.args.jsonData, tt.args.rules)
			if len(vErrors) != tt.expectedErrorsCount {
				t.Errorf("ValidateJson() vErrors = %v, expectedErrorsCount %v, args %v", vErrors, tt.expectedErrorsCount, tt.args)
			}
		})
	}
}

func Test_validation_registerField(t *testing.T) {
	type args struct {
		fieldName string
	}
	tests := []struct {
		name                string
		args                args
		expectedFieldsCount int
	}{
		{
			name:                "test register field",
			args:                args{fieldName: "test"},
			expectedFieldsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(Rules{})
			v.registerField(tt.args.fieldName)
			if len(v.fieldsExist) != tt.expectedFieldsCount {
				t.Errorf("registerField(): can't register field, got %v, want %v", len(v.fieldsExist), tt.expectedFieldsCount)
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
		name                string
		args                args
		expectedErrorsCount int
	}{
		{
			name: "test add error",
			args: args{
				name: "test",
				err:  errors.New("just test err"),
			},
			expectedErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(Rules{})
			v.addError(tt.args.name, tt.args.err)
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("registerField(): can't register field, got %v, want %v", len(v.fieldsExist), tt.expectedErrorsCount)
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
		name               string
		fields             fields
		args               args
		rulesCountExpected int
	}{
		{
			name:               "test get field rules",
			fields:             fields{rules: Rules{"test": {"required", "king:string"}}},
			args:               args{fieldName: "test"},
			rulesCountExpected: 2,
		},
		{
			name:               "test get field rules using .*",
			fields:             fields{rules: Rules{"parent.*": {"required", "king:string"}}},
			args:               args{fieldName: "parent.test"},
			rulesCountExpected: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.fields.rules)
			if fieldRules := v.getFieldRules(tt.args.fieldName); len(fieldRules) != tt.rulesCountExpected {
				t.Errorf("getFieldRules() error getting field rules, got = %v, want %v", len(fieldRules), tt.rulesCountExpected)
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
		name               string
		fields             fields
		args               args
		rulesCountExpected int
	}{
		{
			name:               "test get parent rules",
			fields:             fields{rules: Rules{"test": {"required", "king:string"}}},
			args:               args{fieldName: "test"},
			rulesCountExpected: 2,
		},
		{
			name:               "test get parent rules using .*",
			fields:             fields{rules: Rules{"parent.*": {"required", "king:string"}}},
			args:               args{fieldName: "parent.test"},
			rulesCountExpected: 2,
		},
		{
			name:               "test get parent rules with empty name",
			fields:             fields{rules: Rules{"parent.*": {"required", "king:string"}}},
			args:               args{fieldName: ""},
			rulesCountExpected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.fields.rules)
			if fieldRules := v.getParentRules(tt.args.fieldName); len(fieldRules) != tt.rulesCountExpected {
				t.Errorf("getParentRules() error getting parent rules, got = %v, want %v", len(fieldRules), tt.rulesCountExpected)
			}
		})
	}
}

func Test_validation_addTagRules(t *testing.T) {
	type Parent struct {
		Name string `validation:"required|string"`
		Age  int    `validation:"required|string"`
	}
	type args struct {
		val     interface{}
		rules   Rules
		parName string
	}
	tests := []struct {
		name           string
		args           args
		lengthExpected int
	}{
		{
			name: "test add tag rules",
			args: args{
				val:     Parent{},
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			lengthExpected: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.addTagRules(tt.args.val, reflect.TypeOf(tt.args.val), tt.args.parName)
			if len(v.rules) != tt.lengthExpected {
				t.Errorf("addTagRules() rules = %v len = %v, want %v", tt.args.rules, len(tt.args.rules), tt.lengthExpected)
			}
		})
	}
}

func Test_validation_validateStruct(t *testing.T) {
	type Child struct {
		Name string `validation:"required|kind:string"`
	}
	type Parent struct {
		Name         string `validation:"required|kind:string"`
		Age          int    `validation:"required"`
		StringKeyMap map[string]interface{}
		Child
	}
	type args struct {
		val   interface{}
		rules Rules
	}
	tests := []struct {
		name                string
		args                args
		expectedErrorsCount int
	}{
		{
			name: "validate struct",
			args: args{
				val:   Parent{Name: "Mina", Age: 26},
				rules: Rules{"Name": {"required", "kind:string"}, "Child.Name": {""}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "validate struct with unsuitable data",
			args: args{
				val:   Parent{Name: "Mina"},
				rules: Rules{"Name": {"required", "kind:string"}, "Age": {"required"}},
			},
			expectedErrorsCount: 2,
		},
		{
			name: "validate nested struct",
			args: args{
				val: Parent{
					Name:  "Ikhnaton",
					Child: Child{Name: "Tut"},
				},
				rules: Rules{"Name": {"required", "kind:string"}, "Age": {""}, "Child.Name": {"required", "kind:string"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "validate nested struct with unsuitable data",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
				},
				rules: Rules{"Name": {"required", "kind:string"}, "Child.Name": {"required", "kind:string"}},
			},
			expectedErrorsCount: 2,
		},
		{
			name: "validate nested struct (validation tag)",
			args: args{
				val: Parent{
					Name:  "Ikhnaton",
					Age:   1,
					Child: Child{Name: "tut"},
				},
				rules: Rules{},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "validate nested struct with unsuitable data (validation tag)",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
				},
				rules: Rules{},
			},
			expectedErrorsCount: 2,
		},
		{
			name: "validate struct containing map",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
					Age:  2,
				},
				rules: Rules{"StringKeyMap": {"required"}},
			},
			expectedErrorsCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.addTagRules(tt.args.val, reflect.TypeOf(tt.args.val), "")
			v.validateStruct(tt.args.val, "")
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateStruct() errors = %v expectedErrorsCount %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
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
		name                string
		args                args
		expectedErrorsCount int
		wantPanic           bool
	}{
		{
			name: "test validate map",
			args: args{
				val:   map[string]interface{}{"rName": "Ramses", "city": "Tiba"},
				rules: Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"rName": "Ramses", "age": 90},
				rules: Rules{"rName": {"required", "kind:int"}},
			},
			expectedErrorsCount: 1,
		},
		{
			name: "test validate nested map",
			args: args{
				val:   map[string]interface{}{"user": map[string]interface{}{"rName": "Kyriakos M.", "country": "Egypt"}},
				rules: Rules{"user": {"required"}, "user.rName": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"user": map[string]interface{}{"rName": 1, "country": "Egypt"}},
				rules: Rules{"user": {"required"}, "user.rName": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			expectedErrorsCount: 1,
		},
		{
			name: "test validate nested map with non-string key",
			args: args{
				val:   map[string]interface{}{"user": map[int]interface{}{1: 2}},
				rules: Rules{"user": {"required"}},
			},
			expectedErrorsCount: 0,
			wantPanic:           true,
		},
		{
			name: "test validate map containing struct",
			args: args{
				val:   map[string]interface{}{"user": User{Name: 5}},
				rules: Rules{"user": {"required"}, "user.Name": {"required"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate map containing struct with unsuitable data",
			args: args{
				val:   map[string]interface{}{"user": User{Name: 5}},
				rules: Rules{"user": {"required"}, "user.Name": {"required", "kind:string"}},
			},
			expectedErrorsCount: 1,
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
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateMap() errors = %v, wantValidationErrors %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
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
		name                string
		args                args
		expectedErrorsCount int
	}{
		{
			name: "test validate by type - struct",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf(Parent{}),
				val:         Parent{Name: "Barssom"},
				rules:       Rules{"Parent.Name": {"required"}},
				fieldsExist: map[string]bool{"Parent.Name": true},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate by type - struct with unsuitable data",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf(Parent{}),
				val:         Parent{},
				rules:       Rules{"Parent.Name": {"required"}},
				fieldsExist: map[string]bool{},
			},
			expectedErrorsCount: 1,
		},
		{
			name: "test validate by type - map",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf(map[string]interface{}{}),
				val:         map[string]interface{}{"Name": "Pola"},
				rules:       Rules{"Parent.Name": {"required"}},
				fieldsExist: map[string]bool{"Parent.Name": true},
			},
			expectedErrorsCount: 0,
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
			expectedErrorsCount: 1,
		},
		{
			name: "test validate by type - slice",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf([]interface{}{}),
				val:         []interface{}{"Pola"},
				rules:       Rules{"Parent.0": {"kind:string"}},
				fieldsExist: map[string]bool{"Parent.0": true},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate by slice - slice with unsuitable data",
			args: args{
				name:        "Parent",
				typ:         reflect.TypeOf([]interface{}{}),
				val:         []interface{}{"Pola"},
				rules:       Rules{"Parent.0": {"kind:int"}},
				fieldsExist: map[string]bool{"Parent.0": true},
			},
			expectedErrorsCount: 1,
		},
		{
			name: "test validate by type - field",
			args: args{
				name:        "age",
				typ:         reflect.TypeOf(1),
				val:         1,
				rules:       Rules{"age": {"required"}},
				fieldsExist: map[string]bool{"age": true},
			},
			expectedErrorsCount: 0,
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
			expectedErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateByType(tt.args.name, tt.args.typ, tt.args.val)
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateByType()errors = %v, expectedErrorsCount %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
			}
		})
	}
}

func Test_validation_validateStructFields(t *testing.T) {
	type Parent struct {
		Name string `validation:"required|string"`
	}
	type args struct {
		t       reflect.Type
		v       reflect.Value
		parName string
		rules   Rules
		errors  map[string]string
	}
	tests := []struct {
		name                string
		args                args
		expectedErrorsCount int
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
			expectedErrorsCount: 0,
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
			expectedErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateStructFields(tt.args.t, tt.args.v, tt.args.parName)
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateStructFields() errors = %v, expectedErrorsCount %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
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
		name                string
		args                args
		expectedErrorsCount int
	}{
		{
			name: "test validate map fields",
			args: args{
				val:     map[string]interface{}{"Name": "Pola"},
				parName: "",
				rules:   Rules{"Parent.Name": {"required"}},
				errors:  map[string]string{},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				val:     map[string]interface{}{"Name": ""},
				parName: "Parent",
				rules:   Rules{"Parent.Name": {"required"}},
				errors:  map[string]string{},
			},
			expectedErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateMapFields(tt.args.val, tt.args.parName)
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateMapFields() errors = %v, expectedErrorsCount %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
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
		name                string
		args                args
		expectedErrorsCount int
	}{
		{
			name: "test validate non exist required fields",
			args: args{
				rules:       Rules{"field1": {"required", "string"}, "field3": {"string", "required"}},
				fieldsExist: map[string]bool{"field1": true, "field2": true},
				errors:      map[string]string{},
			},
			expectedErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.fieldsExist = tt.args.fieldsExist
			v.validateNonExistRequiredFields()
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateNonExistRequiredFields() errors = %v len = %v, want %v", v.errors, len(v.errors), tt.expectedErrorsCount)
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
		name                string
		args                args
		expectedErrorsCount int
		wantPanic           bool
	}{
		{
			name: "test requestToMap with multipart/form-data",
			args: args{
				r:     formDataRequest(),
				rules: Rules{"field1": {"required"}, "field2": {"required"}, "file": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 0,
		},
		{
			name: "test requestToMap with application/x-www-form-urlencoded",
			args: args{
				r:     urlencodedRequest(),
				rules: Rules{"lang": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 0,
		},
		{
			name: "test requestToMap with application/json",
			args: args{
				r:     jsonRequest(),
				rules: Rules{"lang": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 0,
		},
		{
			name: "test requestToMap with url params",
			args: args{
				r:     paramsRequest(),
				rules: Rules{"lang": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 0,
		},
		{
			name: "test requestToMap with empty json",
			args: args{
				r:     emptyJSONRequest(),
				rules: Rules{"lang": {"required"}},
			},
			wantPanic:           true,
			expectedErrorsCount: 0,
		},
		{
			name: "test requestToMap with unsuitable data",
			args: args{
				r:     formDataRequest(),
				rules: Rules{"lang": {"required"}},
			},
			wantPanic:           false,
			expectedErrorsCount: 1,
		},
		{
			name: "test requestToMap with rule does not exist",
			args: args{
				r:     formDataRequest(),
				rules: Rules{"lang": {"bla bla"}},
			},
			wantPanic:           true,
			expectedErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) && !tt.wantPanic {
					t.Errorf("ValidateRequest() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			if got := ValidateRequest(tt.args.r, tt.args.rules); len(got) != tt.expectedErrorsCount {
				t.Errorf("ValidateRequest() = %v, erros count expected %v", got, tt.expectedErrorsCount)
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
		name                string
		args                args
		expectedErrorsCount int
		wantPanic           bool
	}{
		{
			name: "test validate slice",
			args: args{
				val:   []interface{}{"Ramses"},
				rules: Rules{".*": {"required", "kind:string"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate slice with unsuitable data",
			args: args{
				val:   []interface{}{"Ramses"},
				rules: Rules{"0": {"required", "kind:int"}},
			},
			expectedErrorsCount: 1,
		},
		{
			name: "test validate nested slice",
			args: args{
				val:   []interface{}{[]interface{}{"Egypt"}},
				rules: Rules{".*": {"kind:slice"}, "0.*": {"kind:string"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate nested slice with unsuitable data",
			args: args{
				val:   []interface{}{[]interface{}{1973}},
				rules: Rules{".*": {"kind:slice"}, "0.*": {"kind:string"}},
			},
			expectedErrorsCount: 1,
		},
		{
			name: "test validate nested slice with unsuitable type",
			args: args{
				val:   []interface{}{[]string{"Mina"}},
				rules: Rules{"*": {"required"}},
			},
			expectedErrorsCount: 0,
			wantPanic:           true,
		},
		{
			name: "test validate slice containing map",
			args: args{
				val:   []interface{}{map[string]interface{}{"user": 1973}},
				rules: Rules{".*": {"required"}, "0.Name": {"required"}},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate slice containing struct",
			args: args{
				val:   []interface{}{User{Name: 1973}},
				rules: Rules{".*": {"required"}, "0.Name": {"required", "kind:string"}},
			},
			expectedErrorsCount: 1,
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
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateSlice() errors = %v, wantValidationErrors %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
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
		name                string
		args                args
		expectedErrorsCount int
	}{
		{
			name: "test validate map fields",
			args: args{
				val:     []interface{}{"Pola"},
				parName: "",
				rules:   Rules{"Parent.0": {"required"}},
				errors:  map[string]string{},
			},
			expectedErrorsCount: 0,
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				val:     []interface{}{"Pola"},
				parName: "Parent",
				rules:   Rules{"Parent.0": {"kind:int"}},
				errors:  map[string]string{},
			},
			expectedErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.validateSliceFields(tt.args.val, tt.args.parName)
			if len(v.errors) != tt.expectedErrorsCount {
				t.Errorf("validateSliceFields() errors = %v, expectedErrorsCount %v, args %v", v.errors, tt.expectedErrorsCount, tt.args)
			}
		})
	}
}
