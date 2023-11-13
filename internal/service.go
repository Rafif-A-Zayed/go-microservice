package user

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)

// Service describes the Order service.
type Service interface {
	Create(ctx context.Context, user User) (string, error)
	GetByID(ctx context.Context, id string) (User, error)
	ChangeStatus(ctx context.Context, id string, status string) error
}
