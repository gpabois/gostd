// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	iter "github.com/gpabois/gostd/iter"
	decoder "github.com/gpabois/gostd/serde/decoder"

	mock "github.com/stretchr/testify/mock"

	reflect "reflect"

	result "github.com/gpabois/gostd/result"
)

// Decoder is an autogenerated mock type for the Decoder type
type Decoder struct {
	mock.Mock
}

type Decoder_Expecter struct {
	mock *mock.Mock
}

func (_m *Decoder) EXPECT() *Decoder_Expecter {
	return &Decoder_Expecter{mock: &_m.Mock}
}

// DecodePrimaryType provides a mock function with given fields: data, typ
func (_m *Decoder) DecodePrimaryType(data interface{}, typ reflect.Type) result.Result[reflect.Value] {
	ret := _m.Called(data, typ)

	var r0 result.Result[reflect.Value]
	if rf, ok := ret.Get(0).(func(interface{}, reflect.Type) result.Result[reflect.Value]); ok {
		r0 = rf(data, typ)
	} else {
		r0 = ret.Get(0).(result.Result[reflect.Value])
	}

	return r0
}

// Decoder_DecodePrimaryType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecodePrimaryType'
type Decoder_DecodePrimaryType_Call struct {
	*mock.Call
}

// DecodePrimaryType is a helper method to define mock.On call
//   - data interface{}
//   - typ reflect.Type
func (_e *Decoder_Expecter) DecodePrimaryType(data interface{}, typ interface{}) *Decoder_DecodePrimaryType_Call {
	return &Decoder_DecodePrimaryType_Call{Call: _e.mock.On("DecodePrimaryType", data, typ)}
}

func (_c *Decoder_DecodePrimaryType_Call) Run(run func(data interface{}, typ reflect.Type)) *Decoder_DecodePrimaryType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(reflect.Type))
	})
	return _c
}

func (_c *Decoder_DecodePrimaryType_Call) Return(_a0 result.Result[reflect.Value]) *Decoder_DecodePrimaryType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Decoder_DecodePrimaryType_Call) RunAndReturn(run func(interface{}, reflect.Type) result.Result[reflect.Value]) *Decoder_DecodePrimaryType_Call {
	_c.Call.Return(run)
	return _c
}

// DecodeTime provides a mock function with given fields: data, typ
func (_m *Decoder) DecodeTime(data interface{}, typ reflect.Type) result.Result[reflect.Value] {
	ret := _m.Called(data, typ)

	var r0 result.Result[reflect.Value]
	if rf, ok := ret.Get(0).(func(interface{}, reflect.Type) result.Result[reflect.Value]); ok {
		r0 = rf(data, typ)
	} else {
		r0 = ret.Get(0).(result.Result[reflect.Value])
	}

	return r0
}

// Decoder_DecodeTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecodeTime'
type Decoder_DecodeTime_Call struct {
	*mock.Call
}

// DecodeTime is a helper method to define mock.On call
//   - data interface{}
//   - typ reflect.Type
func (_e *Decoder_Expecter) DecodeTime(data interface{}, typ interface{}) *Decoder_DecodeTime_Call {
	return &Decoder_DecodeTime_Call{Call: _e.mock.On("DecodeTime", data, typ)}
}

func (_c *Decoder_DecodeTime_Call) Run(run func(data interface{}, typ reflect.Type)) *Decoder_DecodeTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(reflect.Type))
	})
	return _c
}

func (_c *Decoder_DecodeTime_Call) Return(_a0 result.Result[reflect.Value]) *Decoder_DecodeTime_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Decoder_DecodeTime_Call) RunAndReturn(run func(interface{}, reflect.Type) result.Result[reflect.Value]) *Decoder_DecodeTime_Call {
	_c.Call.Return(run)
	return _c
}

// GetCursor provides a mock function with given fields:
func (_m *Decoder) GetCursor() result.Result[interface{}] {
	ret := _m.Called()

	var r0 result.Result[interface{}]
	if rf, ok := ret.Get(0).(func() result.Result[interface{}]); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(result.Result[interface{}])
	}

	return r0
}

// Decoder_GetCursor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCursor'
type Decoder_GetCursor_Call struct {
	*mock.Call
}

