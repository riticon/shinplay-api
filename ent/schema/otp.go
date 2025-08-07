package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/shinplay/pkg/publicid"
)

// OTP holds the schema definition for the OTP entity.
type OTP struct {
	ent.Schema
}

func (OTP) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "otps"},
	}
}

// Fields of the OTP.
func (OTP) Fields() []ent.Field {
	return []ent.Field{
		field.String("otp").MaxLen(6).NotEmpty().DefaultFunc(func() string {
			return publicid.MustWith(4, publicid.Numberic())
		}).Unique(),
		field.Time("expires_at").Default(time.Now().Add(5 * time.Minute)),
	}
}

// Edges of the OTP.
func (OTP) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("otps").
			Unique().
			Required(), // otp must belong to a user
	}
}

func (OTP) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
		mixin.UpdateTime{},
	}
}
