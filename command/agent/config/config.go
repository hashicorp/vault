package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/parseutil"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Config is the configuration for the vault server.
type Config struct {
	AutoAuth      *AutoAuth   `hcl:"auto_auth"`
	ExitAfterAuth bool        `hcl:"exit_after_auth"`
	PidFile       string      `hcl:"pid_file"`
	Listeners     []*Listener `hcl:"listeners"`
	Cache         *Cache      `hcl:"cache"`
	Vault         *Vault      `hcl:"vault"`
}

type Vault struct {
	Address       string `hcl:"address"`
	CACert        string `hcl:"ca_cert"`
	CAPath        string `hcl:"ca_path"`
	TLSSkipVerify bool   `hcl:"tls_skip_verify"`
	ClientCert    string `hcl:"client_cert"`
	ClientKey     string `hcl:"client_key"`
}

type Cache struct {
	UseAutoAuthToken bool `hcl:"use_auto_auth_token"`
}

type Listener struct {
	Type   string
	Config map[string]interface{}
}

type AutoAuth struct {
	Method *Method `hcl:"-"`
	Sinks  []*Sink `hcl:"sinks"`

	// NOTE: This is unsupported outside of testing and may disappear at any
	// time.
	EnableReauthOnNewCredentials bool `hcl:"enable_reauth_on_new_credentials"`
}

type Method struct {
	Type       string
	MountPath  string        `hcl:"mount_path"`
	WrapTTLRaw interface{}   `hcl:"wrap_ttl"`
	WrapTTL    time.Duration `hcl:"-"`
	Config     map[string]interface{}
}

type Sink struct {
	Type       string
	WrapTTLRaw interface{}   `hcl:"wrap_ttl"`
	WrapTTL    time.Duration `hcl:"-"`
	DHType     string        `hcl:"dh_type"`
	DHPath     string        `hcl:"dh_path"`
	AAD        string        `hcl:"aad"`
	AADEnvVar  string        `hcl:"aad_env_var"`
	Config     map[string]interface{}
}

// LoadConfig loads the configuration at the given path, regardless if
// its a file or directory.
func LoadConfig(path string, logger log.Logger) (*Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("location is a directory, not a file")
	}

	// Read the file
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse!
	obj, err := hcl.Parse(string(d))
	if err != nil {
		return nil, err
	}

	// Start building the result
	var result Config
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	if err := parseAutoAuth(&result, list); err != nil {
		return nil, errwrap.Wrapf("error parsing 'auto_auth': {{err}}", err)
	}

	err = parseListeners(&result, list)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing 'listeners': {{err}}", err)
	}

	err = parseCache(&result, list)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing 'cache':{{err}}", err)
	}

	if result.Cache != nil {
		if len(result.Listeners) < 1 {
			return nil, fmt.Errorf("at least one listener required when cache enabled")
		}

		if result.Cache.UseAutoAuthToken && result.AutoAuth == nil {
			return nil, fmt.Errorf("cache.use_auto_auth_token is true but auto_auth not configured")
		}
	}

	err = parseVault(&result, list)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing 'vault':{{err}}", err)
	}

	return &result, nil
}

func parseVault(result *Config, list *ast.ObjectList) error {
	name := "vault"

	vaultList := list.Filter(name)
	if len(vaultList.Items) == 0 {
		return nil
	}

	if len(vaultList.Items) > 1 {
		return fmt.Errorf("one and only one %q block is required", name)
	}

	item := vaultList.Items[0]

	var v Vault
	err := hcl.DecodeObject(&v, item.Val)
	if err != nil {
		return err
	}

	result.Vault = &v

	return nil
}

func parseCache(result *Config, list *ast.ObjectList) error {
	name := "cache"

	cacheList := list.Filter(name)
	if len(cacheList.Items) == 0 {
		return nil
	}

	if len(cacheList.Items) > 1 {
		return fmt.Errorf("one and only one %q block is required", name)
	}

	item := cacheList.Items[0]

	var c Cache
	err := hcl.DecodeObject(&c, item.Val)
	if err != nil {
		return err
	}

	result.Cache = &c
	return nil
}

