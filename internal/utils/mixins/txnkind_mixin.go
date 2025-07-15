package mixins

import (
	"fmt"
	"frog-go/internal/core/domain"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type RecordTypeMixin struct {
	mixin.Schema
	Name string
}

func (m RecordTypeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("record_type").
			NotEmpty().
			Default(string(domain.TypeExpense)).
			Validate(func(s string) error {
				if !domain.RecordType(s).IsValid() {
					return fmt.Errorf("invalid record_type: %q", s)
				}
				return nil
			}),
	}
}
