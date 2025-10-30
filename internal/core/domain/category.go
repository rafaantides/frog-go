package domain

import (
	appError "frog-go/internal/core/errors"
	"frog-go/internal/utils"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID                  uuid.UUID
	Name                string
	Description         *string
	Color               *string
	SuggestedPercentage *int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewCategory(
	id *uuid.UUID,
	name string,
	description, color *string,
	suggestedPercentage *int,
	createdAt *time.Time,
	updatedAt *time.Time,
) (*Category, error) {
	if name == "" {
		return nil, appError.EmptyField("name")
	}

	return &Category{
		ID:                  utils.UUIDOrZero(id),
		Name:                name,
		Description:         description,
		Color:               color,
		SuggestedPercentage: suggestedPercentage,
		CreatedAt:           utils.TimeOrNow(createdAt),
		UpdatedAt:           utils.TimeOrNow(updatedAt),
	}, nil
}
