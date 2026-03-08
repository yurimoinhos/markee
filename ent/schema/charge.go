package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Charge stores each billing request emitted by the platform.
type Charge struct {
	ent.Schema
}

func (Charge) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("owner_user_id").NotEmpty(),
		field.String("customer_id").NotEmpty(),
		field.String("contract_id").Optional().Nillable(),
		field.String("milestone_id").Optional().Nillable(),
		field.String("charge_type").Default("monthly").Validate(func(s string) error {
			valid := map[string]bool{"monthly": true, "milestone": true, "one_time": true}
			if !valid[s] {
				return fmt.Errorf("charge_type inválido: %q", s)
			}
			return nil
		}),
		field.Uint64("amount_cents").GoType(common.Money{}.Value).Min(1),
		field.String("currency").Default("BRL"),
		field.String("payment_method").Default("pix").Validate(func(s string) error {
			valid := map[string]bool{"pix": true, "bank_transfer": true, "credit_card": true, "crypto": true, "boleto": true}
			if !valid[s] {
				return fmt.Errorf("payment_method inválido: %q", s)
			}
			return nil
		}),
		field.String("status").Default("pending").Validate(func(s string) error {
			valid := map[string]bool{"pending": true, "paid": true, "overdue": true, "cancelled": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.Time("due_date").Default(time.Now),
		field.Time("paid_at").Optional().Nillable(),
		field.String("external_id").Optional().Nillable(),
		field.String("payment_link").Optional().Nillable().MaxLen(500),
		field.String("qr_code").Optional().Nillable(),
		field.String("description").Optional().Nillable().MaxLen(500),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
