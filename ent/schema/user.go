package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/shinplay/pkg/publicid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
		mixin.UpdateTime{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("auth_id").MaxLen(24).DefaultFunc(publicid.Must).NotEmpty().Unique(),
		field.String("username").MaxLen(40).Optional().Unique(),
		field.String("first_name").Optional(),
		field.String("last_name").Optional(),
		field.String("email").Optional().Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
