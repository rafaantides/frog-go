package schemas

import (
	"frog-go/internal/utils/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Category struct {
	ent.Schema
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
		mixins.TxnKindMixin{},
	}
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().MaxLen(255).NotEmpty(),
		field.String("description").Optional().Nillable(),
		field.String("color").MaxLen(7).Optional().Nillable(),
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind"),
	}
}