// GetCursor is a helper method to define mock.On call
func (_e *Decoder_Expecter) GetCursor() *Decoder_GetCursor_Call {
	return &Decoder_GetCursor_Call{Call: _e.mock.On("GetCursor")}
}

func (_c *Decoder_GetCursor_Call) Run(run func()) *Decoder_GetCursor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Decoder_GetCursor_Call) Return(_a0 result.Result[interface{}]) *Decoder_GetCursor_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Decoder_GetCursor_Call) RunAndReturn(run func() result.Result[interface{}]) *Decoder_GetCursor_Call {
	_c.Call.Return(run)
	return _c
}

// IsNull provides a mock function with given fields: data
func (_m *Decoder) IsNull(data interface{}) bool {
	ret := _m.Called(data)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Decoder_IsNull_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsNull'
type Decoder_IsNull_Call struct {
	*mock.Call
}

// IsNull is a helper method to define mock.On call
//   - data interface{}
func (_e *Decoder_Expecter) IsNull(data interface{}) *Decoder_IsNull_Call {
	return &Decoder_IsNull_Call{Call: _e.mock.On("IsNull", data)}
}

func (_c *Decoder_IsNull_Call) Run(run func(data interface{})) *Decoder_IsNull_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Decoder_IsNull_Call) Return(_a0 bool) *Decoder_IsNull_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Decoder_IsNull_Call) RunAndReturn(run func(interface{}) bool) *Decoder_IsNull_Call {
	_c.Call.Return(run)
	return _c
}

// IterMap provides a mock function with given fields: data
func (_m *Decoder) IterMap(data interface{}) result.Result[iter.Iterator[decoder.Element]] {
	ret := _m.Called(data)

	var r0 result.Result[iter.Iterator[decoder.Element]]
	if rf, ok := ret.Get(0).(func(interface{}) result.Result[iter.Iterator[decoder.Element]]); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(result.Result[iter.Iterator[decoder.Element]])
	}

	return r0
}

// Decoder_IterMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IterMap'
type Decoder_IterMap_Call struct {
	*mock.Call
}

// IterMap is a helper method to define mock.On call
//   - data interface{}
func (_e *Decoder_Expecter) IterMap(data interface{}) *Decoder_IterMap_Call {
	return &Decoder_IterMap_Call{Call: _e.mock.On("IterMap", data)}
}

func (_c *Decoder_IterMap_Call) Run(run func(data interface{})) *Decoder_IterMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Decoder_IterMap_Call) Return(_a0 result.Result[iter.Iterator[decoder.Element]]) *Decoder_IterMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Decoder_IterMap_Call) RunAndReturn(run func(interface{}) result.Result[iter.Iterator[decoder.Element]]) *Decoder_IterMap_Call {
	_c.Call.Return(run)
	return _c
}

// IterSlice provides a mock function with given fields: data
func (_m *Decoder) IterSlice(data interface{}) result.Result[iter.Iterator[interface{}]] {
	ret := _m.Called(data)

	var r0 result.Result[iter.Iterator[interface{}]]
	if rf, ok := ret.Get(0).(func(interface{}) result.Result[iter.Iterator[interface{}]]); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(result.Result[iter.Iterator[interface{}]])
	}

	return r0
}

// Decoder_IterSlice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IterSlice'
type Decoder_IterSlice_Call struct {
	*mock.Call
}

// IterSlice is a helper method to define mock.On call
//   - data interface{}
func (_e *Decoder_Expecter) IterSlice(data interface{}) *Decoder_IterSlice_Call {
	return &Decoder_IterSlice_Call{Call: _e.mock.On("IterSlice", data)}
}

func (_c *Decoder_IterSlice_Call) Run(run func(data interface{})) *Decoder_IterSlice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Decoder_IterSlice_Call) Return(_a0 result.Result[iter.Iterator[interface{}]]) *Decoder_IterSlice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Decoder_IterSlice_Call) RunAndReturn(run func(interface{}) result.Result[iter.Iterator[interface{}]]) *Decoder_IterSlice_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewDecoder interface {
	mock.TestingT
	Cleanup(func())
}

// NewDecoder creates a new instance of Decoder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDecoder(t mockConstructorTestingTNewDecoder) *Decoder {
	mock := &Decoder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
