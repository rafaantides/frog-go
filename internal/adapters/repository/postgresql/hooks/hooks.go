package hooks

import (
	"context"
	"fmt"

	"frog-go/internal/ent"
	"frog-go/internal/ent/category"
)

func SetCategoryFromTitleHook(client *ent.Client, categorizer *Categorizer) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			dm, ok := m.(*ent.DebtMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type: %T", m)
			}

			if !dm.Op().Is(ent.OpCreate) {
				return next.Mutate(ctx, m)
			}

			title, exists := dm.Title()
			if !exists {
				return nil, fmt.Errorf("title is required to categorize debt")
			}

			categoryName := categorizer.Categorize(title)

			data, err := client.Category.
				Query().
				Where(category.NameEQ(categoryName)).
				Only(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to find category '%s': %w", categoryName, err)
			}

			dm.SetCategoryID(data.ID)

			return next.Mutate(ctx, dm)
		})
	}
}