package validation

import (
	"reflect"
	"testing"
)

func Test_initialize(t *testing.T) {
	type args struct {
		rules Rules
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test initialize",
			args: args{
				rules: Rules{"test": {"required"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialize(tt.args.rules)
			if !reflect.DeepEqual(tt.args.rules, validationRules) {
				t.Errorf("initialize(): error initialize not wroking")
			}
			if reflect.ValueOf(validationErrors).IsZero() {
				t.Errorf("initialize(): validationErrors is nil")
			}
			if reflect.ValueOf(fieldsExists).IsZero() {
				t.Errorf("initialize(): fieldsExists is nil")
			}
		})
	}
}

func Test_ValidateField(t *testing.T) {
	type args struct {
		fieldName  string
		fieldValue interface{}
		fieldRules []string
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		wantValidationErrors bool
	}{
		{
			name: "test validate field",
			args: args{
				fieldName:  "Name",
				fieldValue: "Kyria",
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:              false,
			wantValidationErrors: false,
		},
		{
			name: "test validate field with unsuitable data",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"required", "kind:string"},
			},
			wantErr:              false,
			wantValidationErrors: true,
		},
		{
			name: "test validate field with not exists rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"bla:bla"},
			},
			wantErr:              true,
			wantValidationErrors: false,
		},
		{
			name: "test validate field with empty rule",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{},
			},
			wantErr:              false,
			wantValidationErrors: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationError := ValidateField(tt.args.fieldName, tt.args.fieldValue, tt.args.fieldRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateField() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.wantValidationErrors, tt.args)
			}
			if (validationError != "") != tt.wantValidationErrors {
				t.Errorf("ValidateField() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.wantValidationErrors, tt.args)
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
				t.Errorf("ValidateJson() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("ValidateJson() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func Test_validateMap(t *testing.T) {
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
			initialize(tt.args.validationRules)
			err := validateMap(tt.args.mapData, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMap() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateMap() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func Test_validateStruct(t *testing.T) {
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
			initialize(tt.args.validationRules)
			addValidationTagRules(reflect.TypeOf(tt.args.structData), "")
			err := validateStruct(tt.args.structData, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStruct() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("validateStruct() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func TestValidate(t *testing.T) {
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
			err, validationErrors := Validate(tt.args.val, tt.args.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() err = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(len(validationErrors), tt.validationErrorsCount) {
				t.Errorf("Validate() validationErrors = %v, want %v", validationErrors, tt.validationErrorsCount)
			}
		})
	}
}
