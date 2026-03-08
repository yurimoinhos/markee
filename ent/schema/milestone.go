package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Milestone tracks planned deliverables and optional billing amount.
type Milestone struct {
	ent.Schema
}

func (Milestone) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("project_id").NotEmpty(),
		field.String("contract_id").Optional().Nillable(),
		field.String("title").NotEmpty().MaxLen(200),
		field.String("deliverables").Optional().Nillable(),
		field.Uint64("amount_cents").Optional().Nillable().GoType(common.Money{}.Value),
		field.Time("due_date").Optional().Nillable(),
		field.String("status").Default("pending").Validate(func(s string) error {
			valid := map[string]bool{"pending": true, "in_progress": true, "done": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
