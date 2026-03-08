package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Worklog records hours worked for a project/milestone.
type Worklog struct {
	ent.Schema
}

func (Worklog) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("project_id").NotEmpty(),
		field.String("milestone_id").Optional().Nillable(),
		field.String("user_id").NotEmpty(),
		field.Float("hours").Positive(),
		field.String("description").Optional().Nillable().MaxLen(1000),
		field.Time("worked_at").Default(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
