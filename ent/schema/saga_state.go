package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// SagaState persiste o estado corrente de uma saga no Postgres.
// Permite recuperação pós-crash e auditoria do fluxo.
type SagaState struct {
	ent.Schema
}

func (SagaState) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string {
			return common.GenID().Value
		}),
		field.String("saga_type").NotEmpty().Comment("ex: payment_saga"),
		field.String("order_id").NotEmpty(),
		field.String("user_id").NotEmpty(),
		field.Uint64("amount_cents").GoType(common.Money{}.Value).Min(1),
		field.String("current_step").Default("pending"),
		field.String("status").
			Default("pending").
			Validate(func(s string) error {
				valid := map[string]bool{
					"pending":          true,
					"balance_reserved": true,
					"payment_sent":     true,
					"confirmed":        true,
					"compensating":     true,
					"cancelled":        true,
				}
				if !valid[s] {
					return fmt.Errorf("saga status inválido: %q", s)
				}
				return nil
			}),
		field.String("idempotency_key").NotEmpty().Unique(),
		field.String("provider_ref").Optional().Nillable(),
		field.String("failure_reason").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
