package structs

import (
	"fmt"
	"io"
	"sync"
)

// StreamingRpcHeader is the first struct serialized after entering the
// streaming RPC mode. The header is used to dispatch to the correct method.
type StreamingRpcHeader struct {
	// Method is the name of the method to invoke.
	Method string
}

// StreamingRpcAck is used to acknowledge receiving the StreamingRpcHeader and
// routing to the requested handler.
type StreamingRpcAck struct {
	// Error is used to return whether an error occurred establishing the
	// streaming RPC. This error occurs before entering the RPC handler.
	Error string
}

// StreamingRpcHandler defines the handler for a streaming RPC.
type StreamingRpcHandler func(conn io.ReadWriteCloser)

// StreamingRpcRegistry is used to add and retrieve handlers
type StreamingRpcRegistry struct {
	registry map[string]StreamingRpcHandler
}

// NewStreamingRpcRegistry creates a new registry. All registrations of
// handlers should be done before retrieving handlers.
func NewStreamingRpcRegistry() *StreamingRpcRegistry {
	return &StreamingRpcRegistry{
		registry: make(map[string]StreamingRpcHandler),
	}
}

// Register registers a new handler for the given method name
func (s *StreamingRpcRegistry) Register(method string, handler StreamingRpcHandler) {
	s.registry[method] = handler
}

// GetHandler returns a handler for the given method or an error if it doesn't exist.
func (s *StreamingRpcRegistry) GetHandler(method string) (StreamingRpcHandler, error) {
	h, ok := s.registry[method]
	if !ok {
		return nil, fmt.Errorf("%s: %q", ErrUnknownMethod, method)
	}

	return h, nil
}

// Bridge is used to just link two connections together and copy traffic
func Bridge(a, b io.ReadWriteCloser) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(a, b)
		a.Close()
		b.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(b, a)
		a.Close()
		b.Close()
	}()
	wg.Wait()
}
