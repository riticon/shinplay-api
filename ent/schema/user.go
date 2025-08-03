package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
		field.String("auth_id").DefaultFunc(publicid.Must).NotEmpty().Unique().MaxLen(24),
		field.String("username").MaxLen(40).Optional().Unique(),
		field.String("email").Optional().Unique(),
		field.String("phone_number").Optional().Unique().MaxLen(15).MinLen(7),
		field.String("first_name").Optional(),
		field.String("last_name").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", Session.Type),
		edge.To("otps", OTP.Type),
	}
}
