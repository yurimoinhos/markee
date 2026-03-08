package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// AutomationRun stores execution logs for financial automations.
type AutomationRun struct {
	ent.Schema
}

func (AutomationRun) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("owner_user_id").NotEmpty(),
		field.String("automation_type").NotEmpty(),
		field.String("status").Default("success").Validate(func(s string) error {
			valid := map[string]bool{"success": true, "failed": true, "skipped": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.String("details").Optional().Nillable(),
		field.Time("ran_at").Default(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
