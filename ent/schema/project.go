package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Project links delivery execution to a service contract.
type Project struct {
	ent.Schema
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("owner_user_id").NotEmpty(),
		field.String("contract_id").NotEmpty(),
		field.String("name").NotEmpty().MaxLen(200),
		field.String("status").Default("active").Validate(func(s string) error {
			valid := map[string]bool{"active": true, "on_hold": true, "done": true, "cancelled": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
