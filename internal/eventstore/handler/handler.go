package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
)

type Handler struct {
	Eventstore *eventstore.Eventstore
	Sub        *eventstore.Subscription
	EventQueue chan eventstore.EventReader
}

func NewHandler(es *eventstore.Eventstore) Handler {
	return Handler{
		Eventstore: es,
		EventQueue: make(chan eventstore.EventReader, 100),
	}
}

func (h Handler) Subscribe(aggregates ...eventstore.AggregateType) {
	h.Sub = eventstore.Subscribe(h.EventQueue, aggregates...)
}
