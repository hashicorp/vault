package server

import (
	// We must import sha512 so that it registers with the runtime so that
	// certificates that use it can be parsed.
	_ "crypto/sha512"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

	"github.com/hashicorp/vault/helper/tlsutil"
)

// ListenerFactory is the factory function to create a listener.
type ListenerFactory func(map[string]string, io.Writer) (net.Listener, map[string]string, ReloadFunc, error)

// BuiltinListeners is the list of built-in listener types.
var BuiltinListeners = map[string]ListenerFactory{
	"tcp":   tcpListenerFactory,
	"atlas": atlasListenerFactory,
}

// NewListener creates a new listener of the given type with the given
// configuration. The type is looked up in the BuiltinListeners map.
func NewListener(t string, config map[string]string, logger io.Writer) (net.Listener, map[string]string, ReloadFunc, error) {
	f, ok := BuiltinListeners[t]
	if !ok {
		return nil, nil, nil, fmt.Errorf("unknown listener type: %s", t)
	}

	return f(config, logger)
}

func listenerWrapTLS(
	ln net.Listener,
	props map[string]string,
	config map[string]string) (net.Listener, map[string]string, ReloadFunc, error) {
	props["tls"] = "disabled"

	if v, ok := config["tls_disable"]; ok {
		disabled, err := strconv.ParseBool(v)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("invalid value for 'tls_disable': %v", err)
		}
		if disabled {
			return ln, props, nil, nil
		}
	}

	_, ok := config["tls_cert_file"]
	if !ok {
		return nil, nil, nil, fmt.Errorf("'tls_cert_file' must be set")
	}

	_, ok = config["tls_key_file"]
	if !ok {
		return nil, nil, nil, fmt.Errorf("'tls_key_file' must be set")
	}

	cg := &certificateGetter{
		id: config["address"],
	}

	if err := cg.reload(config); err != nil {
		return nil, nil, nil, fmt.Errorf("error loading TLS cert: %s", err)
	}

	tlsvers, ok := config["tls_min_version"]
	if !ok {
		tlsvers = "tls12"
	}

	tlsConf := &tls.Config{}
	tlsConf.GetCertificate = cg.getCertificate
	tlsConf.NextProtos = []string{"http/1.1"}
	tlsConf.MinVersion, ok = tlsutil.TLSLookup[tlsvers]
	if !ok {
		return nil, nil, nil, fmt.Errorf("'tls_min_version' value %s not supported, please specify one of [tls10,tls11,tls12]", tlsvers)
	}
	tlsConf.ClientAuth = tls.RequestClientCert

	ln = tls.NewListener(ln, tlsConf)
	props["tls"] = "enabled"
	return ln, props, cg.reload, nil
}

type certificateGetter struct {
	sync.RWMutex

	cert *tls.Certificate

	id string
}

func (cg *certificateGetter) reload(config map[string]string) error {
	if config["address"] != cg.id {
		return nil
	}

	cert, err := tls.LoadX509KeyPair(config["tls_cert_file"], config["tls_key_file"])
	if err != nil {
		return err
	}

	cg.Lock()
	defer cg.Unlock()

	cg.cert = &cert

	return nil
}

func (cg *certificateGetter) getCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cg.RLock()
	defer cg.RUnlock()

	if cg.cert == nil {
		return nil, fmt.Errorf("nil certificate")
	}

	return cg.cert, nil
}
