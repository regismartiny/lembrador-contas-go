package user_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/user_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type UserInputDTO struct {
	Name  string `json:"name" binding:"required,min=1"`
	Email string `json:"email" binding:"required,min=5"`
}

type CreateUserOutputDTO struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05"`
}

type UserUseCaseInterface interface {
	CreateUser(
		ctx context.Context,
		userInput UserInputDTO) *internal_error.InternalError
	FindUserById(
		ctx context.Context,
		id string) (*UserOutputDTO, *internal_error.InternalError)
	FindUsers(
		ctx context.Context,
		status user_entity.UserStatus,
		name, email string) ([]UserOutputDTO, *internal_error.InternalError)
}

type UserUseCase struct {
	userRepository user_entity.UserRepositoryInterface
}

func NewUserUseCase(
	userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		userRepository: userRepository,
	}
}

func (u *UserUseCase) CreateUser(
	ctx context.Context,
	userInput UserInputDTO) *internal_error.InternalError {

	user, err := user_entity.CreateUser(userInput.Name, userInput.Email)
	if err != nil {
		return err
	}

	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
