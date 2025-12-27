package event

import (
	"context"
	"log"

	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
)

type InMemoryPublisher struct {
	handlers map[string][]func(context.Context, sharedEvent.Event) error
}

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{
		handlers: make(map[string][]func(context.Context, sharedEvent.Event) error),
	}
}

func (p *InMemoryPublisher) Register(
	eventName string,
	handler func(context.Context, sharedEvent.Event) error,
) {
	p.handlers[eventName] = append(p.handlers[eventName], handler)
}

func (p *InMemoryPublisher) Publish(
	ctx context.Context,
	events ...sharedEvent.Event,
) error {
	for _, e := range events {
		if hs, ok := p.handlers[e.Name()]; ok {
			for _, h := range hs {
				go func(handler func(context.Context, sharedEvent.Event) error, ev sharedEvent.Event) {
					if err := handler(ctx, ev); err != nil {
						log.Println("event handler error:", err)
					}
				}(h, e)
			}
		}
	}
	return nil
}
