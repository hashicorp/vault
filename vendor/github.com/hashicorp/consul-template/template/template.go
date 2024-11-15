// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/hashicorp/consul-template/config"
	dep "github.com/hashicorp/consul-template/dependency"
)

var (
	// ErrTemplateContentsAndSource is the error returned when a template
	// specifies both a "source" and "content" argument, which is not valid.
	ErrTemplateContentsAndSource = errors.New("template: cannot specify both 'source' and 'contents'")

	// ErrTemplateMissingContentsAndSource is the error returned when a template
	// does not specify either a "source" or "content" argument, which is not
	// valid.
	ErrTemplateMissingContentsAndSource = errors.New("template: must specify exactly one of 'source' or 'contents'")

	// ErrMissingReaderFunction is the error returned when the template
	// configuration is missing a reader function.
	ErrMissingReaderFunction = errors.New("template: missing a reader function")
)

// Template is the internal representation of an individual template to process.
// The template retains the relationship between its contents and is
// responsible for it's own execution.
type Template struct {
	// contents is the string contents for the template. It is either given
	// during template creation or read from disk when initialized.
	contents string

	// source is the original location of the template. This may be undefined if
	// the template was dynamically defined.
	source string

	// destination file/path to which the template is rendered
	destination string

	// leftDelim and rightDelim are the template delimiters.
	leftDelim  string
	rightDelim string

	// hexMD5 stores the hex version of the MD5
	hexMD5 string

	// errMissingKey causes the template processing to exit immediately if a map
	// is indexed with a key that does not exist.
	errMissingKey bool

	// errFatal determines whether template errors should cause the process to
	// exit, or just log and continue.
	errFatal bool

	// FuncMap is a map of external functions that this template is
	// permitted to run. Allows users to add functions to the library
	// and selectively opaque existing ones.
	extFuncMap template.FuncMap

	// functionDenylist are functions not permitted to be executed
	// when we render this template
	functionDenylist []string

	// sandboxPath adds a prefix to any path provided to the `file` function
	// and causes an error if a relative path tries to traverse outside that
	// prefix.
	sandboxPath string

	// local reference to configuration for this template
	config *config.TemplateConfig
}

// NewTemplateInput is used as input when creating the template.
type NewTemplateInput struct {
	// Source is the location on disk to the file.
	Source string

	// Destination is the file on disk to render/write the template output.
	Destination string

	// Contents are the raw template contents.
	Contents string

	// ErrMissingKey causes the template parser to exit immediately with an error
	// when a map is indexed with a key that does not exist.
	ErrMissingKey bool

	// ErrFatal determines whether template errors should cause the process to
	// exit, or just log and continue.
	ErrFatal bool

	// LeftDelim and RightDelim are the template delimiters.
	LeftDelim  string
	RightDelim string

	// ExtFuncMap is a map of external functions that this template is
	// permitted to run. Allows users to add functions to the library
	// and selectively opaque existing ones.
	ExtFuncMap template.FuncMap

	// FunctionDenylist are functions not permitted to be executed
	// when we render this template
	FunctionDenylist []string

	// SandboxPath adds a prefix to any path provided to the `file` function
	// and causes an error if a relative path tries to traverse outside that
	// prefix.
	SandboxPath string

	// Config keeps local reference to config struct
	Config *config.TemplateConfig

	// ReaderFunc is called to read in any source file
	ReaderFunc config.Reader
}

// NewTemplate creates and parses a new Consul Template template at the given
// path. If the template does not exist, an error is returned. During
// initialization, the template is read and is parsed for dependencies. Any
// errors that occur are returned.
func NewTemplate(i *NewTemplateInput) (*Template, error) {
	if i == nil {
		i = &NewTemplateInput{}
	}

	// Validate that we are either given the path or the explicit contents
	if i.Source != "" && i.Contents != "" {
		return nil, ErrTemplateContentsAndSource
	} else if i.Source == "" && i.Contents == "" {
		return nil, ErrTemplateMissingContentsAndSource
	}

	var t Template
	t.source = i.Source
	t.contents = i.Contents
	t.leftDelim = i.LeftDelim
	t.rightDelim = i.RightDelim
	t.errMissingKey = i.ErrMissingKey
	t.errFatal = i.ErrFatal
	t.extFuncMap = i.ExtFuncMap
	t.functionDenylist = i.FunctionDenylist
	t.sandboxPath = i.SandboxPath
	t.destination = i.Destination
	t.config = i.Config

	if i.ExtFuncMap != nil {
		t.extFuncMap = make(map[string]any, len(i.ExtFuncMap))
		maps.Copy(t.extFuncMap, i.ExtFuncMap)
	}

	if i.Source != "" {
		if i.ReaderFunc == nil {
			return nil, ErrMissingReaderFunction
		}
		contents, err := i.ReaderFunc(i.Source)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read template")
		}
		t.contents = string(contents)
	}

	// Compute the MD5, encode as hex
	hash := md5.Sum([]byte(t.contents))
	t.hexMD5 = hex.EncodeToString(hash[:])

	return &t, nil
}

