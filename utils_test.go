package validation

import (
	"testing"
)

func Test_IsEmpty(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test check if empty value is empty",
			args: args{val: ""},
			want: true,
		},
		{
			name: "test check if non-empty value is empty",
			args: args{val: "t"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.val); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
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

func Test_IsUint(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is uint rule",
			args: args{val: uint(15)},
			want: true,
		},
		{
			name: "test is uint rule with signed int value",
			args: args{val: -15},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUint(tt.args.val); got != tt.want {
				t.Errorf("IsUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsUint8(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is uint8 rule",
			args: args{val: uint8(200)},
			want: true,
		},
		{
			name: "test is uint8 rule with non-uint8 value",
			args: args{val: uint16(65534)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUint8(tt.args.val); got != tt.want {
				t.Errorf("IsUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsUint16(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is uint16 rule",
			args: args{val: uint16(65534)},
			want: true,
		},
		{
			name: "test is uint16 rule with non-uint16 value",
			args: args{val: uint32(4294967294)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUint16(tt.args.val); got != tt.want {
				t.Errorf("IsUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsUint32(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is uint32 rule",
			args: args{val: uint32(4294967294)},
			want: true,
		},
		{
			name: "test is uint32 rule with non-uint32 value",
			args: args{val: uint64(18446744073709551614)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUint32(tt.args.val); got != tt.want {
				t.Errorf("IsUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsUint64(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is uint64 rule",
			args: args{val: uint64(18446744073709551614)},
			want: true,
		},
		{
			name: "test is uint64 rule with non-uint32 value",
			args: args{val: int64(9223372036854775806)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUint64(tt.args.val); got != tt.want {
				t.Errorf("IsUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsFloat32(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is float32 rule",
			args: args{val: float32(2.2)},
			want: true,
		},
		{
			name: "test is float32 rule with non-float32 value",
			args: args{val: 55},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFloat32(tt.args.val); got != tt.want {
				t.Errorf("IsFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsFloat64(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is float64 rule",
			args: args{val: 2.2},
			want: true,
		},
		{
			name: "test is float64 rule with non-float64 value",
			args: args{val: 55},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFloat64(tt.args.val); got != tt.want {
				t.Errorf("IsFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsComplex64(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is complex64 rule",
			args: args{val: complex64(2 + 2i)},
			want: true,
		},
		{
			name: "test is complex64 rule with non-complex64 value",
			args: args{val: 456456 + 456456i},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsComplex64(tt.args.val); got != tt.want {
				t.Errorf("IsComplex64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsComplex128(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is complex128 rule",
			args: args{val: 2 + 2i},
			want: true,
		},
		{
			name: "test is complex128 rule with non-complex128 value",
			args: args{val: complex64(456456 + 456456i)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsComplex128(tt.args.val); got != tt.want {
				t.Errorf("IsComplex128() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsBool(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is bool rule with true",
			args: args{val: true},
			want: true,
		},
		{
			name: "test is bool rule with false",
			args: args{val: false},
			want: true,
		},
		{
			name: "test is bool rule with non-bool value",
			args: args{val: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBool(tt.args.val); got != tt.want {
				t.Errorf("IsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsSlice(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is slice rule",
			args: args{val: []int{4, 2}},
			want: true,
		},
		{
			name: "test is slice rule with non-slice value",
			args: args{val: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSlice(tt.args.val); got != tt.want {
				t.Errorf("IsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsArray(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is array rule",
			args: args{val: [2]int{4, 2}},
			want: true,
		},
		{
			name: "test is array rule with non-array value",
			args: args{val: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsArray(tt.args.val); got != tt.want {
				t.Errorf("IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsStruct(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is struct rule",
			args: args{val: struct{}{}},
			want: true,
		},
		{
			name: "test is struct rule with non-struct value",
			args: args{val: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStruct(tt.args.val); got != tt.want {
				t.Errorf("IsStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsMap(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test is map rule",
			args: args{val: make(map[interface{}]interface{})},
			want: true,
		},
		{
			name: "test is map rule with non-map value",
			args: args{val: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMap(tt.args.val); got != tt.want {
				t.Errorf("IsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsInteger(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsInteger with uint",
			args: args{
				val: uint(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with int",
			args: args{
				val: 1,
			},
			want: true,
		},
		{
			name: "test IsInteger with uint8",
			args: args{
				val: uint8(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with int8",
			args: args{
				val: int8(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with uint16",
			args: args{
				val: uint16(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with int16",
			args: args{
				val: int16(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with uint32",
			args: args{
				val: uint32(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with int32",
			args: args{
				val: int32(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with uint64",
			args: args{
				val: uint64(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with int64",
			args: args{
				val: int64(1),
			},
			want: true,
		},
		{
			name: "test IsInteger with float32",
			args: args{
				val: float32(1.1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInteger(tt.args.val); got != tt.want {
				t.Errorf("IsInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}
