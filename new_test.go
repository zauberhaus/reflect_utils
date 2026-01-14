// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils_test

import (
	"context"
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
		{
			name:         "context.Context",
			typ:          reflect.TypeOf((*context.Context)(nil)).Elem(),
			expectedVal:  context.Background(),
			checkDefault: true,
		},
		{
			name:      "any",
			typ:       reflect.TypeOf((*any)(nil)).Elem(),
			expectNil: true,
		},
		{
			name:      "generic interface",
			typ:       reflect.TypeOf((*error)(nil)).Elem(),
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

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("basic type", func(t *testing.T) {
		i, err := utils.New[int]()
		assert.NoError(t, err)
		assert.Equal(t, 0, i)
	})

	t.Run("struct", func(t *testing.T) {
		type S struct{ Name string }
		s, err := utils.New[S]()
		assert.NoError(t, err)
		assert.Equal(t, S{}, s)
	})

	t.Run("pointer", func(t *testing.T) {
		p, err := utils.New[*int]()
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, 0, *p)
	})
}