// ID returns the identifier for this template.
func (t *Template) ID() string {
	return t.hexMD5
}

// Contents returns the raw contents of the template.
func (t *Template) Contents() string {
	return t.contents
}

// Config returns the template's config
func (t *Template) Config() *config.TemplateConfig {
	return t.config
}

// Source returns the filepath source of this template.
func (t *Template) Source() string {
	if t.source == "" {
		return "(dynamic)"
	}
	return t.source
}

// ErrFatal indicates whether errors in this template should be fatal.
func (t *Template) ErrFatal() bool {
	return t.errFatal
}

// ExecuteInput is used as input to the template's execute function.
type ExecuteInput struct {
	// Brain is the brain where data for the template is stored.
	Brain *Brain

	// Env is a custom environment provided to the template for envvar resolution.
	// Values specified here will take precedence over any values in the
	// environment when using the `env` function.
	Env []string

	// Config is a copy of the Runner's consul-template configuration. It is
	// provided to allow for functions that might need to adapt based on certain
	// configuration values
	Config *config.Config
}

// ExecuteResult is the result of the template execution.
type ExecuteResult struct {
	// Used is the set of dependencies that were used.
	Used *dep.Set

	// Missing is the set of dependencies that were missing.
	Missing *dep.Set

	// Output is the rendered result.
	Output []byte
}

// Execute evaluates this template in the provided context.
func (t *Template) Execute(i *ExecuteInput) (*ExecuteResult, error) {
	if i == nil {
		i = &ExecuteInput{}
	}

	var used, missing dep.Set

	tmpl := template.New("")
	tmpl.Delims(t.leftDelim, t.rightDelim)

	tmpl.Funcs(funcMap(&funcMapInput{
		newTmpl:          tmpl,
		brain:            i.Brain,
		env:              i.Env,
		used:             &used,
		missing:          &missing,
		extFuncMap:       t.extFuncMap,
		functionDenylist: t.functionDenylist,
		sandboxPath:      t.sandboxPath,
		destination:      t.destination,
		config:           i.Config,
	}))

	if t.errMissingKey {
		tmpl.Option("missingkey=error")
	} else {
		tmpl.Option("missingkey=zero")
	}

	tmpl, err := tmpl.Parse(t.contents)
	if err != nil {
		return nil, errors.Wrap(err, "parse")
	}

	// Execute the template into the writer
	var b bytes.Buffer
	if err := tmpl.Execute(&b, nil); err != nil {
		return nil, errors.Wrap(redactinator(&used, i.Brain, err), "execute")
	}

	return &ExecuteResult{
		Used:    &used,
		Missing: &missing,
		Output:  b.Bytes(),
	}, nil
}

func redactinator(used *dep.Set, b *Brain, err error) error {
	pairs := make([]string, 0, used.Len())
	for _, d := range used.List() {
		if data, ok := b.Recall(d); ok {
			if vd, ok := data.(*dep.Secret); ok {
				for _, v := range vd.Data {
					pairs = append(pairs, fmt.Sprintf("%v", v), "[redacted]")
				}
			}
			if nVar, ok := data.(*dep.NomadVarItems); ok {
				for _, v := range nVar.Values() {
					pairs = append(pairs, fmt.Sprintf("%v", v), "[redacted]")
				}
			}
		}
	}
	return fmt.Errorf(strings.NewReplacer(pairs...).Replace(err.Error()))
}

// funcMapInput is input to the funcMap, which builds the template functions.
type funcMapInput struct {
	newTmpl          *template.Template
	brain            *Brain
	env              []string
	extFuncMap       map[string]interface{}
	functionDenylist []string
	sandboxPath      string
	destination      string
	used             *dep.Set
	missing          *dep.Set
	config           *config.Config
}

