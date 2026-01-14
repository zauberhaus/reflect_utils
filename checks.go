// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils

import (
	"reflect"
)

func HasMethod(value any, method string) bool {
	if value == (any)(nil) {
		return false
	}

	if method == "" {
		return false
	}

	var t reflect.Type

	if v, ok := value.(reflect.Type); ok {
		t = v
	} else if v, ok := value.(reflect.Value); ok {
		t = v.Type()
	} else {
		t = reflect.TypeOf(value)
	}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.Name == method {
			return true
		}
	}

	return false
}

func HasField(value any, field string) bool {
	if value == (any)(nil) {
		return false
	}

	if v, ok := value.(reflect.Value); ok {
		for v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return false
		}

		if _, ok := v.Type().FieldByName(field); !ok {
			return false
		}

		return true

	} else if t, ok := value.(reflect.Type); ok {
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			return false
		}

		if _, ok := t.FieldByName(field); !ok {
			return false
		}

		return true
	} else {
		t := reflect.TypeOf(value)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			return false
		}

		if _, ok := t.FieldByName(field); !ok {
			return false
		}

		return true
	}
}

func IsStruct(value any) bool {
	if value == (any)(nil) {
		return false
	}

	if v, ok := value.(reflect.Value); ok {
		for v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		return v.Kind() == reflect.Struct
	} else if v, ok := value.(reflect.Type); ok {
		for v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		return v.Kind() == reflect.Struct
	} else {
		v := reflect.TypeOf(value)
		for v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		return v.Kind() == reflect.Struct
	}
}

func IsEnum(value any) bool {
	if value == (any)(nil) {
		return false
	}

	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	name := t.Name()

	if name != "" {
		return HasMethod(value, "IsA"+name)
	}

	return false
}

func IsNil(value any) bool {
	if value == nil {
		return true
	}

	if v, ok := value.(reflect.Value); ok {
		switch v.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
			return v.IsNil()
		}
	} else {
		v = reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
			return v.IsNil()
		}
	}

	return false
}

func IsEmpty(object any) bool {

	// get nil case out of the way
	if object == (any)(nil) {
		return true
	}

	var objValue reflect.Value
	if val, ok := object.(reflect.Value); ok {
		objValue = val
	} else {
		objValue = reflect.ValueOf(object)
	}

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return IsEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		return objValue.IsZero()
	}
}

func IsComparable(object any) bool {
	var objValue reflect.Value
	if val, ok := object.(reflect.Value); ok {
		objValue = val
	} else {
		objValue = reflect.ValueOf(object)
	}

	if objValue.Comparable() {
		return true
	}

	return false
}
