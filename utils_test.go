package validation

import "testing"

func Test_getRuleValue(t *testing.T) {
	tests := []struct {
		name string
		rule string
		want string
	}{
		{
			name: "test get rule value from rule does have value",
			rule: "val:test",
			want: "test",
		},
		{
			name: "test get rule value from rule does not have value",
			rule: "val",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRuleValue(tt.rule); got != tt.want {
				t.Errorf("getRuleValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
