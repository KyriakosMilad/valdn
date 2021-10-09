package validation

import (
	"testing"
)

func TestValidateField(t *testing.T) {
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
				fieldRules: []string{"required", "string"},
			},
			wantErr:              false,
			wantValidationErrors: false,
		},
		{
			name: "test validate field with unsuitable data",
			args: args{
				fieldName:  "Name",
				fieldValue: 55,
				fieldRules: []string{"required", "string"},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateField(tt.args.fieldName, tt.args.fieldValue, tt.args.fieldRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateField() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.wantValidationErrors, tt.args)
			}
			if (len(validationErrors) > 0) != tt.wantValidationErrors {
				t.Errorf("ValidateField() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.wantValidationErrors, tt.args)
			}
		})
	}
}

func TestValidateJson(t *testing.T) {
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
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: `test validate json with unsuitable data`,
			args: args{
				jsonData:        `{"name":"Ramses", "age":90}`,
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
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

func TestValidateMap(t *testing.T) {
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
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				mapData:         map[string]interface{}{"name": "Ramses", "age": 90},
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "test validate nested map",
			args: args{
				mapData:         map[string]interface{}{"user": map[string]interface{}{"name": "Kyriakos M.", "country": "Egypt"}},
				validationRules: map[string][]string{"user": {"required"}, "user.name": {"required", "string"}, "user.country": {"required", "string"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "test validate nested map with unsuitable data",
			args: args{
				mapData:         map[string]interface{}{"user": map[string]interface{}{"name": 1, "country": "Egypt"}},
				validationRules: map[string][]string{"user": {"required"}, "user.name": {"required", "string"}, "user.country": {"required", "string"}},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateMap(tt.args.mapData, tt.args.validationRules, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMap() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("ValidateMap() error = %v, validationErrors = %v, wantErr %v, wantValidationErrors %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	type Child struct {
		Name string `validation:"required|string"`
	}
	type Parent struct {
		Name            string `validation:"required|string"`
		Age             int    `validation:"required"`
		StringKeyMap    map[string]interface{}
		NonStringKeyMap map[int]interface{}
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
				validationRules: map[string][]string{"Name": {"required", "string"}, "Child.Name": {""}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 0,
		},
		{
			name: "validate struct with unsuitable data",
			args: args{
				structData:      Parent{Name: "Mina"},
				validationRules: map[string][]string{"Name": {"required", "string"}, "Age": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 1,
		},
		{
			name: "validate nested struct",
			args: args{
				structData: Parent{
					Name:  "Ikhnaton",
					Child: Child{Name: "Tut"},
				},
				validationRules: map[string][]string{"Name": {"required", "string"}, "Age": {""}, "Child.Name": {"required", "string"}},
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
				validationRules: map[string][]string{"Name": {"required", "string"}, "Child.Name": {"required", "string"}},
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
			expectedValidationErrorsCount: 1,
		},
		{
			name: "validate struct includes string-key-map",
			args: args{
				structData: Parent{
					Name: "Ikhnaton",
					Age:  2,
				},
				validationRules: map[string][]string{"stringKeyMap": {"required"}, "stringKeyMap.val": {"required"}},
			},
			wantErr:                       false,
			expectedValidationErrorsCount: 2,
		},
		{
			name: "validate struct includes non-string-key-map",
			args: args{
				structData: Parent{
					Name:            "Ikhnaton",
					Age:             2,
					NonStringKeyMap: map[int]interface{}{1: 2},
				},
				validationRules: map[string][]string{"NonStringKeyMap": {"required"}},
			},
			wantErr:                       true,
			expectedValidationErrorsCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateStruct(tt.args.structData, tt.args.validationRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
			if len(validationErrors) != tt.expectedValidationErrorsCount {
				t.Errorf("ValidateStruct() error = %v, validationErrors = %v, wantErr %v, expectedValidationErrorsCount %v, args %v", err, validationErrors, tt.wantErr, tt.expectedValidationErrorsCount, tt.args)
			}
		})
	}
}