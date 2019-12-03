// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import virtualmachine "github.com/vterdunov/janna/internal/virtualmachine"

// TaskStatuser is an autogenerated mock type for the TaskStatuser type
type TaskStatuser struct {
	mock.Mock
}

// Get provides a mock function with given fields:
func (_m *TaskStatuser) Get() map[string]interface{} {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *TaskStatuser) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Str provides a mock function with given fields: keyvals
func (_m *TaskStatuser) Str(keyvals ...string) virtualmachine.TaskStatuser {
	_va := make([]interface{}, len(keyvals))
	for _i := range keyvals {
		_va[_i] = keyvals[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 virtualmachine.TaskStatuser
	if rf, ok := ret.Get(0).(func(...string) virtualmachine.TaskStatuser); ok {
		r0 = rf(keyvals...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(virtualmachine.TaskStatuser)
		}
	}

	return r0
}

// StrArr provides a mock function with given fields: key, arr
func (_m *TaskStatuser) StrArr(key string, arr []string) virtualmachine.TaskStatuser {
	ret := _m.Called(key, arr)

	var r0 virtualmachine.TaskStatuser
	if rf, ok := ret.Get(0).(func(string, []string) virtualmachine.TaskStatuser); ok {
		r0 = rf(key, arr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(virtualmachine.TaskStatuser)
		}
	}

	return r0
}