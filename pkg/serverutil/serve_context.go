package serverutil

import (
	"context"
)

// ServeContext encapsulates services and hooks for the Serve
type ServeContext struct {
	ctx  context.Context
	port string

	// services can be GRPCServer or GRPCGateway (or both)
	services []interface{}

	grpcEnabled    bool
	gatewayEnabled bool
	pubsubEnabled  bool

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

// WithGRPC enables GRPC capabilities of the server
func WithGRPC() ServeContextOption {
	return func(c *ServeContext) {
		c.grpcEnabled = true
	}
}

// WithGRPCGateway enables gateway capabilities of the server
func WithGRPCGateway() ServeContextOption {
	return func(c *ServeContext) {
		c.gatewayEnabled = true
	}
}

// WithPubSub enables PubSub capabilities of the server. Enabling PubSub
// also enables the gateway option
func WithPubSub() ServeContextOption {
	return func(c *ServeContext) {
		c.pubsubEnabled = true
		c.gatewayEnabled = true
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
