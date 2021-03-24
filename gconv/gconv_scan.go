// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"

	"github.com/gogf/gf/errors/gerror"
)

// Scan automatically calls MapToMap, MapToMaps, Struct or Structs function according to
// the type of parameter `pointer` to implement the converting.
// It calls function MapToMap if `pointer` is type of *map to do the converting.
// It calls function MapToMaps if `pointer` is type of *[]map/*[]*map to do the converting.
// It calls function Struct if `pointer` is type of *struct/**struct to do the converting.
// It calls function Structs if `pointer` is type of *[]struct/*[]*struct to do the converting.
func Scan(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var (
		pointerType = reflect.TypeOf(pointer)
		pointerKind = pointerType.Kind()
	)
	if pointerKind != reflect.Ptr {
		return gerror.Newf("params should be type of pointer, but got: %v", pointerKind)
	}
	var (
		pointerElem     = pointerType.Elem()
		pointerElemKind = pointerElem.Kind()
	)
	switch pointerElemKind {
	case reflect.Map:
		return MapToMap(params, pointer, mapping...)
	case reflect.Array, reflect.Slice:
		var (
			sliceElem     = pointerElem.Elem()
			sliceElemKind = sliceElem.Kind()
		)
		for sliceElemKind == reflect.Ptr {
			sliceElem = sliceElem.Elem()
			sliceElemKind = sliceElem.Kind()
		}
		if sliceElemKind == reflect.Map {
			return MapToMaps(params, pointer, mapping...)
		}
		return Structs(params, pointer, mapping...)
	default:
		return Struct(params, pointer, mapping...)
	}
}

// ScanDeep automatically calls StructDeep or StructsDeep function according to the type of
// parameter `pointer` to implement the converting..
// It calls function StructDeep if `pointer` is type of *struct/**struct to do the converting.
// It calls function StructsDeep if `pointer` is type of *[]struct/*[]*struct to do the converting.
// Deprecated, use Scan instead.
func ScanDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return gerror.Newf("params should be type of pointer, but got: %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Array, reflect.Slice:
		return StructsDeep(params, pointer, mapping...)
	default:
		return StructDeep(params, pointer, mapping...)
	}
}
