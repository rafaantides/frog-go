package schemas

import (
	"frog-go/internal/utils/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).Optional(),
		field.String("username").Unique().NotEmpty().MaxLen(255),
		field.String("email").Unique().Optional().MaxLen(255),
		field.String("password_hash").Sensitive().NotEmpty().MaxLen(255),
		field.Bool("is_active").Default(true),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transactions", Transaction.Type).Ref("user"),
		edge.From("invoices", Invoice.Type).Ref("user"),
	}
}
