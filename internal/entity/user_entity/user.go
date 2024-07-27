package user_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type User struct {
	Id        string
	Name      string
	Email     string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserStatus uint8

const (
	Active UserStatus = iota
	Inactive
)

func (s UserStatus) Name() string {
	return userStatusNames[s]
}

var userStatusNames = []string{
	"active",
	"inactive",
}

func GetUserStatusByName(name string) (UserStatus, *internal_error.InternalError) {
	for k, v := range userStatusNames {
		if v == name {
			return UserStatus(k), nil
		}
	}

	return UserStatus(0), internal_error.NewBadRequestError("invalid user status name")
}

func CreateUser(
	name string,
	email string) (*User, *internal_error.InternalError) {

	user :=
		&User{
			Id:        uuid.New().String(),
			Name:      name,
			Email:     email,
			Status:    Active,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) Validate() *internal_error.InternalError {
	if len(user.Name) <= 1 ||
		len(user.Email) <= 5 {
		return internal_error.NewBadRequestError("invalid user object")
	}

	return nil
}

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, userEntity *User) *internal_error.InternalError
	FindUserById(ctx context.Context, userId string) (*User, *internal_error.InternalError)
	FindUsers(
		ctx context.Context,
		status UserStatus,
		name, email string) ([]User, *internal_error.InternalError)
}
