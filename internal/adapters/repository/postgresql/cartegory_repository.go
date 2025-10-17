package postgresql

import (
	"context"
	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/ent"
	"frog-go/internal/ent/category"
	"frog-go/internal/utils/pagination"

	"github.com/google/uuid"
)

const categoryEntity = "categories"

func (p *PostgreSQL) GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	row, err := p.Client.Category.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, appError.ErrNotFound
		}
		return nil, appError.FailedToFind(categoryEntity, err)
	}
	return dto.NewCategoryResponse(row.ID, row.Name, row.Description, row.Color, row.SuggestedPercentage), nil
}

func (p *PostgreSQL) GetCategoryIDByName(ctx context.Context, name *string) (*uuid.UUID, error) {
	if name == nil {
		return nil, nil
	}

	data, err := p.Client.Category.Query().Where(category.NameEQ(*name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, appError.ErrNotFound
		}
		return nil, appError.FailedToFind(categoryEntity, err)
	}

	id := data.ID
	return &id, nil
}

func (p *PostgreSQL) CreateCategory(ctx context.Context, input domain.Category) (*dto.CategoryResponse, error) {
	row, err := p.Client.Category.
		Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableColor(input.Color).
		SetNillableSuggestedPercentage(input.SuggestedPercentage).
		Save(ctx)

	if err != nil {
		return nil, appError.FailedToSave(categoryEntity, err)
	}

	return dto.NewCategoryResponse(row.ID, row.Name, row.Description, row.Color, row.SuggestedPercentage), nil
}

func (p *PostgreSQL) UpdateCategory(ctx context.Context, id uuid.UUID, input domain.Category) (*dto.CategoryResponse, error) {
	row, err := p.Client.Category.
		UpdateOneID(id).
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableColor(input.Color).
		SetNillableSuggestedPercentage(input.SuggestedPercentage).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, appError.ErrNotFound
		}
		return nil, appError.FailedToUpdate(categoryEntity, err)
	}

	return dto.NewCategoryResponse(row.ID, row.Name, row.Description, row.Color, row.SuggestedPercentage), nil
}

func (p *PostgreSQL) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	err := p.Client.Category.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return appError.ErrNotFound
		}
		return appError.FailedToDelete(categoryEntity, err)
	}
	return nil
}

func (p *PostgreSQL) ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, error) {
	query := p.Client.Category.Query()
	query = applyCategoryFilters(query, pgn)

	if pgn.OrderDirection == config.OrderAsc {
		query = query.Order(ent.Asc(pgn.OrderBy))
	} else {
		query = query.Order(ent.Desc(pgn.OrderBy))
	}

	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	rows, err := query.All(ctx)
	if err != nil {
		return []dto.CategoryResponse{}, err
	}

	response := make([]dto.CategoryResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, *dto.NewCategoryResponse(row.ID, row.Name, row.Description, row.Color, row.SuggestedPercentage))
	}
	return response, nil

}

func (p *PostgreSQL) CountCategories(ctx context.Context, pgn *pagination.Pagination) (int, error) {
	query := p.Client.Category.Query()
	query = applyCategoryFilters(query, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func applyCategoryFilters(query *ent.CategoryQuery, pgn *pagination.Pagination) *ent.CategoryQuery {
	if pgn.Search != "" {
		query = query.Where(
			category.Or(
				category.NameContainsFold(pgn.Search),
				category.DescriptionContainsFold(pgn.Search),
			),
		)
	}
	return query
}
