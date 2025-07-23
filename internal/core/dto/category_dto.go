package dto

import (
	"frog-go/internal/core/domain"

	"github.com/google/uuid"
)

type CategoryRequest struct {
	Name                string  `json:"name"`
	Description         *string `json:"description"`
	Color               *string `json:"color"`
	SuggestedPercentage *int    `json:"suggested_percentage"`
}

type CategoryResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         *string   `json:"description"`
	Color               *string   `json:"color"`
	SuggestedPercentage *int      `json:"suggested_percentage"`
}

func NewCategoryResponse(id uuid.UUID, name string, description, color *string, suggestedPercentage *int) *CategoryResponse {
	return &CategoryResponse{
		ID:                  id,
		Name:                name,
		Description:         description,
		Color:               color,
		SuggestedPercentage: suggestedPercentage,
	}
}

func (r *CategoryRequest) ToDomain() (*domain.Category, error) {
	return domain.NewCategory(r.Name, r.Description, r.Color, r.SuggestedPercentage)

}
