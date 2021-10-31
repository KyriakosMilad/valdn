package validation

import "testing"

func Test_IsZero(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test check if zero value is zero",
			args: args{val: ""},
			want: true,
		},
		{
			name: "test check if non-zero value is zero",
			args: args{val: "t"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsZero(tt.args.val); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsString(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test check if string is string",
			args: args{val: "s"},
			want: true,
		},
		{
			name: "test check if non-string is string",
			args: args{val: 5},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsString(tt.args.val); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsInt(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is int rule",
			args: args{val: 1},
			want: true,
		},
		{
			name: "test is int rule with non-int value",
			args: args{val: "1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt(tt.args.val); got != tt.want {
				t.Errorf("IsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsInt8(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is int8 rule",
			args: args{val: int8(5)},
			want: true,
		},
		{
			name: "test is int8 rule with non-int8 value",
			args: args{val: int16(200)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt8(tt.args.val); got != tt.want {
				t.Errorf("IsInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsInt16(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is int16 rule",
			args: args{val: int16(200)},
			want: true,
		},
		{
			name: "test is int16 rule with non-int16 value",
			args: args{val: int32(2147483646)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt16(tt.args.val); got != tt.want {
				t.Errorf("IsInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsInt32(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is int32 rule",
			args: args{val: int32(2147483646)},
			want: true,
		},
		{
			name: "test is int32 rule with non-int32 value",
			args: args{val: int64(9223372036854775806)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt32(tt.args.val); got != tt.want {
				t.Errorf("IsInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsInt64(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is int64 rule",
			args: args{val: int64(9223372036854775806)},
			want: true,
		},
		{
			name: "test is int64 rule with non-int64 value",
			args: args{val: int8(2)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt64(tt.args.val); got != tt.want {
				t.Errorf("IsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
