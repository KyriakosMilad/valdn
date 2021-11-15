package validation

import (
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
		name                string
		args                args
		wantErr             bool
		wantValidationError bool
	}{
		{
			name: "test validate field",
			args: args{
				fieldName:  "Name",
				fieldValue: "Kyria",
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:             false,
			wantValidationError: false,
		},
		{
			name: "test validate field with unsuitable data",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:             false,
			wantValidationError: true,
		},
		{
			name: "test validate field with not exist rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"bla:bla"},
			},
			wantErr:             true,
			wantValidationError: false,
		},
		{
			name: "test validate field with empty rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{},
			},
			wantErr:             false,
			wantValidationError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationError := Validate(tt.args.fieldName, tt.args.fieldValue, tt.args.fieldRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, errors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationError, tt.wantErr, tt.wantValidationError, tt.args)
			}
			if (validationError != "") != tt.wantValidationError {
				t.Errorf("Validate() error = %v, errors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationError, tt.wantErr, tt.wantValidationError, tt.args)
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
		name        string
		args        args
		wantErr     bool
		errorsCount int
	}{
		{
			name: "test validate struct with unsuitable data",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"required"}},
			},
			wantErr:     false,
			errorsCount: 2,
		},
		{
			name: "test validate struct with non exist rule",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"bla"}},
			},
			wantErr:     true,
			errorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"Name": ""},
				rules: Rules{"Name": {"required"}, "Age": {"required"}},
			},
			wantErr:     false,
			errorsCount: 2,
		},
		{
			name: "test validate map with non exist rule",
			args: args{
				val:   map[string]interface{}{"Name": ""},
				rules: Rules{"Name": {"bla"}},
			},
			wantErr:     true,
			errorsCount: 0,
		},
		{
			name: "test validate with non-struct and non-map value",
			args: args{
				val:   "test",
				rules: Rules{"Name": {"bla"}},
			},
			wantErr:     true,
			errorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, errors := ValidateNested(tt.args.val, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNested() err = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(len(errors), tt.errorsCount) {
				t.Errorf("ValidateNested() errors = %v, want %v", errors, tt.errorsCount)
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
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: `test validate json`,
			args: args{
				jsonData: `{"rName":"Ramses", "city":"Tiba"}`,
				rules:    Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: `test validate json with unsuitable data`,
			args: args{
				jsonData: `{"rName":"Ramses", "age":90}`,
				rules:    Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: `test validate json with non-json-string`,
			args: args{
				jsonData: `{"rName:"}`,
				rules:    Rules{},
			},
			wantErr:                       true,
			expectedValidationErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, errors := ValidateJson(tt.args.jsonData, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJson() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(errors) != tt.expectedValidationErrorsCount {
				t.Errorf("ValidateJson() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
		fieldsCountExpected int
	}{
		{
			name:                "test register field",
			args:                args{fieldName: "test"},
			fieldsCountExpected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(Rules{})
			v.registerField(tt.args.fieldName)
			if len(v.fieldsExist) != tt.fieldsCountExpected {
				t.Errorf("registerField(): can't register field, got %v, want %v", len(v.fieldsExist), tt.fieldsCountExpected)
			}
		})
	}
}

func Test_validation_addError(t *testing.T) {
	type args struct {
		fieldName string
		err       string
	}
	tests := []struct {
		name                string
		args                args
		errorsCountExpected int
	}{
		{
			name: "test add error",
			args: args{
				fieldName: "test",
				err:       "just test err",
			},
			errorsCountExpected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(Rules{})
			v.addError(tt.args.fieldName, tt.args.err)
			if len(v.errors) != tt.errorsCountExpected {
				t.Errorf("registerField(): can't register field, got %v, want %v", len(v.fieldsExist), tt.errorsCountExpected)
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

func Test_validation_addValidationTagRules(t *testing.T) {
	type Parent struct {
		Name string `validation:"required|string"`
		Age  int    `validation:"required|string"`
	}
	type args struct {
		t       reflect.Type
		rules   Rules
		parName string
	}
	tests := []struct {
		name           string
		args           args
		lengthExpected int
	}{
		{
			name: "test add validation tag rules",
			args: args{
				t:       reflect.TypeOf(Parent{}),
				rules:   Rules{"test": {"required"}},
				parName: "",
			},
			lengthExpected: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.addValidationTagRules(tt.args.t, tt.args.parName)
			if len(v.rules) != tt.lengthExpected {
				t.Errorf("addValidationTagRules() rules = %v len = %v, want %v", tt.args.rules, len(tt.args.rules), tt.lengthExpected)
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
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "validate struct",
			args: args{
				val:   Parent{Name: "Mina", Age: 26},
				rules: Rules{"Name": {"required", "kind:string"}, "Child.Name": {""}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate struct with unsuitable data",
			args: args{
				val:   Parent{Name: "Mina"},
				rules: Rules{"Name": {"required", "kind:string"}, "Age": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate nested struct with unsuitable data",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
				},
				rules: Rules{"Name": {"required", "kind:string"}, "Child.Name": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate nested struct with unsuitable data (validation tag)",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
				},
				rules: Rules{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
		{
			name: "validate struct containing string-key-map",
			args: args{
				val: Parent{
					Name: "Ikhnaton",
					Age:  2,
				},
				rules: Rules{"StringKeyMap": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.addValidationTagRules(reflect.TypeOf(tt.args.val), "")
			err := v.validateStruct(tt.args.val, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStruct() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(v.errors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateStruct() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate map",
			args: args{
				val:   map[string]interface{}{"rName": "Ramses", "city": "Tiba"},
				rules: Rules{"rName": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"rName": "Ramses", "age": 90},
				rules: Rules{"rName": {"required", "kind:int"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate nested map",
			args: args{
				val:   map[string]interface{}{"user": map[string]interface{}{"rName": "Kyriakos M.", "country": "Egypt"}},
				rules: Rules{"user": {"required"}, "user.rName": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"user": map[string]interface{}{"rName": 1, "country": "Egypt"}},
				rules: Rules{"user": {"required"}, "user.rName": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate nested map with non-string key",
			args: args{
				val:   map[string]interface{}{"user": map[int]interface{}{1: 2}},
				rules: Rules{"user": {"required"}},
			},
			wantErr:                       true,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map containing struct",
			args: args{
				val:   map[string]interface{}{"user": User{Name: 5}},
				rules: Rules{"user": {"required"}, "user.Name": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map containing struct with unsuitable data",
			args: args{
				val:   map[string]interface{}{"user": User{Name: 5}},
				rules: Rules{"user": {"required"}, "user.Name": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			err := v.validateMap(tt.args.val, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMap() error = %v, errors = %v, wantErr %v, wantValidationErrors %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(v.errors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateMap() error = %v, errors = %v, wantErr %v, wantValidationErrors %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
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
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			err := v.validateByType(tt.args.name, tt.args.typ, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateByType() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(v.errors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateByType() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
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
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
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
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			err := v.validateStructFields(tt.args.t, tt.args.v, tt.args.parName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStructFields() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(v.errors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateStructFields() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
		name                          string
		args                          args
		wantErr                       bool
		expectedValidationErrorsCount int
	}{
		{
			name: "test validate map fields",
			args: args{
				val:     map[string]interface{}{"Name": "Pola"},
				parName: "",
				rules:   Rules{"Parent.Name": {"required"}},
				errors:  map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				val:     map[string]interface{}{"Name": ""},
				parName: "Parent",
				rules:   Rules{"Parent.Name": {"required"}},
				errors:  map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			err := v.validateMapFields(tt.args.val, tt.args.parName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMapFields() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(v.errors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateMapFields() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, v.errors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
		errorsCountExpected int
	}{
		{
			name: "test validate non exist required fields",
			args: args{
				rules:       Rules{"field1": {"required", "string"}, "field3": {"string", "required"}},
				fieldsExist: map[string]bool{"field1": true, "field2": true},
				errors:      map[string]string{},
			},
			errorsCountExpected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.rules)
			v.fieldsExist = tt.args.fieldsExist
			v.validateNonExistRequiredFields()
			if len(v.errors) != tt.errorsCountExpected {
				t.Errorf("validateNonExistRequiredFields() errors = %v len = %v, want %v", v.errors, len(v.errors), tt.errorsCountExpected)
			}
		})
	}
}
