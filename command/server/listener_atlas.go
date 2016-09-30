package server

import (
	"io"
	"net"

	"github.com/hashicorp/scada-client/scada"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

type SCADAListener struct {
	ln            net.Listener
	scadaProvider *scada.Provider
}

func (s *SCADAListener) Accept() (net.Conn, error) {
	return s.ln.Accept()
}

func (s *SCADAListener) Close() error {
	s.scadaProvider.Shutdown()
	return s.ln.Close()
}

func (s *SCADAListener) Addr() net.Addr {
	return s.ln.Addr()
}

func atlasListenerFactory(config map[string]string, logger io.Writer) (net.Listener, map[string]string, vault.ReloadFunc, error) {
	scadaConfig := &scada.Config{
		Service:      "vault",
		Version:      version.GetVersion().VersionNumber(),
		ResourceType: "vault-cluster",
		Meta: map[string]string{
			"node_id":      config["node_id"],
			"cluster_name": config["cluster_name"],
		},
		Atlas: scada.AtlasConfig{
			Endpoint:       config["endpoint"],
			Infrastructure: config["infrastructure"],
			Token:          config["token"],
		},
	}

	provider, list, err := scada.NewHTTPProvider(scadaConfig, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	ln := &SCADAListener{
		ln:            list,
		scadaProvider: provider,
	}

	props := map[string]string{
		"addr":           "Atlas/SCADA",
		"infrastructure": scadaConfig.Atlas.Infrastructure,
	}

	// The outer connection is already TLS-enabled; this is just the listener
	// that reaches back inside that connection
	config["tls_disable"] = "1"

	return listenerWrapTLS(ln, props, config)
}
