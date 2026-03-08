package http

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aggi-tech/aggipay/contracts"
)

// RegisterConsumers registra os handlers de eventos que o módulo order consome.
func RegisterConsumers(consumer contracts.EventConsumer, svc Service) {
	consumer.Subscribe(contracts.TopicOrderConfirmed, func(ctx context.Context, raw []byte) error {
		var evt contracts.OrderConfirmed
		if err := json.Unmarshal(raw, &evt); err != nil {
			return err
		}
		log.Printf("[order] OrderConfirmed recebido: order=%s saga=%s provider=%s", evt.OrderID, evt.SagaID, evt.ProviderRef)
		return svc.ConfirmOrder(ctx, evt.OrderID, evt.ProviderRef)
	})

	consumer.Subscribe(contracts.TopicOrderCancelled, func(ctx context.Context, raw []byte) error {
		var evt contracts.OrderCancelled
		if err := json.Unmarshal(raw, &evt); err != nil {
			return err
		}
		log.Printf("[order] OrderCancelled recebido: order=%s saga=%s reason=%s", evt.OrderID, evt.SagaID, evt.Reason)
		return svc.CancelOrder(ctx, evt.OrderID, evt.Reason)
	})
}

// Service é o subconjunto de application.Service necessário para os consumers.
type Service interface {
	ConfirmOrder(ctx context.Context, id, providerRef string) error
	CancelOrder(ctx context.Context, id, reason string) error
}
