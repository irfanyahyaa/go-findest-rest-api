package mock

import (
	"github.com/stretchr/testify/mock"
	"go-findest-rest-api/dto"
)

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

func (m *MockDatabaseRepository[T]) Find(filter string) ([]T, error) {
	args := m.Called(filter)
	if args.Get(0) != nil {
		return args.Get(0).([]T), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDatabaseRepository[T]) Save(value interface{}, conds ...interface{}) (*T, error) {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	args = m.Called(conds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*T), args.Error(1)
}

func (m *MockDatabaseRepository[T]) AverageTransaction() ([]dto.AverageTransactionAttr, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]dto.AverageTransactionAttr), args.Error(1)
	}
	return nil, args.Error(1)
}
