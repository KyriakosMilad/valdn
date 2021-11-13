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
			name: "test validate field with not exists rule",
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
		name                  string
		args                  args
		wantErr               bool
		validationErrorsCount int
	}{
		{
			name: "test validate struct with unsuitable data",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"required"}},
			},
			wantErr:               false,
			validationErrorsCount: 2,
		},
		{
			name: "test validate struct with non exist rule",
			args: args{
				val:   Child{},
				rules: Rules{"Age": {"bla"}},
			},
			wantErr:               true,
			validationErrorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				val:   map[string]interface{}{"Name": ""},
				rules: Rules{"Name": {"required"}, "Age": {"required"}},
			},
			wantErr:               false,
			validationErrorsCount: 2,
		},
		{
			name: "test validate map with non exist rule",
			args: args{
				val:   map[string]interface{}{"Name": ""},
				rules: Rules{"Name": {"bla"}},
			},
			wantErr:               true,
			validationErrorsCount: 0,
		},
		{
			name: "test validate with non-struct and non-map value",
			args: args{
				val:   "test",
				rules: Rules{"Name": {"bla"}},
			},
			wantErr:               true,
			validationErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateNested(tt.args.val, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNested() err = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(len(validationErrors), tt.validationErrorsCount) {
				t.Errorf("ValidateNested() errors = %v, want %v", validationErrors, tt.validationErrorsCount)
			}
		})
	}
}

