package contracts

import "context"

// EventPublisher é a interface que todo módulo usa para publicar eventos no barramento.
// A implementação concreta fica em platform/eventbus/rabbitmq.
type EventPublisher interface {
	// Publish serializa event como JSON e o publica no tópico informado.
	Publish(ctx context.Context, topic string, event any) error
}

// HandlerFunc é a assinatura do handler que processa uma mensagem bruta do barramento.
type HandlerFunc func(ctx context.Context, raw []byte) error

// EventConsumer é a interface que todo módulo usa para consumir eventos do barramento.
type EventConsumer interface {
	// Subscribe registra um handler para o tópico informado.
	// Deve ser chamado antes de Start.
	Subscribe(topic string, handler HandlerFunc)

	// Start inicia o loop de consumo. Bloqueia até ctx ser cancelado.
	Start(ctx context.Context) error
}
