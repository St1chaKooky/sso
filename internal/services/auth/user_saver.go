package auth

import "context"

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}
