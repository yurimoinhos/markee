package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/aggi-tech/aggipay/platform/common"
)

// IntegrationWebhookEvent stores incoming provider webhooks for idempotent processing.
type IntegrationWebhookEvent struct {
	ent.Schema
}

func (IntegrationWebhookEvent) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(common.ID{}.Value).DefaultFunc(func() string { return common.GenID().Value }),
		field.String("provider").NotEmpty(),
		field.String("event_id").NotEmpty(),
		field.String("signature").Optional().Nillable().MaxLen(500),
		field.String("payload").NotEmpty(),
		field.String("status").Default("processed").Validate(func(s string) error {
			valid := map[string]bool{"processed": true, "ignored": true, "failed": true}
			if !valid[s] {
				return fmt.Errorf("status inválido: %q", s)
			}
			return nil
		}),
		field.Time("received_at").Default(time.Now),
		field.Time("processed_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (IntegrationWebhookEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider", "event_id").Unique(),
	}
}

func (IntegrationWebhookEvent) Annotations() []schema.Annotation {
	return nil
}
