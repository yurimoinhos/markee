package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Customer stores customer financial profile and preferences.
type Customer struct {
	ent.Schema
}

func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("owner_user_id").NotEmpty(),
		field.String("name").NotEmpty().MaxLen(120),
		field.String("company").Optional().Nillable().MaxLen(160),
		field.String("cpf_cnpj").NotEmpty().MaxLen(20),
		field.String("email").NotEmpty().MaxLen(120),
		field.String("phone").Optional().Nillable().MaxLen(30),
		field.String("address").Optional().Nillable().MaxLen(300),
		field.String("preferred_payment_method").Default("pix").Validate(func(s string) error {
			valid := map[string]bool{"pix": true, "bank_transfer": true, "credit_card": true, "crypto": true, "boleto": true}
			if !valid[s] {
				return fmt.Errorf("preferred_payment_method inválido: %q", s)
			}
			return nil
		}),
		field.String("financial_status").Default("regular").Validate(func(s string) error {
			valid := map[string]bool{"regular": true, "watch": true, "delinquent": true}
			if !valid[s] {
				return fmt.Errorf("financial_status inválido: %q", s)
			}
			return nil
		}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
