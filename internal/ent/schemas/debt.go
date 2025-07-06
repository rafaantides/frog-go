package schemas

import (
	"frog-go/internal/utils/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Debt struct {
	ent.Schema
}

func (Debt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
		mixins.MoneyMixin{Name: "amount"},
	}
}

func (Debt) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MaxLen(255).NotEmpty(),
		field.Time("purchase_date"),
		field.Time("due_date"),
		field.UUID("category_id", uuid.UUID{}),
	}
}

func (Debt) Edges() []ent.Edge {
	return []ent.Edge{
		// TODO: ver se tem como deixar category obrigatorio na modelagem, acredito q talvez n de por estar usando um hook para popular no create
		edge.To("category", Category.Type).
			Unique().
			StorageKey(edge.Column("category_id")),
	}
}

func (Debt) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("purchase_date", "category_id"),
		index.Fields("due_date", "category_id"),
	}
}
