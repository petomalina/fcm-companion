package serverutil

import (
	"context"
	"github.com/petomalina/xrpc/pkg/multiplexer"
)

// ServeContext encapsulates services and hooks for the Serve
type ServeContext struct {
	ctx  context.Context
	port string

	// services can be GRPCServer or GRPCGateway (or both)
	services []interface{}

	// handlerFactories are handlers that should be created for
	// the services
	handlerFactories []multiplexer.HandlerFactory

	onListen func()
	onExit   func()
}

// ServeContextOption allows setup of the ServeContext
type ServeContextOption func(c *ServeContext)

// WithContext sets the ServeContext context. The context is then used within the
// grpc gateway connections.
func WithContext(ctx context.Context) ServeContextOption {
	return func(c *ServeContext) {
		c.ctx = ctx
	}
}

// WithPort sets the port of the context. If no port is set, one will be picked
// by the system
func WithPort(port string) ServeContextOption {
	return func(c *ServeContext) {
		c.port = port
	}
}

// WithServices adds services that should be registered within the multiplexer
func WithServices(s ...interface{}) ServeContextOption {
	return func(c *ServeContext) {
		c.services = append(c.services, s...)
	}
}

// WithHandlers adds new HandlerFactories. Serve and Serve apply all factories
// to all services that are eligible for that.
func WithHandlers(h ...multiplexer.HandlerFactory) ServeContextOption {
	return func(c *ServeContext) {
		c.handlerFactories = append(c.handlerFactories, h...)
	}
}

// WithOnListen adds the onListen callback fired when the listening starts
func WithOnListen(cb func()) ServeContextOption {
	return func(c *ServeContext) {
		c.onListen = cb
	}
}

// WithOnExit adds the onExit callback fired when the service is exiting
func WithOnExit(cb func()) ServeContextOption {
	return func(c *ServeContext) {
		c.onExit = cb
	}
}