func Test_ValidateJson(t *testing.T) {
	type args struct {
		jsonData        string
		validationRules map[string][]string
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
				jsonData:        `{"name":"Ramses", "city":"Tiba"}`,
				validationRules: map[string][]string{"name": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: `test validate json with unsuitable data`,
			args: args{
				jsonData:        `{"name":"Ramses", "age":90}`,
				validationRules: map[string][]string{"name": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: `test validate json with non-json-string`,
			args: args{
				jsonData:        `{"name:"}`,
				validationRules: map[string][]string{},
			},
			wantErr:                       true,
			expectedValidationErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateJson(tt.args.jsonData, tt.args.validationRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJson() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("ValidateJson() error = %v, errors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
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
			if len(v.fieldsExists) != tt.fieldsCountExpected {
				t.Errorf("registerField(): can't register field, got %v, want %v", len(v.fieldsExists), tt.fieldsCountExpected)
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
				t.Errorf("registerField(): can't register field, got %v, want %v", len(v.fieldsExists), tt.errorsCountExpected)
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
		t               reflect.Type
		validationRules map[string][]string
		parentName      string
	}
	tests := []struct {
		name           string
		args           args
		lengthExpected int
	}{
		{
			name: "test add validation tag rules",
			args: args{
				t:               reflect.TypeOf(Parent{}),
				validationRules: map[string][]string{"test": {"required"}},
				parentName:      "",
			},
			lengthExpected: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			v.addValidationTagRules(tt.args.t, tt.args.parentName)
			if len(v.rules) != tt.lengthExpected {
				t.Errorf("addValidationTagRules() rules = %v len = %v, want %v", tt.args.validationRules, len(tt.args.validationRules), tt.lengthExpected)
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
		structData      interface{}
		validationRules map[string][]string
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
				structData:      Parent{Name: "Mina", Age: 26},
				validationRules: map[string][]string{"Name": {"required", "kind:string"}, "Child.Name": {""}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate struct with unsuitable data",
			args: args{
				structData:      Parent{Name: "Mina"},
				validationRules: map[string][]string{"Name": {"required", "kind:string"}, "Age": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
		{
			name: "validate nested struct",
			args: args{
				structData: Parent{
					Name:  "Ikhnaton",
					Child: Child{Name: "Tut"},
				},
				validationRules: map[string][]string{"Name": {"required", "kind:string"}, "Age": {""}, "Child.Name": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate nested struct with unsuitable data",
			args: args{
				structData: Parent{
					Name: "Ikhnaton",
				},
				validationRules: map[string][]string{"Name": {"required", "kind:string"}, "Child.Name": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
		{
			name: "validate nested struct (validation tag)",
			args: args{
				structData: Parent{
					Name:  "Ikhnaton",
					Age:   1,
					Child: Child{Name: "tut"},
				},
				validationRules: map[string][]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate nested struct with unsuitable data (validation tag)",
			args: args{
				structData: Parent{
					Name: "Ikhnaton",
				},
				validationRules: map[string][]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
		{
			name: "validate struct includes string-key-map",
			args: args{
				structData: Parent{
					Name: "Ikhnaton",
					Age:  2,
				},
				validationRules: map[string][]string{"StringKeyMap": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			v.addValidationTagRules(reflect.TypeOf(tt.args.structData), "")
			err := v.validateStruct(tt.args.structData, "")
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
		mapData         map[string]interface{}
		validationRules map[string][]string
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
				mapData:         map[string]interface{}{"name": "Ramses", "city": "Tiba"},
				validationRules: map[string][]string{"name": {"required", "kind:string"}, "city": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				mapData:         map[string]interface{}{"name": "Ramses", "age": 90},
				validationRules: map[string][]string{"name": {"required", "kind:int"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate nested map",
			args: args{
				mapData:         map[string]interface{}{"user": map[string]interface{}{"name": "Kyriakos M.", "country": "Egypt"}},
				validationRules: map[string][]string{"user": {"required"}, "user.name": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				mapData:         map[string]interface{}{"user": map[string]interface{}{"name": 1, "country": "Egypt"}},
				validationRules: map[string][]string{"user": {"required"}, "user.name": {"required", "kind:string"}, "user.country": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate nested map with non-string key",
			args: args{
				mapData:         map[string]interface{}{"user": map[int]interface{}{1: 2}},
				validationRules: map[string][]string{"user": {"required"}},
			},
			wantErr:                       true,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map includes struct",
			args: args{
				mapData:         map[string]interface{}{"user": User{Name: 5}},
				validationRules: map[string][]string{"user": {"required"}, "user.Name": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map includes struct with unsuitable data",
			args: args{
				mapData:         map[string]interface{}{"user": User{Name: 5}},
				validationRules: map[string][]string{"user": {"required"}, "user.Name": {"required", "kind:string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			err := v.validateMap(tt.args.mapData, "")
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
		fieldName       string
		fieldType       reflect.Type
		fieldValue      interface{}
		validationRules map[string][]string
		fieldsExists    map[string]bool
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
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(Parent{}),
				fieldValue:      Parent{Name: "Barssom"},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
				fieldsExists:    map[string]bool{"Parent.Name": true},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate by type - struct with unsuitable data",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(Parent{}),
				fieldValue:      Parent{},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
				fieldsExists:    map[string]bool{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate by type - map",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(map[string]interface{}{}),
				fieldValue:      map[string]interface{}{"Name": "Pola"},
				validationRules: map[string][]string{"Parent.Name": {"required"}},
				fieldsExists:    map[string]bool{"Parent.Name": true},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate by type - map with unsuitable data",
			args: args{
				fieldName:       "Parent",
				fieldType:       reflect.TypeOf(map[string]interface{}{}),
				fieldValue:      map[string]interface{}{"Name": 22},
				validationRules: map[string][]string{"Parent.Name": {"kind:string"}},
				fieldsExists:    map[string]bool{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate by type - field",
			args: args{
				fieldName:       "age",
				fieldType:       reflect.TypeOf(1),
				fieldValue:      1,
				validationRules: map[string][]string{"age": {"required"}},
				fieldsExists:    map[string]bool{"age": true},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate by type - field with unsuitable data",
			args: args{
				fieldName:       "age",
				fieldType:       reflect.TypeOf(1),
				fieldValue:      0,
				validationRules: map[string][]string{"age": {"required"}},
				fieldsExists:    map[string]bool{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			err := v.validateByType(tt.args.fieldName, tt.args.fieldType, tt.args.fieldValue)
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
		t                reflect.Type
		v                reflect.Value
		parentName       string
		validationRules  map[string][]string
		validationErrors map[string]string
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
				t:                reflect.TypeOf(Parent{Name: "Sherry"}),
				v:                reflect.ValueOf(Parent{Name: "Sherry"}),
				parentName:       "",
				validationRules:  map[string][]string{"Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate struct fields with unsuitable data",
			args: args{
				t:                reflect.TypeOf(Parent{}),
				v:                reflect.ValueOf(Parent{}),
				parentName:       "",
				validationRules:  map[string][]string{"Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			err := v.validateStructFields(tt.args.t, tt.args.v, tt.args.parentName)
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
		mapData          map[string]interface{}
		parentName       string
		validationRules  map[string][]string
		validationErrors map[string]string
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
				mapData:          map[string]interface{}{"Name": "Pola"},
				parentName:       "",
				validationRules:  map[string][]string{"Parent.Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map fields with unsuitable data",
			args: args{
				mapData:          map[string]interface{}{"Name": ""},
				parentName:       "Parent",
				validationRules:  map[string][]string{"Parent.Name": {"required"}},
				validationErrors: map[string]string{},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			err := v.validateMapFields(tt.args.mapData, tt.args.parentName)
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
		validationRules  map[string][]string
		fieldsExists     map[string]bool
		validationErrors map[string]string
	}
	tests := []struct {
		name                          string
		args                          args
		validationErrorsCountExpected int
	}{
		{
			name: "test validate non exist required fields",
			args: args{
				validationRules:  map[string][]string{"field1": {"required", "string"}, "field3": {"string", "required"}},
				fieldsExists:     map[string]bool{"field1": true, "field2": true},
				validationErrors: map[string]string{},
			},
			validationErrorsCountExpected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createNewValidation(tt.args.validationRules)
			v.fieldsExists = tt.args.fieldsExists
			v.validateNonExistRequiredFields()
			if len(v.errors) != tt.validationErrorsCountExpected {
				t.Errorf("validateNonExistRequiredFields() errors = %v len = %v, want %v", v.errors, len(v.errors), tt.validationErrorsCountExpected)
			}
		})
	}
}
