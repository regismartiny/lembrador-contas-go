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

type UserStatus int

func (s UserStatus) String() string {
	switch s {
	case Active:
		return "active"
	case Inactive:
		return "inactive"
	default:
		return "unknown"
	}
}

const (
	Active UserStatus = iota
	Inactive
)

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
	FindUserById(ctx context.Context, userId string) (*User, *internal_error.InternalError)
	CreateUser(ctx context.Context, userEntity *User) *internal_error.InternalError
}
