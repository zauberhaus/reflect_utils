// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package reflect_utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
)

func New[T any](v ...T) (T, error) {
	t := reflect.TypeFor[T]()
	val, err := NewOf(t)
	if err != nil {
		return *new(T), err
	}

	if result, ok := val.(T); ok {
		return result, nil
	}

	return *new(T), nil
}

func NewOf(t reflect.Type) (any, error) {
	if t == nil {
		return nil, nil
	}

	switch t.Kind() {
	case reflect.Ptr:
		n := t.Elem()
		v := reflect.New(n)
		if v.CanInterface() {
			y := v.Interface()
			rv := reflect.ValueOf(&y).Elem()
			t := rv.Elem().Type().Elem()
			rv.Set(reflect.New(t))
			return y, nil
		}
		return nil, fmt.Errorf("can't interface %v", t)
	case reflect.Slice:
		v := reflect.MakeSlice(t, 0, 0)
		if v.CanInterface() {
			return v.Interface(), nil
		}
		return nil, fmt.Errorf("can't interface %v", t)
	case reflect.Map:
		v := reflect.MakeMap(t)
		if v.CanInterface() {
			return v.Interface(), nil
		}
		return nil, fmt.Errorf("can't interface %v", t)
	case reflect.Struct:
		val := reflect.New(t)
		if val.CanInterface() {
			return val.Interface(), nil
		}
		return nil, fmt.Errorf("can't interface %v", t)

	case reflect.Interface:
		name := fmt.Sprintf("%v", t)
		switch name {
		case "context.Context":
			return context.Background(), nil
		case "interface {}":
			return any(nil), nil
		default:
			val := reflect.New(t).Elem()
			if val.CanInterface() {
				return val.Interface(), nil
			}

			return nil, fmt.Errorf("can't interface %v", t)
		}
	default:
		val := reflect.New(t).Elem()
		if val.CanInterface() {
			return val.Interface(), nil
		}

		return nil, fmt.Errorf("can't interface %v", t)
	}
}

func HasDefault(t reflect.Type) bool {
	if t == nil {
		return false
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		if defaultVal, ok := t.Field(i).Tag.Lookup("default"); ok && defaultVal != "-" {
			return true
		}
	}

	return false
}

func NewWithDefaultsOf(t reflect.Type) (any, error) {

	val, err := NewOf(t)
	if err != nil {
		return nil, err
	}

	if IsStruct(t) {
		err := defaults.Set(val)
		if err != nil {
			return nil, fmt.Errorf("set defaults failed: %w", err)
		}

		if !IsPointer(t) {
			v := reflect.ValueOf(val)
			e := v.Elem()
			return e.Interface(), nil
		}
	}

	return val, nil
}
