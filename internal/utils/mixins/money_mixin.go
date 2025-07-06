package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type MoneyMixin struct {
	mixin.Schema
	Name string
}

func (m MoneyMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Float(m.Name).
			SchemaType(map[string]string{"postgres": "decimal(10,2)"}),
	}
}
