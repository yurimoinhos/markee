package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string {
			return common.GenID().Value
		}),
		field.String("first_name").MinLen(3).MaxLen(40).NotEmpty(),
		field.String("last_name").MinLen(3).MaxLen(70).NotEmpty(),
		field.Uint64("balance").GoType(common.Money{}.Value).Min(0),
		field.String("phone_number").MinLen(8).MaxLen(15).Nillable().Optional(),
		field.String("email").MinLen(5).MaxLen(100).NotEmpty().Unique().Immutable(),
		field.String("password_hash").Nillable().Optional().Sensitive(),
		field.String("oauth_provider").Nillable().Optional(),
		field.String("oauth_sub").Nillable().Optional().Unique(),
		field.Bool("active").Default(false),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Nillable().Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", Group.Type).Ref("users"),
		edge.To("roles", Role.Type),
	}
}
