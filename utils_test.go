// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	utils "github.com/zauberhaus/reflect_utils"
)

func TestNewOf(t *testing.T) {
	t.Parallel()

	type MyStruct struct {
		Name string `default:"test"`
	}

	tests := []struct {
		name         string
		typ          reflect.Type
		expectedVal  any
		expectNil    bool
		checkDefault bool
	}{
		{
			name:        "int",
			typ:         reflect.TypeOf(0),
			expectedVal: 0,
		},
		{
			name:        "pointer to int",
			typ:         reflect.TypeOf(new(int)),
			expectedVal: new(int),
		},
		{
			name:        "string",
			typ:         reflect.TypeOf(""),
			expectedVal: "",
		},
		{
			name:        "pointer to string",
			typ:         reflect.TypeOf(new(string)),
			expectedVal: new(string),
		},
		{
			name:        "slice",
			typ:         reflect.TypeOf([]int{}),
			expectedVal: []int{},
		},
		{
			name:        "map",
			typ:         reflect.TypeOf(map[string]string{}),
			expectedVal: map[string]string{},
		},
		{
			name:         "struct with default",
			typ:          reflect.TypeOf(MyStruct{}),
			expectedVal:  MyStruct{Name: "test"},
			checkDefault: true,
		},
		{
			name:         "pointer to struct with default",
			typ:          reflect.TypeOf(&MyStruct{}),
			expectedVal:  &MyStruct{Name: "test"},
			checkDefault: true,
		},
		{
			name:      "nil type",
			typ:       nil,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.NewWithDefaultsOf(tt.typ)
			require.NoError(t, err)

			if tt.expectNil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			if tt.checkDefault {
				assert.Equal(t, tt.expectedVal, result)
			} else {
				assert.IsType(t, tt.expectedVal, result)
				if reflect.TypeOf(tt.expectedVal).Kind() != reflect.Pointer {
					assert.Equal(t, tt.expectedVal, result)
				}
			}
		})
	}
}

func TestHasDefault(t *testing.T) {
	t.Parallel()

	type StructWithDefault struct {
		Name string `default:"hello"`
	}
	type StructWithoutDefault struct {
		Name string
	}
	type StructWithEmptyDefault struct {
		Name string `default:""`
	}
	type StructWithIgnoredDefault struct {
		Name string `default:"-"`
	}

	assert.True(t, utils.HasDefault(reflect.TypeOf(StructWithDefault{})))
	assert.True(t, utils.HasDefault(reflect.TypeOf(&StructWithDefault{})))
	assert.False(t, utils.HasDefault(reflect.TypeOf(StructWithoutDefault{})))
	assert.True(t, utils.HasDefault(reflect.TypeOf(StructWithEmptyDefault{})))
	assert.False(t, utils.HasDefault(reflect.TypeOf(StructWithIgnoredDefault{})))
	assert.False(t, utils.HasDefault(reflect.TypeOf(1)))
	assert.False(t, utils.HasDefault(nil))
}

func TestNewDefaultOf(t *testing.T) {
	t.Parallel()

	type MyStruct struct {
		Name    string `default:"test"`
		Number  int    `default:"123"`
		Pointer *int
	}

	t.Run("struct", func(t *testing.T) {
		result, err := utils.NewWithDefaultsOf(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)
		assert.Equal(t, MyStruct{Name: "test", Number: 123, Pointer: nil}, result)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		result, err := utils.NewWithDefaultsOf(reflect.TypeOf(&MyStruct{}))
		assert.NoError(t, err)
		assert.Equal(t, &MyStruct{Name: "test", Number: 123, Pointer: nil}, result)
	})

	t.Run("non-struct", func(t *testing.T) {
		result, err := utils.NewWithDefaultsOf(reflect.TypeOf(0))
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})
}

type MethodStruct struct{}

//lint:ignore U1000 MyMethod is used for reflection testing.
func (m MethodStruct) MyMethod() {}

//lint:ignore U1000 MyPointerMethod is used for reflection testing.
func (m *MethodStruct) MyPointerMethod() {}

func TestHasMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		value    any
		method   string
		expected bool
	}{
		{"value with value receiver", MethodStruct{}, "MyMethod", true},
		{"value with pointer receiver", MethodStruct{}, "MyPointerMethod", false},
		{"pointer with value receiver", &MethodStruct{}, "MyMethod", true},
		{"pointer with pointer receiver", &MethodStruct{}, "MyPointerMethod", true},
		{"non-existent method", MethodStruct{}, "NonExistent", false},
		{"nil value", nil, "MyMethod", false},
		{"empty method name", MethodStruct{}, "", false},
		{"reflect.Type with value receiver", reflect.TypeOf(MethodStruct{}), "MyMethod", true},
		{"reflect.Type with pointer receiver", reflect.TypeOf(MethodStruct{}), "MyPointerMethod", false},
		{"reflect.Type (ptr) with pointer receiver", reflect.TypeOf(&MethodStruct{}), "MyPointerMethod", true},
		{"reflect.Value with value receiver", reflect.ValueOf(MethodStruct{}), "MyMethod", true},
		{"reflect.Value (ptr) with pointer receiver", reflect.ValueOf(&MethodStruct{}), "MyPointerMethod", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, utils.HasMethod(tc.value, tc.method))
		})
	}
}

