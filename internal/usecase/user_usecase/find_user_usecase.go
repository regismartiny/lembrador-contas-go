package user_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/user_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

func FindUserUseCase(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		userRepository,
	}
}

type UserOutputDTO struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05"`
}

func (u *UserUseCase) FindUserById(
	ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.userRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:        userEntity.Id,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Status:    user_entity.UserStatus(userEntity.Status).Name(),
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}

func (u *UserUseCase) FindUsers(
	ctx context.Context,
	status user_entity.UserStatus,
	name, email string) ([]UserOutputDTO, *internal_error.InternalError) {
	userEntities, err := u.userRepository.FindUsers(
		ctx, status, name, email)
	if err != nil {
		return nil, err
	}

	var userOutputs []UserOutputDTO
	for _, value := range userEntities {
		userOutputs = append(userOutputs, UserOutputDTO{
			Id:        value.Id,
			Name:      value.Name,
			Email:     value.Email,
			Status:    user_entity.UserStatus(value.Status).Name(),
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		})
	}

	return userOutputs, nil
}
