package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/shinplay/pkg/publicid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("auth_id").Unique().Immutable().Default(publicid.Must()).
			Comment("Unique identifier for the user, generated using publicid package."),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
