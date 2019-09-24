package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/y0ssar1an/q"
)

// Config is the configuration for the vault server.
type Config struct {
	AutoAuth      *AutoAuth         `hcl:"auto_auth"`
	ExitAfterAuth bool              `hcl:"exit_after_auth"`
	PidFile       string            `hcl:"pid_file"`
	Listeners     []*Listener       `hcl:"listeners"`
	Cache         *Cache            `hcl:"cache"`
	Vault         *Vault            `hcl:"vault"`
	Templates     []*TemplateConfig `hcl:"templates"`
}

type Vault struct {
	Address          string      `hcl:"address"`
	CACert           string      `hcl:"ca_cert"`
	CAPath           string      `hcl:"ca_path"`
	TLSSkipVerify    bool        `hcl:"-"`
	TLSSkipVerifyRaw interface{} `hcl:"tls_skip_verify"`
	ClientCert       string      `hcl:"client_cert"`
	ClientKey        string      `hcl:"client_key"`
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
	Namespace  string        `hcl:"namespace"`
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

type TemplateConfig struct {
	// Backup determines if this template should retain a backup. The default
	// value is false.
	Backup *bool `hcl:"backup"`

	// Command is the arbitrary command to execute after a template has
	// successfully rendered. This is DEPRECATED. Use Exec instead.
	Command *string `hcl:"command"`

	// CommandTimeout is the amount of time to wait for the command to finish
	// before force-killing it. This is DEPRECATED. Use Exec instead.
	CommandTimeoutRaw interface{}    `hcl:"command_timeout"`
	CommandTimeout    *time.Duration `hcl:"-"`

	// Contents are the raw template contents to evaluate. Either this or Source
	// must be specified, but not both.
	Contents *string `hcl:"contents"`

	// CreateDestDirs tells Consul Template to create the parent directories of
	// the destination path if they do not exist. The default value is true.
	CreateDestDirs *bool `hcl:"create_dest_dirs"`

	// Destination is the location on disk where the template should be rendered.
	// This is required unless running in debug/dry mode.
	Destination *string `hcl:"destination"`

	// ErrMissingKey is used to control how the template behaves when attempting
	// to index a struct or map key that does not exist.
	ErrMissingKey *bool `hcl:"error_on_missing_key"`

	// // Exec is the configuration for the command to run when the template renders
	// // successfully.
	// Exec *ExecConfig `hcl:"exec"`

	// Perms are the file system permissions to use when creating the file on
	// disk. This is useful for when files contain sensitive information, such as
	// secrets from Vault.
	PermsRaw interface{}  `hcl:"perms"`
	Perms    *os.FileMode `hcl:"-"`

	// Source is the path on disk to the template contents to evaluate. Either
	// this or Contents should be specified, but not both.
	Source *string `hcl:"source"`

	// // Wait configures per-template quiescence timers.
	// Wait *WaitConfig `hcl:"wait"`

	// LeftDelim and RightDelim are optional configurations to control what
	// delimiter is utilized when parsing the template.
	LeftDelim  *string `hcl:"left_delimiter"`
	RightDelim *string `hcl:"right_delimiter"`

	// FunctionBlacklist is a list of functions that this template is not
	// permitted to run.
	FunctionBlacklist []string `hcl:"function_blacklist"`

	// SandboxPath adds a prefix to any path provided to the `file` function
	// and causes an error if a relative path tries to traverse outside that
	// prefix.
	SandboxPath *string `hcl:"sandbox_path"`
}

// LoadConfig loads the configuration at the given path, regardless if
// its a file or directory.
func LoadConfig(path string) (*Config, error) {
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

	err = parseTemplates(&result, list)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing 'templates':{{err}}", err)
	}

	if result.Cache != nil {
		if len(result.Listeners) < 1 {
			return nil, fmt.Errorf("at least one listener required when cache enabled")
		}

		if result.Cache.UseAutoAuthToken {
			if result.AutoAuth == nil {
				return nil, fmt.Errorf("cache.use_auto_auth_token is true but auto_auth not configured")
			}
			if result.AutoAuth.Method.WrapTTL > 0 {
				return nil, fmt.Errorf("cache.use_auto_auth_token is true and auto_auth uses wrapping")
			}
		}
	}

	if result.AutoAuth != nil {
		if len(result.AutoAuth.Sinks) == 0 && (result.Cache == nil || !result.Cache.UseAutoAuthToken) {
			return nil, fmt.Errorf("auto_auth requires at least one sink or cache.use_auto_auth_token=true ")
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

	if v.TLSSkipVerifyRaw != nil {
		v.TLSSkipVerify, err = parseutil.ParseBool(v.TLSSkipVerifyRaw)
		if err != nil {
			return err
		}
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

	if result.AutoAuth.Method.WrapTTL > 0 {
		if len(result.AutoAuth.Sinks) != 1 {
			return fmt.Errorf("error parsing auto_auth: wrapping enabled on auth method and 0 or many sinks defined")
		}

		if result.AutoAuth.Sinks[0].WrapTTL > 0 {
			return fmt.Errorf("error parsing auto_auth: wrapping enabled both on auth method and sink")
		}
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

	// Canonicalize namespace path if provided
	m.Namespace = namespace.Canonicalize(m.Namespace)

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

func parseTemplates(result *Config, list *ast.ObjectList) error {
	name := "template"

	templateList := list.Filter(name)
	if len(templateList.Items) < 1 {
		return nil
	}

	var tcs []*TemplateConfig

	for _, item := range templateList.Items {
		var tc TemplateConfig
		if err := hcl.DecodeObject(&tc, item.Val); err != nil {
			q.Q("error here:", err)
			return err
		}

		if tc.CommandTimeoutRaw != nil {
			timeout, err := parseutil.ParseDurationSecond(tc.CommandTimeoutRaw)
			if err != nil {
				return err
			}
			tc.CommandTimeout = &timeout
			tc.CommandTimeoutRaw = nil
		}

		q.Q(tc.PermsRaw)
		perms := os.FileMode(0644)
		if tc.PermsRaw != nil {
			switch tc.PermsRaw.(type) {
			case int:
				perms = os.FileMode(tc.PermsRaw.(int))
			default:
				return errors.New("error parsing perms")
			}
			tc.PermsRaw = nil
			tc.Perms = &perms
		}
		q.Q(tc.Perms.String())

		// check source vs contents
		// check command / timeout

		tcs = append(tcs, &tc)
	}

	result.Templates = tcs
	return nil
}
