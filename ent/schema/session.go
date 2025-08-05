package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/shinplay/pkg/publicid"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").NotEmpty().Unique().DefaultFunc(func() string {
			return publicid.MustWith(30, publicid.AlphaNumeric())
		}),
		field.Text("refresh_token").Comment("Refresh token for the session"),
		field.Time("expires_at").Comment("Session expiration time"),
		field.String("user_agent").Optional().Comment("User agent string of the session"),
		field.String("ip_address").Optional().Comment("IP address of the user"),
	}
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
		mixin.UpdateTime{},
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Unique().
			Required(), // session must belong to a user
	}
}
