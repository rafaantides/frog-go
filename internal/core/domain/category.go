package domain

import (
	"frog-go/internal/core/errors"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Color       *string
	CreatedAt   string
	UpdatedAt   string
}

func NewCategory(name string, description, color *string) (*Category, error) {

	if name == "" {
		return nil, errors.EmptyField("name")
	}

	return &Category{
		Name:        name,
		Description: description,
		Color:       color,
	}, nil
}
