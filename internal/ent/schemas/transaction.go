package schemas

import (
	"fmt"
	"frog-go/internal/core/domain"
	"frog-go/internal/utils/mixins"

	"entgo.io/ent"
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
		mixins.TxnKindMixin{},
		mixins.MoneyMixin{Name: "amount"},
	}
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MaxLen(255).NotEmpty(),
		field.Time("purchase_date"),
		field.Time("due_date").Nillable().Optional(),

		field.String("status").
			NotEmpty().
			Default(string(domain.TxnStatusPending)).
			Validate(func(s string) error {
				if !domain.TxnStatus(s).IsValid() {
					return fmt.Errorf("invalid status: %q", s)
				}
				return nil
			}),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		// TODO: ver se tem como deixar category obrigatorio na modelagem, acredito q talvez n de por estar usando um hook para popular no create
		edge.To("category", Category.Type).
			Unique().
			StorageKey(edge.Column("category_id")),
	}
}

func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("purchase_date"),
		index.Fields("due_date"),
		index.Fields("kind"),
		index.Edges("category"),
		index.Edges("category").Fields("due_date", "kind"),
		index.Edges("category").Fields("purchase_date", "kind"),
	}
}
