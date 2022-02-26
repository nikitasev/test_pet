package repository

import (
	"errors"
	"test_pet/internal/domain/entity"
)

var MissingClientWithId = errors.New("there is no such client")

type User interface {
	GetList(limit, offset int32) ([]entity.User, error)
	Add(name string) (int64, error)
	DeleteById(userId int64) error
}
