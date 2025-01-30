package transport

import "context"

// Transport defines the interface that all MCP transports must implement
type Transport interface {
	// Start starts the transport
	Start() error

	// SendNotification sends a notification to connected clients
	SendNotification(method string, params interface{}) error
}

// HTTPTransport extends Transport for HTTP-based transports
type HTTPTransport interface {
	Transport

	// StartHTTP starts the transport on the given address
	StartHTTP(addr string) error

	// Stop gracefully stops the transport
	Stop(ctx context.Context) error
}

// Options represents configuration options for transports
type Options struct {
	// Address is the network address to listen on (for HTTP transports)
	Address string

	// BufferSize is the size of notification channels
	BufferSize int

	// Additional options can be added here
}

// Option is a function that configures Options
type Option func(*Options)

// WithAddress sets the address option
func WithAddress(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// WithBufferSize sets the buffer size option
func WithBufferSize(size int) Option {
	return func(o *Options) {
		o.BufferSize = size
	}
}

// defaultOptions returns the default transport options
func defaultOptions() Options {
	return Options{
		Address:    ":8080",
		BufferSize: 100,
	}
}
