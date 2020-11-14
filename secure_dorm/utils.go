package secure_dorm

import (
	"reflect"
)

/*
 * Given a pointer to a slice of structs, returns a pointer to a new slice
 * of the same type.
 */
func NewSliceFromSlice(result interface{}) interface{} {
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		panic("CreateNewSliceOfSameType's argument must be of a slice of model structs")
	}

	if reflect.TypeOf(result).Elem().Kind() != reflect.Slice {
		panic("CreateNewSliceOfSameType's argument must be of a slice of model structs")
	}

	if reflect.TypeOf(result).Elem().Elem().Kind() != reflect.Struct {
		panic("CreateNewSliceOfSameType's argument must be of a slice of model structs")
	}

	return reflect.New(reflect.TypeOf(result).Elem()).Interface()
}

/*
 * Given a pointer to a struct, returns a pointer to a new slice of
 * structs of the same type.
 */
func NewSliceFromStruct(result interface{}) interface{} {
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		panic("First's argument must be of a pointer to a model struct")
	}

	if reflect.TypeOf(result).Elem().Kind() != reflect.Struct {
		panic("First's argument must be of a pointer to a model struct")
	}

	elemType := reflect.TypeOf(result).Elem()
	return reflect.New(reflect.SliceOf(elemType)).Interface()
}
