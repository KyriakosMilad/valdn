package validation

import (
	"testing"
)

func TestValidateJson(t *testing.T) {
	type args struct {
		jsonData        string
		validationRules map[string][]string
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		wantValidationErrors bool
	}{
		{
			name: `test validate json`,
			args: args{
				jsonData:        `{"name":"Ramses", "city":"Tiba"}`,
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}},
			},
			wantErr:              false,
			wantValidationErrors: false,
		},
		{
			name: `test validate json with unsuitable data`,
			args: args{
				jsonData:        `{"name":"Ramses", "age":90}`,
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}, "age": {"required", "numerical"}},
			},
			wantErr:              true,
			wantValidationErrors: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateJson(tt.args.jsonData, tt.args.validationRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v, args %v", err, tt.wantErr, tt.args)
			}
			if (len(validationErrors) > 0) != tt.wantValidationErrors {
				t.Errorf("Validate() error = %v, wantValidationErrors %v, args %v", validationErrors, tt.wantErr, tt.args)
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
		name                 string
		args                 args
		wantErr              bool
		wantValidationErrors bool
	}{
		{
			name: "test validate map",
			args: args{
				mapData:         map[string]interface{}{"name": "Ramses", "city": "Tiba"},
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}},
			},
			wantErr:              false,
			wantValidationErrors: false,
		},
		{
			name: "test validate map with unsuitable data",
			args: args{
				mapData:         map[string]interface{}{"name": "Ramses", "age": 90},
				validationRules: map[string][]string{"name": {"required", "string"}, "city": {"required", "string"}, "age": {"required", "numerical"}},
			},
			wantErr:              false,
			wantValidationErrors: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, validationErrors := ValidateMap(tt.args.mapData, tt.args.validationRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v, args %v", err, tt.wantErr, tt.args)
			}
			if (len(validationErrors) > 0) != tt.wantValidationErrors {
				t.Errorf("Validate() error = %v, wantValidationErrors %v, args %v", validationErrors, tt.wantErr, tt.args)
			}
		})
	}
}
