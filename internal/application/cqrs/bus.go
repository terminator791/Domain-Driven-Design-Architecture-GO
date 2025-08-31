package cqrs

import (
	"context"
	"fmt"
	"time"
)

// Command represents a command in CQRS pattern
type Command interface {
	CommandID() string
	CommandType() string
	AggregateID() string
	GetPayload() interface{}
}

// CommandHandler handles commands
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
	CanHandle(cmd Command) bool
}

// Query represents a query in CQRS pattern
type Query interface {
	QueryID() string
	QueryType() string
	GetFilter() interface{}
}

// QueryHandler handles queries
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
	CanHandle(query Query) bool
}

// CommandBus dispatches commands to appropriate handlers
type CommandBus interface {
	Dispatch(ctx context.Context, cmd Command) error
	RegisterHandler(handler CommandHandler)
}

// QueryBus dispatches queries to appropriate handlers
type QueryBus interface {
	Dispatch(ctx context.Context, query Query) (interface{}, error)
	RegisterHandler(handler QueryHandler)
}

// Event represents a domain event
type Event interface {
	EventID() string
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	GetPayload() interface{}
}

// EventHandler handles domain events
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	CanHandle(event Event) bool
}

// EventBus publishes events to handlers
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(handler EventHandler)
}

// SimpleCommandBus is a basic implementation of CommandBus
type SimpleCommandBus struct {
	handlers []CommandHandler
}

func NewSimpleCommandBus() *SimpleCommandBus {
	return &SimpleCommandBus{
		handlers: make([]CommandHandler, 0),
	}
}

func (bus *SimpleCommandBus) RegisterHandler(handler CommandHandler) {
	bus.handlers = append(bus.handlers, handler)
}

func (bus *SimpleCommandBus) Dispatch(ctx context.Context, cmd Command) error {
	for _, handler := range bus.handlers {
		if handler.CanHandle(cmd) {
			return handler.Handle(ctx, cmd)
		}
	}
	return fmt.Errorf("no handler found for command type: %s", cmd.CommandType())
}

// SimpleQueryBus is a basic implementation of QueryBus
type SimpleQueryBus struct {
	handlers []QueryHandler
}

func NewSimpleQueryBus() *SimpleQueryBus {
	return &SimpleQueryBus{
		handlers: make([]QueryHandler, 0),
	}
}

func (bus *SimpleQueryBus) RegisterHandler(handler QueryHandler) {
	bus.handlers = append(bus.handlers, handler)
}

func (bus *SimpleQueryBus) Dispatch(ctx context.Context, query Query) (interface{}, error) {
	for _, handler := range bus.handlers {
		if handler.CanHandle(query) {
			return handler.Handle(ctx, query)
		}
	}
	return nil, fmt.Errorf("no handler found for query type: %s", query.QueryType())
}

// SimpleEventBus is a basic implementation of EventBus
type SimpleEventBus struct {
	handlers []EventHandler
}

func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make([]EventHandler, 0),
	}
}

func (bus *SimpleEventBus) Subscribe(handler EventHandler) {
	bus.handlers = append(bus.handlers, handler)
}

func (bus *SimpleEventBus) Publish(ctx context.Context, event Event) error {
	for _, handler := range bus.handlers {
		if handler.CanHandle(event) {
			// In a real implementation, this might be async
			err := handler.Handle(ctx, event)
			if err != nil {
				// Log error but continue with other handlers
				fmt.Printf("Error handling event %s: %v\n", event.EventType(), err)
			}
		}
	}
	return nil
}