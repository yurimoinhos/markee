package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/aggi-tech/aggipay/platform/common"
)

// ContractVersion stores generated contract artifacts.
type ContractVersion struct {
	ent.Schema
}

func (ContractVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("contract_id").NotEmpty(),
		field.Int("version").Min(1),
		field.String("template_name").NotEmpty().MaxLen(120),
		field.String("editable_content").NotEmpty(),
		field.String("pdf_url").Optional().Nillable().MaxLen(500),
		field.String("editable_url").Optional().Nillable().MaxLen(500),
		field.Time("generated_at").Default(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
