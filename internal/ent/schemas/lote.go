package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Teste struct {
	ent.Schema
}

func (Teste) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
