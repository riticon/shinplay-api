package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shinplay/pkg/publicid"
)

// OTP holds the schema definition for the OTP entity.
type OTP struct {
	ent.Schema
}

// Fields of the OTP.
func (OTP) Fields() []ent.Field {
	return []ent.Field{
		field.String("otp").MaxLen(6).NotEmpty().DefaultFunc(func() string {
			return publicid.MustWith(6, publicid.Numberic())
		}).Unique(),
		field.Time("expires_at").Default(time.Now().Add(5 * time.Minute)),
	}
}

// Edges of the OTP.
func (OTP) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("otps").Unique().Required(),
	}
}
