package mixins

import (
	"fmt"
	"frog-go/internal/core/domain"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type RecordStatusMixin struct {
	mixin.Schema
	Name string
}

func (m RecordStatusMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("status").
			NotEmpty().
			Default(string(domain.StatusPending)).
			Validate(func(s string) error {
				if !domain.TxnStatus(s).IsValid() {
					return fmt.Errorf("invalid status: %q", s)
				}
				return nil
			}),
	}
}
