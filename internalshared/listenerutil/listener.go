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
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/jefferai/isbadcipher"
	"github.com/mitchellh/cli"
)

type Listener struct {
	net.Listener
	Config *configutil.Listener
}

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

func TLSConfig(
	l *configutil.Listener,
	props map[string]string,
	ui cli.Ui) (*tls.Config, reloadutil.ReloadFunc, error) {
	props["tls"] = "disabled"

	if l.TLSDisable {
		return nil, nil, nil
	}

	cg := reloadutil.NewCertificateGetter(l.TLSCertFile, l.TLSKeyFile, "")
	if err := cg.Reload(); err != nil {
		// We try the key without a passphrase first and if we get an incorrect
		// passphrase response, try again after prompting for a passphrase
		if errwrap.Contains(err, x509.IncorrectPasswordError.Error()) {
			var passphrase string
			passphrase, err = ui.AskSecret(fmt.Sprintf("Enter passphrase for %s:", l.TLSKeyFile))
			if err == nil {
				cg = reloadutil.NewCertificateGetter(l.TLSCertFile, l.TLSKeyFile, passphrase)
				if err = cg.Reload(); err == nil {
					goto PASSPHRASECORRECT
				}
			}
		}
		return nil, nil, fmt.Errorf("error loading TLS cert: %w", err)
	}

PASSPHRASECORRECT:
	tlsConf := &tls.Config{
		GetCertificate:           cg.GetCertificate,
		NextProtos:               []string{"h2", "http/1.1"},
		ClientAuth:               tls.RequestClientCert,
		PreferServerCipherSuites: l.TLSPreferServerCipherSuites,
	}

	if l.TLSMinVersion == "" {
		l.TLSMinVersion = "tls12"
	}

	if l.TLSMaxVersion == "" {
		l.TLSMaxVersion = "tls13"
	}

	var ok bool
	tlsConf.MinVersion, ok = tlsutil.TLSLookup[l.TLSMinVersion]
	if !ok {
		return nil, nil, fmt.Errorf("'tls_min_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]", l.TLSMinVersion)
	}

	tlsConf.MaxVersion, ok = tlsutil.TLSLookup[l.TLSMaxVersion]
	if !ok {
		return nil, nil, fmt.Errorf("'tls_max_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]", l.TLSMaxVersion)
	}

	if tlsConf.MaxVersion < tlsConf.MinVersion {
		return nil, nil, fmt.Errorf("'tls_max_version' must be greater than or equal to 'tls_min_version'")
	}

	if len(l.TLSCipherSuites) > 0 {
		// HTTP/2 with TLS 1.2 blacklists several cipher suites.
		// https://tools.ietf.org/html/rfc7540#appendix-A
		//
		// Since the CLI (net/http) automatically uses HTTP/2 with TLS 1.2,
		// we check here if all or some specified cipher suites are blacklisted.
		badCiphers := []string{}
		for _, cipher := range l.TLSCipherSuites {
			if isbadcipher.IsBadCipher(cipher) {
				// Get the name of the current cipher.
				cipherStr, err := tlsutil.GetCipherName(cipher)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid value for 'tls_cipher_suites': %w", err)
				}
				badCiphers = append(badCiphers, cipherStr)
			}
		}
		if len(badCiphers) == len(l.TLSCipherSuites) {
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
		tlsConf.CipherSuites = l.TLSCipherSuites
	}

	if l.TLSRequireAndVerifyClientCert {
		tlsConf.ClientAuth = tls.RequireAndVerifyClientCert
		if l.TLSClientCAFile != "" {
			caPool := x509.NewCertPool()
			data, err := ioutil.ReadFile(l.TLSClientCAFile)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to read tls_client_ca_file: %w", err)
			}

			if !caPool.AppendCertsFromPEM(data) {
				return nil, nil, fmt.Errorf("failed to parse CA certificate in tls_client_ca_file")
			}
			tlsConf.ClientCAs = caPool
		}
	}

	if l.TLSDisableClientCerts {
		if l.TLSRequireAndVerifyClientCert {
			return nil, nil, fmt.Errorf("'tls_disable_client_certs' and 'tls_require_and_verify_client_cert' are mutually exclusive")
		}
		tlsConf.ClientAuth = tls.NoClientCert
	}

	props["tls"] = "enabled"
	return tlsConf, cg.Reload, nil
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
