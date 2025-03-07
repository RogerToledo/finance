// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	repository "github.com/me/finance/pkg/repository"
	mock "github.com/stretchr/testify/mock"
)

// RepositoryAll is an autogenerated mock type for the RepositoryAll type
type RepositoryAll struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *RepositoryAll) All() *repository.Repository {
	ret := _m.Called()

	var r0 *repository.Repository
	if rf, ok := ret.Get(0).(func() *repository.Repository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Repository)
		}
	}

	return r0
}

type mockConstructorTestingTNewRepositoryAll interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepositoryAll creates a new instance of RepositoryAll. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepositoryAll(t mockConstructorTestingTNewRepositoryAll) *RepositoryAll {
	mock := &RepositoryAll{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
