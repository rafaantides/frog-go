package utilsctx

import (
	"context"
	appError "frog-go/internal/core/errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value(UserIDKey)
	if v == nil {
		return uuid.Nil, appError.ErrUserNotFoundInCtx
	}

	id, ok := v.(uuid.UUID)
	if !ok {
		return uuid.Nil, appError.ErrUserNotFoundInCtx
	}

	return id, nil
}
