// Code generated by mockery v1.0.0. DO NOT EDIT.

package usecase

import mock "github.com/stretchr/testify/mock"

// MockAppInfoRepository is an autogenerated mock type for the AppInfoRepository type
type MockAppInfoRepository struct {
	mock.Mock
}

// AppInfo provides a mock function with given fields:
func (_m *MockAppInfoRepository) AppInfo() (*AppInfoResponse, error) {
	ret := _m.Called()

	var r0 *AppInfoResponse
	if rf, ok := ret.Get(0).(func() *AppInfoResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AppInfoResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
