// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package exporter

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// MockExporter is an autogenerated mock type for the Exporter type
type MockExporter struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockExporter) Close() {
	_m.Called()
}

// WriteMetrics provides a mock function with given fields: w
func (_m *MockExporter) WriteMetrics(w io.Writer) error {
	ret := _m.Called(w)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer) error); ok {
		r0 = rf(w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
