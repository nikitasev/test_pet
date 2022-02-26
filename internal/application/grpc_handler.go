package application

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"test_pet/internal/domain/entity"
	"test_pet/internal/domain/repository"
	"test_pet/internal/domain/service"
	"test_pet/pkg/grpc/userapi"
)

type GrpcHandler struct {
	userRepo    repository.User
	eventLogger service.NewUserEventLogger
	cache       service.UserListCache
	logger      *zap.Logger
	userapi.UnimplementedUserServiceServer
}

func NewGrpcHandler(userRepo repository.User, eventLogger service.NewUserEventLogger, cache service.UserListCache, logger *zap.Logger) *GrpcHandler {
	return &GrpcHandler{userRepo: userRepo, eventLogger: eventLogger, cache: cache, logger: logger}
}

func (h *GrpcHandler) AddUserRequest(ctx context.Context, input *userapi.AddUserInput) (*userapi.AddUserOutput, error) {
	var res userapi.AddUserOutput

	userId, err := h.userRepo.Add(input.GetName())
	if err != nil {
		h.logger.Error("failed adding to user list", zap.Error(err))

		return nil, errors.New("internal server error")
	}

	if err := h.eventLogger.Log(userId); err != nil {
		h.logger.Error("failed event logging", zap.Error(err))
	}

	res.Id = userId
	return &res, nil
}

func (h *GrpcHandler) DeleteUserRequest(ctx context.Context, input *userapi.DeleteUserInput) (*userapi.DeleteUserOutput, error) {
	var res userapi.DeleteUserOutput

	if err := h.userRepo.DeleteById(input.GetId()); err != nil {
		if errors.Is(err, repository.MissingClientWithId) {
			msg := "cannot delete client: " + err.Error()
			res.Error = &msg
			return &res, nil
		}

		h.logger.Error("failed user deleting", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	return &res, nil
}

func (h *GrpcHandler) GetListRequest(ctx context.Context, input *userapi.GetListInput) (*userapi.GetListOutput, error) {
	var (
		list []entity.User
		res  userapi.GetListOutput
		err  error
	)

	list, err = h.cache.GetListByParams(input.GetLimit(), input.GetOffset())
	if err == nil {
		res.List = convertListToResponse(list)
		return &res, nil
	}

	if !errors.Is(err, service.MissCacheError) {
		h.logger.Error("cache error", zap.Error(err))
	}

	list, err = h.userRepo.GetList(input.GetLimit(), input.GetOffset())
	if err != nil {
		h.logger.Error("failed loading user list", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	if err := h.cache.SaveList(list, input.GetLimit(), input.GetOffset()); err != nil {
		h.logger.Error("cache error", zap.Error(err))
	}

	res.List = convertListToResponse(list)
	return &res, nil
}

func convertListToResponse(list []entity.User) []*userapi.User {
	var result []*userapi.User

	for _, user := range list {
		result = append(result, &userapi.User{
			Id:   user.Id,
			Name: user.Name,
		})
	}

	return result
}