func parseListeners(result *Config, list *ast.ObjectList) error {
	name := "listener"

	listenerList := list.Filter(name)

	var listeners []*Listener
	for _, item := range listenerList.Items {
		var lnConfig map[string]interface{}
		err := hcl.DecodeObject(&lnConfig, item.Val)
		if err != nil {
			return err
		}

		var lnType string
		switch {
		case lnConfig["type"] != nil:
			lnType = lnConfig["type"].(string)
			delete(lnConfig, "type")
		case len(item.Keys) == 1:
			lnType = strings.ToLower(item.Keys[0].Token.Value().(string))
		default:
			return errors.New("listener type must be specified")
		}

		switch lnType {
		case "unix", "tcp":
		default:
			return fmt.Errorf("invalid listener type %q", lnType)
		}

		listeners = append(listeners, &Listener{
			Type:   lnType,
			Config: lnConfig,
		})
	}

	result.Listeners = listeners

	return nil
}

func parseAutoAuth(result *Config, list *ast.ObjectList) error {
	name := "auto_auth"

	autoAuthList := list.Filter(name)
	if len(autoAuthList.Items) == 0 {
		return nil
	}
	if len(autoAuthList.Items) > 1 {
		return fmt.Errorf("at most one %q block is allowed", name)
	}

	// Get our item
	item := autoAuthList.Items[0]

	var a AutoAuth
	if err := hcl.DecodeObject(&a, item.Val); err != nil {
		return err
	}

	result.AutoAuth = &a

	subs, ok := item.Val.(*ast.ObjectType)
	if !ok {
		return fmt.Errorf("could not parse %q as an object", name)
	}
	subList := subs.List

	if err := parseMethod(result, subList); err != nil {
		return errwrap.Wrapf("error parsing 'method': {{err}}", err)
	}
	if a.Method == nil {
		return fmt.Errorf("no 'method' block found")
	}

	if err := parseSinks(result, subList); err != nil {
		return errwrap.Wrapf("error parsing 'sink' stanzas: {{err}}", err)
	}

	return nil
}

func parseMethod(result *Config, list *ast.ObjectList) error {
	name := "method"

	methodList := list.Filter(name)
	if len(methodList.Items) != 1 {
		return fmt.Errorf("one and only one %q block is required", name)
	}

	// Get our item
	item := methodList.Items[0]

	var m Method
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	if m.Type == "" {
		if len(item.Keys) == 1 {
			m.Type = strings.ToLower(item.Keys[0].Token.Value().(string))
		}
		if m.Type == "" {
			return errors.New("method type must be specified")
		}
	}

	// Default to Vault's default
	if m.MountPath == "" {
		m.MountPath = fmt.Sprintf("auth/%s", m.Type)
	}
	// Standardize on no trailing slash
	m.MountPath = strings.TrimSuffix(m.MountPath, "/")

	if m.WrapTTLRaw != nil {
		var err error
		if m.WrapTTL, err = parseutil.ParseDurationSecond(m.WrapTTLRaw); err != nil {
			return err
		}
		m.WrapTTLRaw = nil
	}

	result.AutoAuth.Method = &m
	return nil
}

func parseSinks(result *Config, list *ast.ObjectList) error {
	name := "sink"

	sinkList := list.Filter(name)
	if len(sinkList.Items) < 1 {
		return nil
	}

	var ts []*Sink

	for _, item := range sinkList.Items {
		var s Sink
		if err := hcl.DecodeObject(&s, item.Val); err != nil {
			return err
		}

		if s.Type == "" {
			if len(item.Keys) == 1 {
				s.Type = strings.ToLower(item.Keys[0].Token.Value().(string))
			}
			if s.Type == "" {
				return errors.New("sink type must be specified")
			}
		}

		if s.WrapTTLRaw != nil {
			var err error
			if s.WrapTTL, err = parseutil.ParseDurationSecond(s.WrapTTLRaw); err != nil {
				return multierror.Prefix(err, fmt.Sprintf("sink.%s", s.Type))
			}
			s.WrapTTLRaw = nil
		}

		switch s.DHType {
		case "":
		case "curve25519":
		default:
			return multierror.Prefix(errors.New("invalid value for 'dh_type'"), fmt.Sprintf("sink.%s", s.Type))
		}

		if s.AADEnvVar != "" {
			s.AAD = os.Getenv(s.AADEnvVar)
			s.AADEnvVar = ""
		}

		switch {
		case s.DHPath == "" && s.DHType == "":
			if s.AAD != "" {
				return multierror.Prefix(errors.New("specifying AAD data without 'dh_type' does not make sense"), fmt.Sprintf("sink.%s", s.Type))
			}
		case s.DHPath != "" && s.DHType != "":
		default:
			return multierror.Prefix(errors.New("'dh_type' and 'dh_path' must be specified together"), fmt.Sprintf("sink.%s", s.Type))
		}

		ts = append(ts, &s)
	}

	result.AutoAuth.Sinks = ts
	return nil
}
