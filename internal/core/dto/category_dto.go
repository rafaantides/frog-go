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

func DomainToCategoryResponse(data domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:                  data.ID,
		Name:                data.Name,
		Description:         data.Description,
		Color:               data.Color,
		SuggestedPercentage: data.SuggestedPercentage,
	}
}

func DomainToCategoryResponseList(domains []domain.Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(domains))
	for i, d := range domains {
		responses[i] = DomainToCategoryResponse(d)
	}
	return responses
}

func (r *CategoryRequest) ToDomain() (*domain.Category, error) {
	return domain.NewCategory(nil, r.Name, r.Description, r.Color, r.SuggestedPercentage, nil, nil)

}
