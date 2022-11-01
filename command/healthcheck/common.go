/*
 * The healthcheck package attempts to allow generic checks of arbitrary
 * engines, while providing a common framework with some performance
 * efficiencies in mind.
 *
 * The core of this package is the Executor context; a caller would
 * provision a set of checks, an API client, and a configuration,
 * which the executor would use to decide which checks to execute
 * and how.
 *
 * Checks are based around a series of remote paths that are fetched by
 * the client; these are broken into two categories: static paths, which
 * can always be fetched; and dynamic paths, which the check fetches based
 * on earlier results.
 *
 * For instance, a basic PKI CA lifetime check will have static fetch against
 * the list of CAs, and a dynamic fetch, using that earlier list, to fetch the
 * PEMs of all CAs.
 *
 * This allows health checks to share data: many PKI checks will need the
 * issuer list and so repeatedly fetching this may result in a performance
 * impact.
 */

package healthcheck

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"
)

type Executor struct {
	Client         *api.Client
	Mount          string
	DefaultEnabled bool

	Config map[string]map[string]interface{}

	Resources map[string]map[logical.Operation]*PathFetch

	Checkers []Check
}

func NewExecutor(client *api.Client, mount string) *Executor {
	return &Executor{
		Client:         client,
		DefaultEnabled: true,
		Mount:          mount,
		Config:         make(map[string]map[string]interface{}),
		Resources:      make(map[string]map[logical.Operation]*PathFetch),
	}
}

func (e *Executor) AddCheck(c Check) {
	e.Checkers = append(e.Checkers, c)
}

func (e *Executor) BuildConfig(external map[string]interface{}) error {
	merged := e.Config

	for index, checker := range e.Checkers {
		name := checker.Name()
		if _, present := merged[name]; name == "" || present {
			return fmt.Errorf("bad checker %v: name is empty or already present: %v", index, name)
		}

		// Fetch the default configuration; if the check returns enabled
		// status, verify it matches our expectations (in the event it should
		// be disabled by default), otherwise, add it in.
		config := checker.DefaultConfig()
		enabled, present := config["enabled"]
		if !present {
			config["enabled"] = e.DefaultEnabled
		} else if enabled.(bool) && !e.DefaultEnabled {
			config["enabled"] = e.DefaultEnabled
		}

		// Now apply any external config for this check.
		if econfig, present := external[name]; present {
			for param, evalue := range econfig.(map[string]interface{}) {
				if _, ok := config[param]; !ok {
					// Assumption: default configs have all possible
					// configuration options. This external config has
					// an unknown option, so we want to error out.
					return fmt.Errorf("unknown configuration option for %v: %v", name, param)
				}

				config[param] = evalue
			}
		}

		// Now apply it and save it.
		if err := checker.LoadConfig(config); err != nil {
			return fmt.Errorf("error saving merged config for %v: %w", name, err)
		}
		merged[name] = config
	}

	return nil
}

func (e *Executor) Execute() (map[string][]*Result, error) {
	ret := make(map[string][]*Result)
	for _, checker := range e.Checkers {
		if !checker.IsEnabled() {
			continue
		}

		if err := checker.FetchResources(e); err != nil {
			return nil, err
		}

		results, err := checker.Evaluate(e)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			result.Endpoint = e.templatePath(result.Endpoint)
			result.StatusDisplay = ResultStatusNameMap[result.Status]
		}

		ret[checker.Name()] = results
	}

	return ret, nil
}

func (e *Executor) templatePath(path string) string {
	return strings.ReplaceAll(path, "{{mount}}", e.Mount)
}

