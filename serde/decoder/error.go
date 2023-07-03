package decoder

import (
	"fmt"
	"reflect"
)

type ExpectingPrimaryTypeError struct {
	Type reflect.Type
}

func (err ExpectingPrimaryTypeError) Error() string {
	return "value type is not a primary type"
}

func NewExpectingPrimaryTypeError(typ reflect.Type) error {
	return ExpectingPrimaryTypeError{typ}
}

type ExpectingSliceError struct {
	Type reflect.Type
}

func (err ExpectingSliceError) Error() string {
	return "value type is not a slice"
}

func NewExpectingSliceError(typ reflect.Type) error {
	return ExpectingMapError{typ}
}

type ExpectingMapError struct {
	Type reflect.Type
}

func (err ExpectingMapError) Error() string {
	return "value type is not a map"
}

func NewExpectingMapError(typ reflect.Type) error {
	return ExpectingMapError{typ}
}

type WrongTypeError struct {
	ExpectedType reflect.Type
	Type         reflect.Type
}

func (err WrongTypeError) Error() string {
	return fmt.Sprintf("unexpected value type, expecting %s, got %s", err.ExpectedType, err.Type)
}

func NewWrongTypeError(typ reflect.Type, expectedTyp reflect.Type) error {
	return WrongTypeError{ExpectedType: expectedTyp, Type: typ}
}
