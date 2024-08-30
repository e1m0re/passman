// Code generated by mockery v2.44.2. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/e1m0re/passman/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// StoreService is an autogenerated mock type for the StoreService type
type StoreService struct {
	mock.Mock
}

// AddItem provides a mock function with given fields: ctx, datumInfo
func (_m *StoreService) AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error) {
	ret := _m.Called(ctx, datumInfo)

	if len(ret) == 0 {
		panic("no return value specified for AddItem")
	}

	var r0 *model.DatumItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.DatumInfo) (*model.DatumItem, error)); ok {
		return rf(ctx, datumInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.DatumInfo) *model.DatumItem); ok {
		r0 = rf(ctx, datumInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DatumItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.DatumInfo) error); ok {
		r1 = rf(ctx, datumInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStoreService creates a new instance of StoreService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStoreService(t interface {
	mock.TestingT
	Cleanup(func())
}) *StoreService {
	mock := &StoreService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
