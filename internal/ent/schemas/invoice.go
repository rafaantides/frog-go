package schemas

import (
	"frog-go/internal/utils/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Invoice struct {
	ent.Schema
}

func (Invoice) Mixin() []ent.Mixin {
	defaultZero := 0.0
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
		mixins.RecordStatusMixin{},
		mixins.MoneyMixin{Name: "amount", Default: &defaultZero},
	}
}

func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MaxLen(255).NotEmpty(),
		field.Time("due_date"),
	}
}

func (Invoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transactions", Transaction.Type).Ref("invoice"),
		edge.To("user", User.Type).Unique().Required().StorageKey(edge.Column("user_id")).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
