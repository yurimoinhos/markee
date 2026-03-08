package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// PaymentRecord stores settled or pending payment transactions.
type PaymentRecord struct {
	ent.Schema
}

func (PaymentRecord) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("charge_id").NotEmpty(),
		field.Uint64("amount_cents").GoType(common.Money{}.Value).Min(1),
		field.String("method").Default("pix"),
		field.String("status").Default("pending").Validate(func(s string) error {
			valid := map[string]bool{"pending": true, "paid": true, "failed": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.Time("due_date").Optional().Nillable(),
		field.Time("paid_at").Optional().Nillable(),
		field.String("receipt_url").Optional().Nillable().MaxLen(500),
		field.String("external_id").Optional().Nillable(),
		field.String("tx_hash").Optional().Nillable().MaxLen(256),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
