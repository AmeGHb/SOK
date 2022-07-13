package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user *User) error
	FindAll(ctx context.Context) (users []User, err error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User, transaction float64, sign string) error
}
