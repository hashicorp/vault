package radius

import (
	"context"
	"errors"
	"net"
)

// ErrServerShutdown is returned from server Serve methods when Shutdown
// has been called and handlers are still completing.
var ErrServerShutdown = errors.New("radius: server is shutting down")

// Handler provides a handler to RADIUS server requests. When a RADIUS request
// is received, ServeRADIUS is called.
type Handler interface {
	ServeRADIUS(w ResponseWriter, r *Request)
}

// HandlerFunc allows a function to implement Handler.
type HandlerFunc func(w ResponseWriter, r *Request)

// ServeRADIUS calls h(w, p).
func (h HandlerFunc) ServeRADIUS(w ResponseWriter, r *Request) {
	h(w, r)
}

// Request is an incoming RADIUS request that is being handled by the server.
type Request struct {
	LocalAddr  net.Addr
	RemoteAddr net.Addr

	*Packet

	ctx context.Context
}

// Context returns the context of the request. If a context has not been set
// using WithContext, the Background context is returned.
func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

// WithContext returns a shallow copy of the request with the new request's
// context set to the given context.
func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil ctx")
	}
	req := new(Request)
	*req = *r
	req.ctx = ctx
	return req
}

// ResponseWriter is used by RADIUS servers when replying to a RADIUS request.
type ResponseWriter interface {
	Write(packet *Packet) error
}

// SecretSource supplies RADIUS servers with the secret that should be used for
// authorizing and decrypting packets.
//
// ctx is canceled if the server's Shutdown method is called.
type SecretSource interface {
	RADIUSSecret(ctx context.Context, remoteAddr net.Addr) ([]byte, error)
}

// StaticSecretSource returns a SecretSource that uses secret for all requests.
func StaticSecretSource(secret []byte) SecretSource {
	return staticSecretSource(secret)
}

type staticSecretSource []byte

func (secret staticSecretSource) RADIUSSecret(ctx context.Context, remoteAddr net.Addr) ([]byte, error) {
	return []byte(secret), nil
}
