package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

var (
	clusterTLSServerLookup = func(ctx context.Context, c *Core, repClusters *ReplicatedClusters, _ *ReplicatedCluster) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			c.logger.Debug("performing server cert lookup")

			switch {
			default:
				currCert := c.localClusterCert.Load().([]byte)
				if len(currCert) == 0 {
					return nil, fmt.Errorf("got forwarding connection but no local cert")
				}

				localCert := make([]byte, len(currCert))
				copy(localCert, currCert)

				return &tls.Certificate{
					Certificate: [][]byte{localCert},
					PrivateKey:  c.localClusterPrivateKey.Load().(*ecdsa.PrivateKey),
					Leaf:        c.localClusterParsedCert.Load().(*x509.Certificate),
				}, nil
			}
		}
	}

	clusterTLSClientLookup = func(ctx context.Context, c *Core, repClusters *ReplicatedClusters, _ *ReplicatedCluster) func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return func(requestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if len(requestInfo.AcceptableCAs) != 1 {
				return nil, fmt.Errorf("expected only a single acceptable CA")
			}

			currCert := c.localClusterCert.Load().([]byte)
			if len(currCert) == 0 {
				return nil, fmt.Errorf("forwarding connection client but no local cert")
			}

			localCert := make([]byte, len(currCert))
			copy(localCert, currCert)

			return &tls.Certificate{
				Certificate: [][]byte{localCert},
				PrivateKey:  c.localClusterPrivateKey.Load().(*ecdsa.PrivateKey),
				Leaf:        c.localClusterParsedCert.Load().(*x509.Certificate),
			}, nil
		}
	}

	clusterTLSServerConfigLookup = func(ctx context.Context, c *Core, repClusters *ReplicatedClusters, repCluster *ReplicatedCluster) func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
		return func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
			//c.logger.Trace("performing server config lookup")

			caPool := x509.NewCertPool()

			ret := &tls.Config{
				ClientAuth:           tls.RequireAndVerifyClientCert,
				GetCertificate:       clusterTLSServerLookup(ctx, c, repClusters, repCluster),
				GetClientCertificate: clusterTLSClientLookup(ctx, c, repClusters, repCluster),
				MinVersion:           tls.VersionTLS12,
				RootCAs:              caPool,
				ClientCAs:            caPool,
				NextProtos:           clientHello.SupportedProtos,
				CipherSuites:         c.clusterCipherSuites,
			}

			parsedCert := c.localClusterParsedCert.Load().(*x509.Certificate)

			if parsedCert == nil {
				return nil, fmt.Errorf("forwarding connection client but no local cert")
			}

			caPool.AddCert(parsedCert)

			return ret, nil
		}
	}
)
