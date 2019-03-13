package listenerutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	osuser "os/user"
	"strconv"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/jefferai/isbadcipher"
	"github.com/mitchellh/cli"
)

type UnixSocketsConfig struct {
	User  string `hcl:"user"`
	Mode  string `hcl:"mode"`
	Group string `hcl:"group"`
}

// rmListener is an implementation of net.Listener that forwards most
// calls to the listener but also removes a file as part of the close. We
// use this to cleanup the unix domain socket on close.
type rmListener struct {
	net.Listener
	Path string
}

func (l *rmListener) Close() error {
	// Close the listener itself
	if err := l.Listener.Close(); err != nil {
		return err
	}

	// Remove the file
	return os.Remove(l.Path)
}

func UnixSocketListener(path string, unixSocketsConfig *UnixSocketsConfig) (net.Listener, error) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove socket file: %v", err)
	}

	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	if unixSocketsConfig != nil {
		err = setFilePermissions(path, unixSocketsConfig.User, unixSocketsConfig.Group, unixSocketsConfig.Mode)
		if err != nil {
			return nil, fmt.Errorf("failed to set file system permissions on the socket file: %s", err)
		}
	}

	// Wrap the listener in rmListener so that the Unix domain socket file is
	// removed on close.
	return &rmListener{
		Listener: ln,
		Path:     path,
	}, nil
}

func WrapTLS(
	ln net.Listener,
	props map[string]string,
	config map[string]interface{},
	ui cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, *tls.Config, error) {
	props["tls"] = "disabled"

	if v, ok := config["tls_disable"]; ok {
		disabled, err := parseutil.ParseBool(v)
		if err != nil {
			return nil, nil, nil, nil, errwrap.Wrapf("invalid value for 'tls_disable': {{err}}", err)
		}
		if disabled {
			return ln, props, nil, nil, nil
		}
	}

	certFileRaw, ok := config["tls_cert_file"]
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("'tls_cert_file' must be set")
	}
	certFile := certFileRaw.(string)
	keyFileRaw, ok := config["tls_key_file"]
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("'tls_key_file' must be set")
	}
	keyFile := keyFileRaw.(string)

	cg := reload.NewCertificateGetter(certFile, keyFile, "")
	if err := cg.Reload(config); err != nil {
		// We try the key without a passphrase first and if we get an incorrect
		// passphrase response, try again after prompting for a passphrase
		if errwrap.Contains(err, x509.IncorrectPasswordError.Error()) {
			var passphrase string
			passphrase, err = ui.AskSecret(fmt.Sprintf("Enter passphrase for %s:", keyFile))
			if err == nil {
				cg = reload.NewCertificateGetter(certFile, keyFile, passphrase)
				if err = cg.Reload(config); err == nil {
					goto PASSPHRASECORRECT
				}
			}
		}
		return nil, nil, nil, nil, errwrap.Wrapf("error loading TLS cert: {{err}}", err)
	}

