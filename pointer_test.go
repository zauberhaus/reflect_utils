// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	utils "github.com/zauberhaus/reflect_utils"
)

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

func TestPtr(t *testing.T) {
	t.Parallel()

	v := 10
	p := utils.Ptr(v)
	assert.NotNil(t, p)
	assert.Equal(t, v, *p)

	s := "test"
	ps := utils.Ptr(s)
	assert.NotNil(t, ps)
	assert.Equal(t, s, *ps)
}