// funcMap is the map of template functions to their respective functions.
func funcMap(i *funcMapInput) template.FuncMap {
	var scratch Scratch

	// Get the Nomad default namespace from the client config
	// this is done here rather than in the function to prevent an
	// import cycle between the dependency and config packages
	nomadNS := "default"
	if i.config != nil && i.config.Nomad != nil && *(i.config.Nomad).Namespace != "" {
		nomadNS = *(i.config.Nomad).Namespace
	}

	r := template.FuncMap{
		// API functions
		"datacenters":      datacentersFunc(i.brain, i.used, i.missing),
		"exportedServices": exportedServicesFunc(i.brain, i.used, i.missing),
		"file":             fileFunc(i.brain, i.used, i.missing, i.sandboxPath),
		"key":              keyFunc(i.brain, i.used, i.missing),
		"keyExists":        keyExistsFunc(i.brain, i.used, i.missing),
		"keyOrDefault":     keyWithDefaultFunc(i.brain, i.used, i.missing),
		"ls":               lsFunc(i.brain, i.used, i.missing, true),
		"safeLs":           safeLsFunc(i.brain, i.used, i.missing),
		"node":             nodeFunc(i.brain, i.used, i.missing),
		"nodes":            nodesFunc(i.brain, i.used, i.missing),
		"partitions":       partitionsFunc(i.brain, i.used, i.missing),
		"peerings":         peeringsFunc(i.brain, i.used, i.missing),
		"secret":           secretFunc(i.brain, i.used, i.missing),
		"secrets":          secretsFunc(i.brain, i.used, i.missing),
		"service":          serviceFunc(i.brain, i.used, i.missing),
		"connect":          connectFunc(i.brain, i.used, i.missing),
		"services":         servicesFunc(i.brain, i.used, i.missing),
		"tree":             treeFunc(i.brain, i.used, i.missing, true),
		"safeTree":         safeTreeFunc(i.brain, i.used, i.missing),
		"caRoots":          connectCARootsFunc(i.brain, i.used, i.missing),
		"caLeaf":           connectLeafFunc(i.brain, i.used, i.missing),
		"pkiCert":          pkiCertFunc(i.brain, i.used, i.missing, i.destination),

		// Nomad Functions.
		"nomadServices":    nomadServicesFunc(i.brain, i.used, i.missing),
		"nomadService":     nomadServiceFunc(i.brain, i.used, i.missing),
		"nomadVarList":     nomadVariablesFunc(i.brain, i.used, i.missing, nomadNS, true),
		"nomadVarListSafe": nomadSafeVariablesFunc(i.brain, i.used, i.missing, nomadNS),
		"nomadVar":         nomadVariableItemsFunc(i.brain, i.used, i.missing, nomadNS),
		"nomadVarExists":   nomadVariableExistsFunc(i.brain, i.used, i.missing, nomadNS),

		// Scratch
		"scratch": func() *Scratch { return &scratch },

		// Helper functions
		"base64Decode":          base64Decode,
		"base64Encode":          base64Encode,
		"base64URLDecode":       base64URLDecode,
		"base64URLEncode":       base64URLEncode,
		"byKey":                 byKey,
		"byPort":                byPort,
		"byTag":                 byTag,
		"contains":              contains,
		"containsAll":           containsSomeFunc(true, true),
		"containsAny":           containsSomeFunc(false, false),
		"containsNone":          containsSomeFunc(true, false),
		"containsNotAll":        containsSomeFunc(false, true),
		"env":                   envFunc(i.env),
		"mustEnv":               mustEnvFunc(i.env),
		"envOrDefault":          envWithDefaultFunc(i.env),
		"executeTemplate":       executeTemplateFunc(i.newTmpl),
		"explode":               explode,
		"explodeMap":            explodeMap,
		"mergeMap":              mergeMap,
		"mergeMapWithOverride":  mergeMapWithOverride,
		"in":                    in,
		"indent":                indent,
		"loop":                  loop,
		"join":                  join,
		"trim":                  trim,
		"trimPrefix":            trimPrefix,
		"trimSuffix":            trimSuffix,
		"trimSpace":             trimSpace,
		"parseBool":             parseBool,
		"parseFloat":            parseFloat,
		"parseInt":              parseInt,
		"parseJSON":             parseJSON,
		"parseUint":             parseUint,
		"parseYAML":             parseYAML,
		"plugin":                plugin,
		"regexReplaceAll":       regexReplaceAll,
		"regexMatch":            regexMatch,
		"replaceAll":            replaceAll,
		"sha256Hex":             sha256Hex,
		"md5sum":                md5sum,
		"hmacSHA256Hex":         hmacSHA256Hex,
		"timestamp":             timestamp,
		"toLower":               toLower,
		"toJSON":                toJSON,
		"toJSONPretty":          toJSONPretty,
		"toUnescapedJSON":       toUnescapedJSON,
		"toUnescapedJSONPretty": toUnescapedJSONPretty,
		"toTitle":               toTitle,
		"toTOML":                toTOML,
		"toUpper":               toUpper,
		"toYAML":                toYAML,
		"split":                 split,
		"splitToMap":            splitToMap,
		"byMeta":                byMeta,
		"sockaddr":              sockaddr,
		"writeToFile":           writeToFile,

		// Math functions
		"add":      add,
		"subtract": subtract,
		"multiply": multiply,
		"divide":   divide,
		"modulo":   modulo,
		"minimum":  minimum,
		"maximum":  maximum,
		// Debug functions
		"spew_dump":    spewDump,
		"spew_printf":  spewPrintf,
		"spew_sdump":   spewSdump,
		"spew_sprintf": spewSprintf,
	}

	// Add the Sprig functions to the funcmap
	for k, v := range sprig.TxtFuncMap() {
		target := "sprig_" + k
		r[target] = v
	}

	// Add external functions
	if i.extFuncMap != nil {
		for name, fn := range i.extFuncMap {
			r[name] = fn
		}
	}

	for _, bf := range i.functionDenylist {
		if _, ok := r[bf]; ok {
			r[bf] = denied
		}
	}

	return r
}
