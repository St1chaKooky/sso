package auth

import (
	"context"
	"sso/internal/domain/models"
)

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
