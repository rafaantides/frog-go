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
			dm, ok := m.(*ent.TransactionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type: %T", m)
			}

			if !dm.Op().Is(ent.OpCreate) {
				return next.Mutate(ctx, m)
			}

			title, exists := dm.Title()
			if !exists {
				return nil, fmt.Errorf("title is required to categorize transaction")
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

func UpdateInvoiceAmountHook(client *ent.Client) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			dm, ok := m.(*ent.TransactionMutation)
			if !ok {
				return next.Mutate(ctx, m)
			}

			invoiceID, hasInvoice := dm.InvoiceID()
			if !hasInvoice {
				return next.Mutate(ctx, m)
			}

			newAmount, hasNewAmount := dm.Amount()
			if !hasNewAmount {
				return next.Mutate(ctx, m)
			}

			var delta float64

			switch {
			case dm.Op().Is(ent.OpCreate):
				delta = newAmount

			case dm.Op().Is(ent.OpUpdateOne) || dm.Op().Is(ent.OpUpdate):
				id, ok := dm.ID()
				if !ok {
					return nil, fmt.Errorf("missing transaction ID during update")
				}
				oldTransaction, err := client.Transaction.Get(ctx, id)
				if err != nil {
					return nil, fmt.Errorf("failed to load old transaction: %w", err)
				}

				// Se o valor não mudou, ignora
				if oldTransaction.Amount == newAmount {
					return next.Mutate(ctx, m)
				}

				delta = newAmount - oldTransaction.Amount
			}

			err := client.Invoice.
				UpdateOneID(invoiceID).
				AddAmount(delta).
				Exec(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to update invoice amount: %w", err)
			}

			return next.Mutate(ctx, m)
		})
	}
}
