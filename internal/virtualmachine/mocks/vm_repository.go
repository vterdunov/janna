// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import virtualmachine "github.com/vterdunov/janna/internal/virtualmachine"

// VMRepository is an autogenerated mock type for the VMRepository type
type VMRepository struct {
	mock.Mock
}

// VMDeploy provides a mock function with given fields: params
func (_m *VMRepository) VMDeploy(params virtualmachine.VMDeployRequest) (virtualmachine.VMDeployResponse, error) {
	ret := _m.Called(params)

	var r0 virtualmachine.VMDeployResponse
	if rf, ok := ret.Get(0).(func(virtualmachine.VMDeployRequest) virtualmachine.VMDeployResponse); ok {
		r0 = rf(params)
	} else {
		r0 = ret.Get(0).(virtualmachine.VMDeployResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(virtualmachine.VMDeployRequest) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VMInfo provides a mock function with given fields: uuid
func (_m *VMRepository) VMInfo(uuid string) (virtualmachine.VMInfoResponse, error) {
	ret := _m.Called(uuid)

	var r0 virtualmachine.VMInfoResponse
	if rf, ok := ret.Get(0).(func(string) virtualmachine.VMInfoResponse); ok {
		r0 = rf(uuid)
	} else {
		r0 = ret.Get(0).(virtualmachine.VMInfoResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}