type MyEnum int

//lint:ignore U1000 IsAMyEnum is used for reflection testing.
func (e MyEnum) IsAMyEnum() {}

type NotAnEnum int

type PtrReceiverEnum string

//lint:ignore U1000 IsAPtrReceiverEnum is used for reflection testing.
func (e *PtrReceiverEnum) IsAPtrReceiverEnum() {}

func TestIsEnum(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		value    any
		expected bool
	}{
		{"valid enum with value receiver", MyEnum(1), true},
		{"not an enum", NotAnEnum(1), false},
		{"valid enum with pointer receiver (on pointer)", utils.Ptr(PtrReceiverEnum("a")), true},
		{"valid enum with pointer receiver (on value)", PtrReceiverEnum("a"), false},
		{"nil value", nil, false},
		{"primitive type int", 123, false},
		{"primitive type string", "hello", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, utils.IsEnum(tc.value))
		})
	}
}

func TestIsPointer(t *testing.T) {
	t.Parallel()
	var i int
	var p *int
	var nilP *int
	p = &i

	assert.True(t, utils.IsPointer(p))
	assert.True(t, utils.IsPointer(nilP))
	assert.False(t, utils.IsPointer(i))
	assert.False(t, utils.IsPointer(nil))

	assert.True(t, utils.IsPointer(reflect.ValueOf(p)))
	assert.False(t, utils.IsPointer(reflect.ValueOf(i)))

	assert.True(t, utils.IsPointer(reflect.TypeOf(p)))
	assert.False(t, utils.IsPointer(reflect.TypeOf(i)))
}

func TestIsStruct(t *testing.T) {
	t.Parallel()
	type s struct{}
	var st s
	var p *s = &s{}
	var nilP *s

	assert.True(t, utils.IsStruct(st))
	assert.True(t, utils.IsStruct(p))
	assert.True(t, utils.IsStruct(nilP))
	assert.False(t, utils.IsStruct(1))
	assert.False(t, utils.IsStruct(nil))

	assert.True(t, utils.IsStruct(reflect.ValueOf(st)))
	assert.True(t, utils.IsStruct(reflect.ValueOf(p)))
	assert.False(t, utils.IsStruct(reflect.ValueOf(1)))

	assert.True(t, utils.IsStruct(reflect.TypeOf(st)))
	assert.True(t, utils.IsStruct(reflect.TypeOf(p)))
	assert.False(t, utils.IsStruct(reflect.TypeOf(1)))
}

func TestIsNil(t *testing.T) {
	t.Parallel()
	var p *int
	var m map[string]int
	var s []int
	var c chan int
	i := 1
	pNonNil := &i

	assert.True(t, utils.IsNil(nil))
	assert.True(t, utils.IsNil(p))
	assert.True(t, utils.IsNil(m))
	assert.True(t, utils.IsNil(s))
	assert.True(t, utils.IsNil(c))

	assert.False(t, utils.IsNil(i))
	assert.False(t, utils.IsNil(pNonNil))
	assert.False(t, utils.IsNil(make(map[string]int)))
	assert.False(t, utils.IsNil(make([]int, 0)))
	assert.False(t, utils.IsNil(make(chan int)))

	assert.True(t, utils.IsNil(reflect.ValueOf(p)))
	assert.False(t, utils.IsNil(reflect.ValueOf(i)))
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}
	var p *int
	i := 0
	pZero := &i

	// True cases
	assert.True(t, utils.IsEmpty(nil), "nil should be empty")
	assert.True(t, utils.IsEmpty(""), "empty string should be empty")
	assert.True(t, utils.IsEmpty(0), "int 0 should be empty")
	assert.True(t, utils.IsEmpty(false), "false should be empty")
	assert.True(t, utils.IsEmpty([]string{}), "empty slice should be empty")
	assert.True(t, utils.IsEmpty(map[string]int{}), "empty map should be empty")
	assert.True(t, utils.IsEmpty(p), "nil pointer should be empty")
	assert.True(t, utils.IsEmpty(pZero), "pointer to zero value should be empty")
	assert.True(t, utils.IsEmpty([2]int{0, 0}), "zero value array should be empty")
	assert.True(t, utils.IsEmpty(struct{}{}), "empty struct should be empty")
	assert.True(t, utils.IsEmpty(make(chan int)), "empty channel should be empty")

	// False cases
	assert.False(t, utils.IsEmpty("hello"), "non-empty string should not be empty")
	assert.False(t, utils.IsEmpty(1), "int 1 should not be empty")
	assert.False(t, utils.IsEmpty(true), "true should not be empty")
	assert.False(t, utils.IsEmpty([]string{"a"}), "non-empty slice should not be empty")
	assert.False(t, utils.IsEmpty(map[string]int{"a": 1}), "non-empty map should not be empty")
	j := 1
	pNonZero := &j
	assert.False(t, utils.IsEmpty(pNonZero), "pointer to non-zero value should not be empty")
	assert.False(t, utils.IsEmpty([2]int{0, 1}), "non-zero value array should not be empty")
	assert.False(t, utils.IsEmpty(chWithValue), "channel with value should not be empty")

	// reflect.Value cases
	assert.True(t, utils.IsEmpty(reflect.ValueOf("")))
	assert.False(t, utils.IsEmpty(reflect.ValueOf("a")))
}
