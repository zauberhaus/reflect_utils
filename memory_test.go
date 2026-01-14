// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	utils "github.com/zauberhaus/reflect_utils"
)

func TestCopyToHeap(t *testing.T) {
	t.Parallel()

	// Value type
	i := 10
	res := utils.CopyToHeap(i)
	assert.IsType(t, &i, res)
	assert.Equal(t, i, *(res.(*int)))

	// Pointer type (should return as is)
	p := &i
	resP := utils.CopyToHeap(p)
	assert.Equal(t, p, resP)

	// Struct
	type S struct{ A int }
	s := S{A: 1}
	resS := utils.CopyToHeap(s)
	assert.IsType(t, &s, resS)
	assert.Equal(t, s, *(resS.(*S)))
}

func TestFromPointer(t *testing.T) {
	t.Parallel()

	// Pointer
	i := 10
	p := &i
	res := utils.FromPointer(p)
	assert.Equal(t, i, res)

	// Value (should return as is)
	resV := utils.FromPointer(i)
	assert.Equal(t, i, resV)
}
