package user_usecase

import (
	"context"
	"lab_fullcyle-auction_go/cmd/internal/entity/user_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
)

type UserUseCase struct {
	userRepo user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUserUseCase(userRepo user_entity.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

type UserUseCaseInterface interface {
	FindUserByID(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
}

func (uc *UserUseCase) FindUserByID(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {
	user, err := uc.userRepo.FindUserByID(id)
	if err != nil {
		return nil, internal_error.NewInternalServerError("failed to find user by id: " + id)
	}

	return &UserOutputDTO{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
