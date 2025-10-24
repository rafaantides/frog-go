package postgresql

import (
	"context"
	"fmt"
	"frog-go/internal/core/domain"
	appError "frog-go/internal/core/errors"
)

const userSessionEntity = "user_sessions"

// TODO: Implement the methods

func (p *PostgreSQL) GetUserSessionByToken(ctx context.Context, token string) (*domain.UserSession, error) {
	return nil, appError.FailedToFind(userSessionEntity, fmt.Errorf("not implemented"))
}

func (p *PostgreSQL) CreateUserSession(ctx context.Context, session domain.UserSession) error {
	return nil

}
func (p *PostgreSQL) DeleteUserSession(ctx context.Context, session domain.UserSession) error {
	return nil
}
