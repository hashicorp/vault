package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"text/template"

	"github.com/pkg/errors"

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
)

// Template is the internal representation of an individual template to process.
// The template retains the relationship between it's contents and is
// responsible for it's own execution.
type Template struct {
	// contents is the string contents for the template. It is either given
	// during template creation or read from disk when initialized.
	contents string

	// source is the original location of the template. This may be undefined if
	// the template was dynamically defined.
	source string

	// leftDelim and rightDelim are the template delimiters.
	leftDelim  string
	rightDelim string

	// hexMD5 stores the hex version of the MD5
	hexMD5 string

	// errMissingKey causes the template processing to exit immediately if a map
	// is indexed with a key that does not exist.
	errMissingKey bool

	// functionDenylist are functions not permitted to be executed
	// when we render this template
	functionDenylist []string

	// sandboxPath adds a prefix to any path provided to the `file` function
	// and causes an error if a relative path tries to traverse outside that
	// prefix.
	sandboxPath string
}

// NewTemplateInput is used as input when creating the template.
type NewTemplateInput struct {
	// Source is the location on disk to the file.
	Source string

	// Contents are the raw template contents.
	Contents string

	// ErrMissingKey causes the template parser to exit immediately with an error
	// when a map is indexed with a key that does not exist.
	ErrMissingKey bool

	// LeftDelim and RightDelim are the template delimiters.
	LeftDelim  string
	RightDelim string

	// FunctionDenylist are functions not permitted to be executed
	// when we render this template
	FunctionDenylist []string

	// SandboxPath adds a prefix to any path provided to the `file` function
	// and causes an error if a relative path tries to traverse outside that
	// prefix.
	SandboxPath string
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
	t.functionDenylist = i.FunctionDenylist
	t.sandboxPath = i.SandboxPath

	if i.Source != "" {
		contents, err := ioutil.ReadFile(i.Source)
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

// Source returns the filepath source of this template.
func (t *Template) Source() string {
	if t.source == "" {
		return "(dynamic)"
	}
	return t.source
}

// ExecuteInput is used as input to the template's execute function.
type ExecuteInput struct {
	// Brain is the brain where data for the template is stored.
	Brain *Brain

	// Env is a custom environment provided to the template for envvar resolution.
	// Values specified here will take precedence over any values in the
	// environment when using the `env` function.
	Env []string
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
		t:                tmpl,
		brain:            i.Brain,
		env:              i.Env,
		used:             &used,
		missing:          &missing,
		functionDenylist: t.functionDenylist,
		sandboxPath:      t.sandboxPath,
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
		return nil, errors.Wrap(err, "execute")
	}

	return &ExecuteResult{
		Used:    &used,
		Missing: &missing,
		Output:  b.Bytes(),
	}, nil
}

// funcMapInput is input to the funcMap, which builds the template functions.
type funcMapInput struct {
	t                *template.Template
	brain            *Brain
	env              []string
	functionDenylist []string
	sandboxPath      string
	used             *dep.Set
	missing          *dep.Set
}

// funcMap is the map of template functions to their respective functions.
func funcMap(i *funcMapInput) template.FuncMap {
	var scratch Scratch

	r := template.FuncMap{
		// API functions
		"datacenters":  datacentersFunc(i.brain, i.used, i.missing),
		"file":         fileFunc(i.brain, i.used, i.missing, i.sandboxPath),
		"key":          keyFunc(i.brain, i.used, i.missing),
		"keyExists":    keyExistsFunc(i.brain, i.used, i.missing),
		"keyOrDefault": keyWithDefaultFunc(i.brain, i.used, i.missing),
		"ls":           lsFunc(i.brain, i.used, i.missing, true),
		"safeLs":       safeLsFunc(i.brain, i.used, i.missing),
		"node":         nodeFunc(i.brain, i.used, i.missing),
		"nodes":        nodesFunc(i.brain, i.used, i.missing),
		"secret":       secretFunc(i.brain, i.used, i.missing),
		"secrets":      secretsFunc(i.brain, i.used, i.missing),
		"service":      serviceFunc(i.brain, i.used, i.missing),
		"connect":      connectFunc(i.brain, i.used, i.missing),
		"services":     servicesFunc(i.brain, i.used, i.missing),
		"tree":         treeFunc(i.brain, i.used, i.missing, true),
		"safeTree":     safeTreeFunc(i.brain, i.used, i.missing),
		"caRoots":      connectCARootsFunc(i.brain, i.used, i.missing),
		"caLeaf":       connectLeafFunc(i.brain, i.used, i.missing),

		// Scratch
		"scratch": func() *Scratch { return &scratch },

		// Helper functions
		"base64Decode":    base64Decode,
		"base64Encode":    base64Encode,
		"base64URLDecode": base64URLDecode,
		"base64URLEncode": base64URLEncode,
		"byKey":           byKey,
		"byTag":           byTag,
		"contains":        contains,
		"containsAll":     containsSomeFunc(true, true),
		"containsAny":     containsSomeFunc(false, false),
		"containsNone":    containsSomeFunc(true, false),
		"containsNotAll":  containsSomeFunc(false, true),
		"env":             envFunc(i.env),
		"executeTemplate": executeTemplateFunc(i.t),
		"explode":         explode,
		"explodeMap":      explodeMap,
		"in":              in,
		"indent":          indent,
		"loop":            loop,
		"join":            join,
		"trimSpace":       trimSpace,
		"parseBool":       parseBool,
		"parseFloat":      parseFloat,
		"parseInt":        parseInt,
		"parseJSON":       parseJSON,
		"parseUint":       parseUint,
		"parseYAML":       parseYAML,
		"plugin":          plugin,
		"regexReplaceAll": regexReplaceAll,
		"regexMatch":      regexMatch,
		"replaceAll":      replaceAll,
		"sha256Hex":       sha256Hex,
		"timestamp":       timestamp,
		"toLower":         toLower,
		"toJSON":          toJSON,
		"toJSONPretty":    toJSONPretty,
		"toTitle":         toTitle,
		"toTOML":          toTOML,
		"toUpper":         toUpper,
		"toYAML":          toYAML,
		"split":           split,
		"byMeta":          byMeta,
		"sockaddr":        sockaddr,
		// Math functions
		"add":      add,
		"subtract": subtract,
		"multiply": multiply,
		"divide":   divide,
		"modulo":   modulo,
		"minimum":  minimum,
		"maximum":  maximum,
	}

	for _, bf := range i.functionDenylist {
		if _, ok := r[bf]; ok {
			r[bf] = denied
		}
	}

	return r
}
