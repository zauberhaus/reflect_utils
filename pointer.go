// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils

import "reflect"

func Ptr[T any](v T) *T {
	return &v
}

func IsPointer(value any) bool {
	if v, ok := value.(reflect.Value); ok {
		return v.Kind() == reflect.Pointer
	} else if v, ok := value.(reflect.Type); ok {
		return v.Kind() == reflect.Pointer
	} else {
		v := reflect.ValueOf(value)
		return v.Kind() == reflect.Pointer
	}
}
