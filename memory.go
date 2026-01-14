// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils

import "reflect"

func CopyToHeap(val any) any {
	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	if t.Kind() == reflect.Pointer {
		return val
	}

	x := reflect.New(t)
	x.Elem().Set(v)
	return x.Interface()
}

func FromPointer(val any) any {
	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	if t.Kind() != reflect.Pointer {
		return val
	}

	v = v.Elem()
	return v.Interface()
}
