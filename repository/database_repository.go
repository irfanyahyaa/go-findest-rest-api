package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type DatabaseRepository[T any] interface {
	First(conds ...interface{}) (*T, error)
	Create(value *T) (*T, error)
	Find(filter string) ([]T, error)
	Save(value interface{}, conds ...interface{}) (*T, error)
}

type DatabaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

func NewDatabaseRepository[T any](db *gorm.DB) DatabaseRepository[T] {
	return &DatabaseRepositoryImpl[T]{
		db: db,
	}
}

func (r *DatabaseRepositoryImpl[T]) First(conds ...interface{}) (*T, error) {
	var entity T
	if err := r.db.First(&entity, conds...).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *DatabaseRepositoryImpl[T]) Create(value *T) (*T, error) {
	if err := r.db.Create(value).Error; err != nil {
		return nil, err
	}

	return value, nil
}

func (r *DatabaseRepositoryImpl[T]) Find(filter string) ([]T, error) {
	var entity []T
	query := "SELECT * FROM transactions"

	if filter != "" {
		query = fmt.Sprintf("%s %s", query, filter)
	}

	if err := r.db.Raw(query).Scan(&entity).Error; err != nil {
		return nil, err
	}

	r.db.Save(entity)

	return entity, nil
}

func (r *DatabaseRepositoryImpl[T]) Save(value interface{}, conds ...interface{}) (*T, error) {
	var entity T
	if err := r.db.Save(value).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&entity, conds...).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}
