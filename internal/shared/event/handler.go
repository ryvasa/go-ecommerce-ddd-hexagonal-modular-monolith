package event

import "context"

type Handler[E Event] interface {
	Handle(ctx context.Context, event E) error
}
