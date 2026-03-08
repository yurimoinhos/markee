package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// ContractSignature stores digital acceptance audit data.
type ContractSignature struct {
	ent.Schema
}

func (ContractSignature) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("contract_id").NotEmpty(),
		field.String("contract_version_id").Optional().Nillable(),
		field.String("provider").NotEmpty().Default("clicksign"),
		field.String("provider_doc_id").Optional().Nillable(),
		field.String("status").Default("pending").Validate(func(s string) error {
			valid := map[string]bool{"pending": true, "signed": true, "rejected": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.String("signer_name").Optional().Nillable(),
		field.String("signer_email").Optional().Nillable(),
		field.Time("accepted_at").Optional().Nillable(),
		field.String("accepted_ip").Optional().Nillable().MaxLen(80),
		field.String("payload_json").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
