package auth

import (
	"context"
	"sso/internal/domain/models"
)

type AppProvider interface {
	App(ctx context.Context, appId int32) (models.App, error)
}
