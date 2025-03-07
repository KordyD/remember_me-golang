// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/kordyd/remember_me-golang/sso/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockAppStorage is an autogenerated mock type for the AppStorage type
type MockAppStorage struct {
	mock.Mock
}

type MockAppStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAppStorage) EXPECT() *MockAppStorage_Expecter {
	return &MockAppStorage_Expecter{mock: &_m.Mock}
}

// GetApp provides a mock function with given fields: ctx, id
func (_m *MockAppStorage) GetApp(ctx context.Context, id string) (models.App, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetApp")
	}

	var r0 models.App
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.App, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.App); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.App)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAppStorage_GetApp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetApp'
type MockAppStorage_GetApp_Call struct {
	*mock.Call
}

// GetApp is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockAppStorage_Expecter) GetApp(ctx interface{}, id interface{}) *MockAppStorage_GetApp_Call {
	return &MockAppStorage_GetApp_Call{Call: _e.mock.On("GetApp", ctx, id)}
}

func (_c *MockAppStorage_GetApp_Call) Run(run func(ctx context.Context, id string)) *MockAppStorage_GetApp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAppStorage_GetApp_Call) Return(_a0 models.App, _a1 error) *MockAppStorage_GetApp_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAppStorage_GetApp_Call) RunAndReturn(run func(context.Context, string) (models.App, error)) *MockAppStorage_GetApp_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAppStorage creates a new instance of MockAppStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAppStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAppStorage {
	mock := &MockAppStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
