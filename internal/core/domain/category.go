package domain

import (
	"frog-go/internal/core/errors"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	Kind        TxnKind `json:"kind"`
	Description *string
	Color       *string
	CreatedAt   string
	UpdatedAt   string
}

func NewCategory(name string, kind *TxnKind, description, color *string) (*Category, error) {
	if name == "" {
		return nil, errors.EmptyField("name")
	}

	kindValue := TxnKindExpense
	if kind != nil {
		kindValue = *kind
	}

	return &Category{
		Name:        name,
		Description: description,
		Color:       color,
		Kind:        kindValue,
	}, nil
}
