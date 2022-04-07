package valdn

import (
	"mime/multipart"
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

func Test_IsKind(t *testing.T) {
	type args struct {
		val  interface{}
		kind string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsKind with string",
			args: args{
				val:  "string",
				kind: "string",
			},
			want: true,
		},
		{
			name: "test IsKind with int",
			args: args{
				val:  -44,
				kind: "int",
			},
			want: true,
		},
		{
			name: "test IsKind with uint",
			args: args{
				val:  uint(44),
				kind: "uint",
			},
			want: true,
		},
		{
			name: "test IsKind with float",
			args: args{
				val:  44.44,
				kind: "float64",
			},
			want: true,
		},
		{
			name: "test IsKind with complex",
			args: args{
				val:  44 + 22i,
				kind: "complex128",
			},
			want: true,
		},
		{
			name: "test IsKind with bool",
			args: args{
				val:  true,
				kind: "bool",
			},
			want: true,
		},
		{
			name: "test IsKind with map",
			args: args{
				val:  map[string]interface{}{},
				kind: "map",
			},
			want: true,
		},
		{
			name: "test IsKind with array",
			args: args{
				val:  [1]string{"array"},
				kind: "array",
			},
			want: true,
		},
		{
			name: "test IsKind with slice",
			args: args{
				val:  []string{"slice"},
				kind: "slice",
			},
			want: true,
		},
		{
			name: "test IsKind with struct",
			args: args{
				val:  struct{}{},
				kind: "struct",
			},
			want: true,
		},
		{
			name: "test IsKind with with unsuitable data",
			args: args{
				val:  1,
				kind: "string",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKind(tt.args.val, tt.args.kind); got != tt.want {
				t.Errorf("IsKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsKindIn(t *testing.T) {
	type args struct {
		val   interface{}
		kinds []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsKindIn",
			args: args{
				val:   44,
				kinds: []string{"int8", "string", "int"},
			},
			want: true,
		},
		{
			name: "test IsKindIn with unsuitable data",
			args: args{
				val:   44,
				kinds: []string{"int8", "string", "int32"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKindIn(tt.args.val, tt.args.kinds); got != tt.want {
				t.Errorf("IsKindIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsType(t *testing.T) {
	type user struct {
		name string
	}
	type args struct {
		val interface{}
		typ string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsType with string",
			args: args{
				val: "string",
				typ: "string",
			},
			want: true,
		},
		{
			name: "test IsType with int",
			args: args{
				val: -44,
				typ: "int",
			},
			want: true,
		},
		{
			name: "test IsType with uint",
			args: args{
				val: uint(44),
				typ: "uint",
			},
			want: true,
		},
		{
			name: "test IsType with float",
			args: args{
				val: 44.44,
				typ: "float64",
			},
			want: true,
		},
		{
			name: "test IsType with complex",
			args: args{
				val: 44 + 22i,
				typ: "complex128",
			},
			want: true,
		},
		{
			name: "test IsType with bool",
			args: args{
				val: true,
				typ: "bool",
			},
			want: true,
		},
		{
			name: "test IsType with map",
			args: args{
				val: map[string]interface{}{},
				typ: "map[string]interface {}",
			},
			want: true,
		},
		{
			name: "test IsType with array",
			args: args{
				val: [1]string{"array"},
				typ: "[1]string",
			},
			want: true,
		},
		{
			name: "test IsType with slice",
			args: args{
				val: []string{"slice"},
				typ: "[]string",
			},
			want: true,
		},
		{
			name: "test IsType with struct",
			args: args{
				val: user{name: "test"},
				typ: "user",
			},
			want: true,
		},
		{
			name: "test IsType with with unsuitable data",
			args: args{
				val: 1,
				typ: "string",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsType(tt.args.val, tt.args.typ); got != tt.want {
				t.Errorf("IsType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsTypeIn(t *testing.T) {
	type args struct {
		val   interface{}
		types []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsTypeIn",
			args: args{
				val:   map[string]interface{}{},
				types: []string{"string", "int", "map[string]interface {}"},
			},
			want: true,
		},
		{
			name: "test IsTypeIn with unsuitable data",
			args: args{
				val:   map[string]interface{}{},
				types: []string{"string", "int", "int8"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTypeIn(tt.args.val, tt.args.types); got != tt.want {
				t.Errorf("IsTypeIn() = %v, want %v", got, tt.want)
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

func Test_IsUnsignedInteger(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsUnsignedInteger with uint",
			args: args{
				val: uint(1),
			},
			want: true,
		},
		{
			name: "test IsUnsignedInteger with uint8",
			args: args{
				val: uint8(1),
			},
			want: true,
		},
		{
			name: "test IsUnsignedInteger with uint16",
			args: args{
				val: uint16(1),
			},
			want: true,
		},
		{
			name: "test IsUnsignedInteger with uint32",
			args: args{
				val: uint32(1),
			},
			want: true,
		},
		{
			name: "test IsUnsignedInteger with uint64",
			args: args{
				val: uint64(1),
			},
			want: true,
		},
		{
			name: "test IsUnsignedInteger with singed int",
			args: args{
				val: -1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnsignedInteger(tt.args.val); got != tt.want {
				t.Errorf("IsUnsignedInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsFloat(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsFloat with float32",
			args: args{
				val: float32(1.1),
			},
			want: true,
		},
		{
			name: "test IsFloat with float64",
			args: args{
				val: 1.1,
			},
			want: true,
		},
		{
			name: "test IsFloat with int",
			args: args{
				val: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFloat(tt.args.val); got != tt.want {
				t.Errorf("IsFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsUnsignedFloat(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsUnsignedFloat",
			args: args{
				val: 4.4,
			},
			want: true,
		},
		{
			name: "test IsUnsignedFloat with unsuitable data",
			args: args{
				val: -1.2,
			},
			want: false,
		},
		{
			name: "test IsUnsignedFloat with non-float value",
			args: args{
				val: 55,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnsignedFloat(tt.args.val); got != tt.want {
				t.Errorf("IsUnsignedFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsComplex(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsComplex with complex64",
			args: args{
				val: complex64(456456 + 456456i),
			},
			want: true,
		},
		{
			name: "test IsComplex with complex128",
			args: args{
				val: 456456 + 456456i,
			},
			want: true,
		},
		{
			name: "test IsComplex with int",
			args: args{
				val: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsComplex(tt.args.val); got != tt.want {
				t.Errorf("IsComplex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsNumeric(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsNumeric with integer",
			args: args{
				val: 1,
			},
			want: true,
		},
		{
			name: "test IsNumeric with float",
			args: args{
				val: 1.1,
			},
			want: true,
		},
		{
			name: "test IsNumeric with complex",
			args: args{
				val: 1 + 2i,
			},
			want: true,
		},
		{
			name: "test IsNumeric with string",
			args: args{
				val: "ss",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumeric(tt.args.val); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsCollection(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsCollection with slice",
			args: args{
				val: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "test IsCollection with array",
			args: args{
				val: [3]int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "test IsCollection with map",
			args: args{
				val: map[string]int{"one": 1},
			},
			want: true,
		},
		{
			name: "test IsCollection with struct",
			args: args{
				val: args{val: nil},
			},
			want: true,
		},
		{
			name: "test IsCollection with string",
			args: args{
				val: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCollection(tt.args.val); got != tt.want {
				t.Errorf("IsCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsEmail(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsEmail",
			args: args{
				s: "email@email.com",
			},
			want: true,
		},
		{
			name: "test IsEmail with invalid email",
			args: args{
				s: "emailemail.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.args.s); got != tt.want {
				t.Errorf("IsEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsJSON(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsJSON",
			args: args{
				s: `{"name":"Ramses", "city":"Tiba"}`,
			},
			want: true,
		},
		{
			name: "test IsJSON with unsuitable data",
			args: args{
				s: `name:"Ramses`,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSON(tt.args.s); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsIPv4(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsIPv4",
			args: args{
				s: "255.255.255.255",
			},
			want: true,
		},
		{
			name: "test IsIPv4 with unsuitable data",
			args: args{
				s: "1.1.1.1.",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv4(tt.args.s); got != tt.want {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsIPv6(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "test IsIPv6",
			args: args{
				s: "2001:db8:3333:4444:5555:6666:1.2.3.4",
			},
			want: true,
		},
		{
			name: "test IsIPv6 with unsuitable data",
			args: args{
				s: "1.1.1.1.",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv6(tt.args.s); got != tt.want {
				t.Errorf("IsIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsIP(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsIP",
			args: args{
				s: "2001:db8:3333:4444:5555:6666:1.2.3.4",
			},
			want: true,
		},
		{
			name: "test IsIP with unsuitabel data",
			args: args{
				s: "6489",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIP(tt.args.s); got != tt.want {
				t.Errorf("IsIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsMAC(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsMAC",
			args: args{
				s: "3D:F2:C9:A6:B3:4F",
			},
			want: true,
		},
		{
			name: "test IsMAC with unsuitable data",
			args: args{
				s: "asfgasg.asgas.asg456",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMAC(tt.args.s); got != tt.want {
				t.Errorf("IsMAC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsURL",
			args: args{
				s: "presidency.eg",
			},
			want: true,
		},
		{
			name: "test IsURL with unsuitable ata",
			args: args{
				s: "presidency_eg",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.args.s); got != tt.want {
				t.Errorf("IsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsFile(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test IsFile",
			args: args{
				v: &multipart.FileHeader{Size: 44},
			},
			want: true,
		},
		{
			name: "test IsFile with unsuitable data",
			args: args{
				v: "bla bla",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.args.v); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
