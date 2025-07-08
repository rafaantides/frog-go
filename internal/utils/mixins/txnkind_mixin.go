package mixins

import (
	"fmt"
	"frog-go/internal/core/domain"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TxnKindMixin struct {
	mixin.Schema
	Name string
}

func (m TxnKindMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").
			NotEmpty().
			Default(string(domain.TxnKindExpense)).
			Validate(func(s string) error {
				if !domain.TxnKind(s).IsValid() {
					return fmt.Errorf("invalid kind: %q", s)
				}
				return nil
			}),
	}
}
