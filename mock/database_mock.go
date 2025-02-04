package mock

import "github.com/stretchr/testify/mock"

type MockDatabaseRepository[T any] struct {
	mock.Mock
}

func (m *MockDatabaseRepository[T]) First(conds ...interface{}) (*T, error) {
	args := m.Called(conds...)
	if args.Get(0) != nil {
		return args.Get(0).(*T), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDatabaseRepository[T]) Create(value *T) (*T, error) {
	args := m.Called(value)
	return args.Get(0).(*T), args.Error(1)
}
