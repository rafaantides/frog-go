package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type MoneyMixin struct {
	mixin.Schema
	Name    string
	Default *float64
}

func (m MoneyMixin) Fields() []ent.Field {
	f := field.Float(m.Name).
		SchemaType(map[string]string{"postgres": "decimal(10,2)"})

	if m.Default != nil {
		f = f.Default(*m.Default)
	}

	return []ent.Field{f}
}
