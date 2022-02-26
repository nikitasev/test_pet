package service

import (
	"errors"
	"test_pet/internal/domain/entity"
)

var MissCacheError = errors.New("pin not match")

type UserListCache interface {
	SaveList(list []entity.User, limit, offset int32) error
	GetListByParams(limit, offset int32) ([]entity.User, error)
}
