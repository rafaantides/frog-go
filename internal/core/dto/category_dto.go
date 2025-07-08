package dto

import (
	"frog-go/internal/core/domain"

	"github.com/google/uuid"
)

type CategoryRequest struct {
	Name        string  `json:"name"`
	Kind        string  `json:"kind" validate:"required,oneof=income expense"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Kind        string    `json:"kind"`
	Description *string   `json:"description"`
	Color       *string   `json:"color"`
}

func NewCategoryResponse(id uuid.UUID, name string, description, color *string) *CategoryResponse {
	return &CategoryResponse{
		ID:          id,
		Name:        name,
		Description: description,
		Color:       color,
	}
}

func (r *CategoryRequest) ToDomain() (*domain.Category, error) {
	kind := domain.TxnKind(r.Kind)
	return domain.NewCategory(r.Name, &kind, r.Description, r.Color)

}
