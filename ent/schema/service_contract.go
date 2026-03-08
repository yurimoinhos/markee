package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// ServiceContract stores commercial terms for software delivery.
type ServiceContract struct {
	ent.Schema
}

func (ServiceContract) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("owner_user_id").NotEmpty(),
		field.String("customer_id").NotEmpty(),
		field.String("contract_type").NotEmpty(),
		field.String("title").NotEmpty().MaxLen(180),
		field.Uint64("amount_cents").GoType(common.Money{}.Value).Min(0),
		field.String("billing_type").Default("monthly").Validate(func(s string) error {
			valid := map[string]bool{"monthly": true, "milestone": true, "one_time": true}
			if !valid[s] {
				return fmt.Errorf("billing_type inválido: %q", s)
			}
			return nil
		}),
		field.String("status").Default("draft").Validate(func(s string) error {
			valid := map[string]bool{"draft": true, "active": true, "expired": true, "cancelled": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.Time("start_date").Default(time.Now),
		field.Time("end_date").Optional().Nillable(),
		field.Int("duration_months").Optional().Nillable(),
		field.String("deliverables").Optional().Nillable(),
		field.String("sla").Optional().Nillable(),
		field.String("penalties").Optional().Nillable(),
		field.String("payment_terms").Optional().Nillable(),
		field.Bool("auto_renew").Default(false),
		field.Time("next_renewal_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
