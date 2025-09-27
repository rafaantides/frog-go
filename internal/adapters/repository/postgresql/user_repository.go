package postgresql

import (
	"context"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/errors"
	"frog-go/internal/ent"
	"frog-go/internal/ent/user"
)

const userEntity = "users"

func (p *PostgreSQL) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := p.Client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.FailedToFind(userEntity, err)
	}

	return &domain.User{
		ID:           row.ID,
		Name:         row.Name,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (p *PostgreSQL) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := p.Client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.FailedToFind(userEntity, err)
	}

	return &domain.User{
		ID:           row.ID,
		Name:         row.Name,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}