PASSPHRASECORRECT:
	var tlsvers string
	tlsversRaw, ok := config["tls_min_version"]
	if !ok {
		tlsvers = "tls12"
	} else {
		tlsvers = tlsversRaw.(string)
	}

	tlsConf := &tls.Config{}
	tlsConf.GetCertificate = cg.GetCertificate
	tlsConf.NextProtos = []string{"h2", "http/1.1"}
	tlsConf.MinVersion, ok = tlsutil.TLSLookup[tlsvers]
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("'tls_min_version' value %q not supported, please specify one of [tls10,tls11,tls12]", tlsvers)
	}
	tlsConf.ClientAuth = tls.RequestClientCert

	if v, ok := config["tls_cipher_suites"]; ok {
		ciphers, err := tlsutil.ParseCiphers(v.(string))
		if err != nil {
			return nil, nil, nil, nil, errwrap.Wrapf("invalid value for 'tls_cipher_suites': {{err}}", err)
		}

		// HTTP/2 with TLS 1.2 blacklists several cipher suites.
		// https://tools.ietf.org/html/rfc7540#appendix-A
		//
		// Since the CLI (net/http) automatically uses HTTP/2 with TLS 1.2,
		// we check here if all or some specified cipher suites are blacklisted.
		badCiphers := []string{}
		for _, cipher := range ciphers {
			if isbadcipher.IsBadCipher(cipher) {
				// Get the name of the current cipher.
				cipherStr, err := tlsutil.GetCipherName(cipher)
				if err != nil {
					return nil, nil, nil, nil, errwrap.Wrapf("invalid value for 'tls_cipher_suites': {{err}}", err)
				}
				badCiphers = append(badCiphers, cipherStr)
			}
		}
		if len(badCiphers) == len(ciphers) {
			ui.Warn(`WARNING! All cipher suites defined by 'tls_cipher_suites' are blacklisted by the
HTTP/2 specification. HTTP/2 communication with TLS 1.2 will not work as intended
and Vault will be unavailable via the CLI.
Please see https://tools.ietf.org/html/rfc7540#appendix-A for further information.`)
		} else if len(badCiphers) > 0 {
			ui.Warn(fmt.Sprintf(`WARNING! The following cipher suites defined by 'tls_cipher_suites' are
blacklisted by the HTTP/2 specification:
%v
Please see https://tools.ietf.org/html/rfc7540#appendix-A for further information.`, badCiphers))
		}
		tlsConf.CipherSuites = ciphers
	}
	if v, ok := config["tls_prefer_server_cipher_suites"]; ok {
		preferServer, err := parseutil.ParseBool(v)
		if err != nil {
			return nil, nil, nil, nil, errwrap.Wrapf("invalid value for 'tls_prefer_server_cipher_suites': {{err}}", err)
		}
		tlsConf.PreferServerCipherSuites = preferServer
	}
	var requireVerifyCerts bool
	var err error
	if v, ok := config["tls_require_and_verify_client_cert"]; ok {
		requireVerifyCerts, err = parseutil.ParseBool(v)
		if err != nil {
			return nil, nil, nil, nil, errwrap.Wrapf("invalid value for 'tls_require_and_verify_client_cert': {{err}}", err)
		}
		if requireVerifyCerts {
			tlsConf.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if tlsClientCaFile, ok := config["tls_client_ca_file"]; ok {
			caPool := x509.NewCertPool()
			data, err := ioutil.ReadFile(tlsClientCaFile.(string))
			if err != nil {
				return nil, nil, nil, nil, errwrap.Wrapf("failed to read tls_client_ca_file: {{err}}", err)
			}

			if !caPool.AppendCertsFromPEM(data) {
				return nil, nil, nil, nil, fmt.Errorf("failed to parse CA certificate in tls_client_ca_file")
			}
			tlsConf.ClientCAs = caPool
		}
	}
	if v, ok := config["tls_disable_client_certs"]; ok {
		disableClientCerts, err := parseutil.ParseBool(v)
		if err != nil {
			return nil, nil, nil, nil, errwrap.Wrapf("invalid value for 'tls_disable_client_certs': {{err}}", err)
		}
		if disableClientCerts && requireVerifyCerts {
			return nil, nil, nil, nil, fmt.Errorf("'tls_disable_client_certs' and 'tls_require_and_verify_client_cert' are mutually exclusive")
		}
		if disableClientCerts {
			tlsConf.ClientAuth = tls.NoClientCert
		}
	}

	ln = tls.NewListener(ln, tlsConf)
	props["tls"] = "enabled"
	return ln, props, cg.Reload, tlsConf, nil
}

// setFilePermissions handles configuring ownership and permissions
// settings on a given file. All permission/ownership settings are
// optional. If no user or group is specified, the current user/group
// will be used. Mode is optional, and has no default (the operation is
// not performed if absent). User may be specified by name or ID, but
// group may only be specified by ID.
func setFilePermissions(path string, user, group, mode string) error {
	var err error
	uid, gid := os.Getuid(), os.Getgid()

	if user != "" {
		if uid, err = strconv.Atoi(user); err == nil {
			goto GROUP
		}

		// Try looking up the user by name
		u, err := osuser.Lookup(user)
		if err != nil {
			return fmt.Errorf("failed to look up user %q: %v", user, err)
		}
		uid, _ = strconv.Atoi(u.Uid)
	}

GROUP:
	if group != "" {
		if gid, err = strconv.Atoi(group); err == nil {
			goto OWN
		}

		// Try looking up the user by name
		g, err := osuser.LookupGroup(group)
		if err != nil {
			return fmt.Errorf("failed to look up group %q: %v", user, err)
		}
		gid, _ = strconv.Atoi(g.Gid)
	}

OWN:
	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("failed setting ownership to %d:%d on %q: %v",
			uid, gid, path, err)
	}

	if mode != "" {
		mode, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return fmt.Errorf("invalid mode specified: %v", mode)
		}
		if err := os.Chmod(path, os.FileMode(mode)); err != nil {
			return fmt.Errorf("failed setting permissions to %d on %q: %v",
				mode, path, err)
		}
	}

	return nil
}