func (e *Executor) FetchIfNotFetched(op logical.Operation, rawPath string) (*PathFetch, error) {
	path := e.templatePath(rawPath)

	byOp, present := e.Resources[path]
	if present && byOp != nil {
		result, present := byOp[op]
		if present && result != nil {
			return result, result.FetchSurfaceError()
		}
	}

	// Must not exist in cache; create it.
	if byOp == nil {
		e.Resources[path] = make(map[logical.Operation]*PathFetch)
	}

	ret := &PathFetch{
		Operation:   op,
		Path:        path,
		ParsedCache: make(map[string]interface{}),
	}

	data := map[string][]string{}
	if op == logical.ListOperation {
		data["list"] = []string{"true"}
	} else if op != logical.ReadOperation {
		return nil, fmt.Errorf("unknown operation: %v on %v", op, path)
	}

	response, err := e.Client.Logical().ReadRawWithData(path, data)
	ret.Response = response
	if err != nil {
		ret.FetchError = err
	} else {
		// Not all secrets will parse correctly. Sometimes we really want
		// to fetch a raw endpoint, sometimes we're run with a bad mount
		// or missing permissions.
		secret, secretErr := e.Client.Logical().ParseRawResponseAndCloseBody(response, err)
		if secretErr != nil {
			ret.SecretParseError = secretErr
		} else {
			ret.Secret = secret
		}
	}

	e.Resources[path][op] = ret
	return ret, ret.FetchSurfaceError()
}

type PathFetch struct {
	Operation        logical.Operation
	Path             string
	Response         *api.Response
	FetchError       error
	Secret           *api.Secret
	SecretParseError error
	ParsedCache      map[string]interface{}
}

func (p *PathFetch) IsOK() bool {
	return p.FetchError == nil && p.Response != nil
}

func (p *PathFetch) IsSecretOK() bool {
	return p.IsOK() && p.SecretParseError == nil && p.Secret != nil
}

func (p *PathFetch) FetchSurfaceError() error {
	if p.IsOK() || p.IsSecretPermissionsError() || p.IsUnsupportedPathError() || p.IsMissingResource() {
		return nil
	}

	if strings.Contains(p.FetchError.Error(), "route entry not found") {
		return fmt.Errorf("Error making API request: was a bad mount given?\n\nOperation: %v\nPath: %v\nOriginal Error:\n%v", p.Operation, p.Path, p.FetchError)
	}

	return p.FetchError
}

func (p *PathFetch) IsSecretPermissionsError() bool {
	return !p.IsOK() && strings.Contains(p.FetchError.Error(), "permission denied")
}

func (p *PathFetch) IsUnsupportedPathError() bool {
	return !p.IsOK() && strings.Contains(p.FetchError.Error(), "unsupported path")
}

func (p *PathFetch) IsMissingResource() bool {
	return !p.IsOK() && strings.Contains(p.FetchError.Error(), "unable to find")
}

type Check interface {
	Name() string
	IsEnabled() bool

	DefaultConfig() map[string]interface{}
	LoadConfig(config map[string]interface{}) error

	FetchResources(e *Executor) error

	Evaluate(e *Executor) ([]*Result, error)
}

type ResultStatus int

const (
	ResultNotApplicable ResultStatus = iota
	ResultOK
	ResultInformational
	ResultWarning
	ResultCritical
	ResultInvalidVersion
	ResultInsufficientPermissions
)

var ResultStatusNameMap = map[ResultStatus]string{
	ResultNotApplicable:           "not_applicable",
	ResultOK:                      "ok",
	ResultInformational:           "informational",
	ResultWarning:                 "warning",
	ResultCritical:                "critical",
	ResultInvalidVersion:          "invalid_version",
	ResultInsufficientPermissions: "insufficient_permissions",
}

var NameResultStatusMap = map[string]ResultStatus{
	"not_applicable":           ResultNotApplicable,
	"ok":                       ResultOK,
	"informational":            ResultInformational,
	"warning":                  ResultWarning,
	"critical":                 ResultCritical,
	"invalid_version":          ResultInvalidVersion,
	"insufficient_permissions": ResultInsufficientPermissions,
}

type Result struct {
	Status        ResultStatus `json:"status_code"`
	StatusDisplay string       `json:"status"`
	Endpoint      string       `json:"endpoint,omitempty"`
	Message       string       `json:"message,omitempty"`
}
