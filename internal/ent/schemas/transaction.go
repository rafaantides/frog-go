package schemas

import (
	"frog-go/internal/utils/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Transaction struct {
	ent.Schema
}

func (Transaction) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
		mixins.RecordTypeMixin{},
		mixins.RecordStatusMixin{},
		mixins.MoneyMixin{Name: "amount"},
	}
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MaxLen(255).NotEmpty(),
		field.Time("record_date"),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique().Required().StorageKey(edge.Column("user_id")).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("invoice", Invoice.Type).Unique().StorageKey(edge.Column("invoice_id")).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("category", Category.Type).Unique().StorageKey(edge.Column("category_id")),
		// TODO: ver se tem como deixar category obrigatorio na modelagem, acredito q talvez n de por estar usando um hook para popular no create
	}
}

func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("record_date"),
		index.Fields("record_type"),
		index.Edges("invoice"),
		index.Edges("category"),
		index.Edges("category").Fields("record_date", "record_type"),
	}
}
