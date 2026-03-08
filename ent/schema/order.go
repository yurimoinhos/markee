package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Order representa um pedido de pagamento associado a uma saga.
type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string {
			return common.GenID().Value
		}),
		field.String("saga_id").NotEmpty(),
		field.String("user_id").NotEmpty(),
		field.Uint64("amount").GoType(common.Money{}.Value).Min(1),
		field.String("status").
			Default("pending").
			Validate(func(s string) error {
				valid := map[string]bool{
					"pending": true, "confirmed": true, "cancelled": true,
				}
				if !valid[s] {
					return fmt.Errorf("status inválido: %q", s)
				}
				return nil
			}),
		field.String("description").MaxLen(500).Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
