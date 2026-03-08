package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// PaymentEvidence stores files and notes attached to a payment.
type PaymentEvidence struct {
	ent.Schema
}

func (PaymentEvidence) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("payment_id").NotEmpty(),
		field.String("file_url").Optional().Nillable().MaxLen(500),
		field.String("note").Optional().Nillable().MaxLen(1000),
		field.String("tx_hash").Optional().Nillable().MaxLen(256),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
