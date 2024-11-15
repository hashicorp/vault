// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"golang.org/x/exp/maps"
)

const (
	// DefaultTemplateCommandTimeout is the amount of time to wait for a command
	// to return.
	DefaultTemplateCommandTimeout = 30 * time.Second
)

var (
	// ErrTemplateStringEmpty is the error returned with the template contents
	// are empty.
	ErrTemplateStringEmpty = errors.New("template: cannot be empty")

	// configTemplateRe is the pattern to split the config template syntax.
	configTemplateRe = regexp.MustCompile("([a-zA-Z]:)?([^:]+)")
)

// TemplateConfig is a representation of a template on disk, as well as the
// associated commands and reload instructions.
type TemplateConfig struct {
	// Backup determines if this template should retain a backup. The default
	// value is false.
	Backup *bool `mapstructure:"backup"`

	// Command is the arbitrary command to execute after a template has
	// successfully rendered. This is DEPRECATED. Use Exec instead.
	Command commandList `mapstructure:"command"`

	// CommandTimeout is the amount of time to wait for the command to finish
	// before force-killing it. This is DEPRECATED. Use Exec instead.
	CommandTimeout *time.Duration `mapstructure:"command_timeout"`

	// Contents are the raw template contents to evaluate. Either this or Source
	// must be specified, but not both.
	Contents *string `mapstructure:"contents"`

	// CreateDestDirs tells Consul Template to create the parent directories of
	// the destination path if they do not exist. The default value is true.
	CreateDestDirs *bool `mapstructure:"create_dest_dirs"`

	// Destination is the location on disk where the template should be rendered.
	// This is required unless running in debug/dry mode.
	Destination *string `mapstructure:"destination"`

	// ErrMissingKey is used to control how the template behaves when attempting
	// to index a struct or map key that does not exist.
	ErrMissingKey *bool `mapstructure:"error_on_missing_key"`

	// ErrFatal determines whether template errors should cause the process to
	// exit, or just log and continue.
	ErrFatal *bool `mapstructure:"error_fatal"`

	// Exec is the configuration for the command to run when the template renders
	// successfully.
	Exec *ExecConfig `mapstructure:"exec"`

	// Perms are the file system permissions to use when creating the file on
	// disk. This is useful for when files contain sensitive information, such as
	// secrets from Vault.
	Perms *os.FileMode `mapstructure:"perms"`

	// User is the username or uid that will be set when creating the file on disk.
	// Useful when simply setting Perms is not enough.
	//
	// Platform dependent: this doesn't work on Windows but it fails gracefully
	// with a warning
	User *string `mapstructure:"user"`

	// Uid is equivalent to User when it's a valid int and it's defined for backward compatibility with v0.28.0
	Uid *int `mapstructure:"uid"`

	// Group is the group name or gid that will be set when creating the file on disk.
	// Useful when simply setting Perms is not enough.
	//
	// Platform dependent: this doesn't work on Windows but it fails gracefully
	// with a warning
	Group *string `mapstructure:"group"`

	// Gid is equivalent to Group when it's a valid int and it's defined for backward compatibility with v0.28.0
	Gid *int `mapstructure:"gid"`

	// Source is the path on disk to the template contents to evaluate. Either
	// this or Contents should be specified, but not both.
	Source *string `mapstructure:"source"`

	// Wait configures per-template quiescence timers.
	Wait *WaitConfig `mapstructure:"wait"`

	// LeftDelim and RightDelim are optional configurations to control what
	// delimiter is utilized when parsing the template.
	LeftDelim  *string `mapstructure:"left_delimiter"`
	RightDelim *string `mapstructure:"right_delimiter"`

	// ExtFuncMap is a map of external functions that this template is
	// permitted to run. Allows users to add functions to the library
	// and selectively opaque existing ones. Omitted from json output
	// to prevent errors when the configuration is marshalled for printing.
	ExtFuncMap template.FuncMap `json:"-"`

	// FunctionDenylist is a list of functions that this template is not
	// permitted to run.
	FunctionDenylist []string `mapstructure:"function_denylist"`

	// FunctionDenylistDeprecated is the backward compatible option for
	// FunctionDenylist for configuration supported by v0.25.0 and older. This
	// should not be used directly, use FunctionDenylist instead. Values from
	// this are combined to FunctionDenylist in Finalize().
	FunctionDenylistDeprecated []string `mapstructure:"function_blacklist" json:"-"`

	// SandboxPath adds a prefix to any path provided to the `file` function
	// and causes an error if a relative path tries to traverse outside that
	// prefix.
	SandboxPath *string `mapstructure:"sandbox_path"`

	// MapToEnvironmentVariable is the name of the environment variable this
	// template should map to. It is currently only used by Vault Agent and
	// will be ignored otherwise. When specified, Vault Agent will render the
	// contents of this template to the given environment variable instead
	// of a file. This field is mutually exclusive with `Destination`.
	MapToEnvironmentVariable *string `mapstructure:"-"`
}

// DefaultTemplateConfig returns a configuration that is populated with the
// default values.
func DefaultTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		Exec:       DefaultExecConfig(),
		Wait:       DefaultWaitConfig(),
		ExtFuncMap: make(template.FuncMap),
	}
}

// Copy returns a deep copy of this configuration.
func (c *TemplateConfig) Copy() *TemplateConfig {
	if c == nil {
		return nil
	}

	var o TemplateConfig

	o.Backup = c.Backup

	o.Command = c.Command

	o.CommandTimeout = c.CommandTimeout

	o.Contents = c.Contents

	o.CreateDestDirs = c.CreateDestDirs

	o.Destination = c.Destination

	o.ErrMissingKey = c.ErrMissingKey

	o.ErrFatal = c.ErrFatal

	if c.Exec != nil {
		o.Exec = c.Exec.Copy()
	}

	o.Perms = c.Perms

	o.Source = c.Source

	o.User = c.User
	o.Group = c.Group

	o.Uid = c.Uid
	o.Gid = c.Gid

	if c.Wait != nil {
		o.Wait = c.Wait.Copy()
	}

	o.LeftDelim = c.LeftDelim
	o.RightDelim = c.RightDelim

	if c.ExtFuncMap != nil {
		o.ExtFuncMap = make(template.FuncMap, len(c.ExtFuncMap))
		maps.Copy(o.ExtFuncMap, c.ExtFuncMap)
	}

	o.FunctionDenylist = append(o.FunctionDenylist, c.FunctionDenylist...)
	o.FunctionDenylistDeprecated = append(o.FunctionDenylistDeprecated, c.FunctionDenylistDeprecated...)

	o.SandboxPath = c.SandboxPath

	o.MapToEnvironmentVariable = c.MapToEnvironmentVariable

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *TemplateConfig) Merge(o *TemplateConfig) *TemplateConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Backup != nil {
		r.Backup = o.Backup
	}

	if o.Command != nil {
		r.Command = o.Command
	}

	if o.CommandTimeout != nil {
		r.CommandTimeout = o.CommandTimeout
	}

	if o.Contents != nil {
		r.Contents = o.Contents
	}

	if o.CreateDestDirs != nil {
		r.CreateDestDirs = o.CreateDestDirs
	}

	if o.Destination != nil {
		r.Destination = o.Destination
	}

	if o.ErrMissingKey != nil {
		r.ErrMissingKey = o.ErrMissingKey
	}

	if o.ErrFatal != nil {
		r.ErrFatal = o.ErrFatal
	}

	if o.Exec != nil {
		r.Exec = r.Exec.Merge(o.Exec)
	}

	if o.Perms != nil {
		r.Perms = o.Perms
	}

	if o.Source != nil {
		r.Source = o.Source
	}

	if o.User != nil {
		r.User = o.User
	}
	if o.Group != nil {
		r.Group = o.Group
	}

	if o.Uid != nil {
		r.Uid = o.Uid
	}

	if o.Gid != nil {
		r.Gid = o.Gid
	}

	if o.Wait != nil {
		r.Wait = r.Wait.Merge(o.Wait)
	}

	if o.LeftDelim != nil {
		r.LeftDelim = o.LeftDelim
	}

	if o.RightDelim != nil {
		r.RightDelim = o.RightDelim
	}

	if o.ExtFuncMap != nil {
		if r.ExtFuncMap == nil {
			r.ExtFuncMap = make(template.FuncMap, len(o.ExtFuncMap))
		}
		for key, fun := range o.ExtFuncMap {
			r.ExtFuncMap[key] = fun
		}
	}

	r.FunctionDenylist = append(r.FunctionDenylist, o.FunctionDenylist...)
	r.FunctionDenylistDeprecated = append(r.FunctionDenylistDeprecated, o.FunctionDenylistDeprecated...)

	if o.SandboxPath != nil {
		r.SandboxPath = o.SandboxPath
	}

	if o.MapToEnvironmentVariable != nil {
		r.MapToEnvironmentVariable = o.MapToEnvironmentVariable
	}

	return r
}

// Finalize ensures the configuration has no nil pointers and sets default
// values.
func (c *TemplateConfig) Finalize() {
	if c.Backup == nil {
		c.Backup = Bool(false)
	}

	if c.Command == nil {
		c.Command = []string{}
	}

	if c.CommandTimeout == nil {
		c.CommandTimeout = TimeDuration(DefaultTemplateCommandTimeout)
	}

	if c.Contents == nil {
		c.Contents = String("")
	}

	if c.CreateDestDirs == nil {
		c.CreateDestDirs = Bool(true)
	}

	if c.Destination == nil {
		c.Destination = String("")
	}

	if c.ErrMissingKey == nil {
		c.ErrMissingKey = Bool(false)
	}

	if c.ErrFatal == nil {
		c.ErrFatal = Bool(true)
	}

	// Backwards compatibility for uid
	if c.User == nil && c.Uid != nil {
		uStr := strconv.Itoa(*c.Uid)
		c.User = &uStr
	}

	// Backwards compatibility for gid
	if c.Group == nil && c.Gid != nil {
		gStr := strconv.Itoa(*c.Gid)
		c.Group = &gStr
	}

	if c.Exec == nil {
		c.Exec = DefaultExecConfig()
	}

	// Backwards compat for specifying command directly
	if c.Exec.Command == nil && c.Command != nil {
		c.Exec.Command = c.Command
	}
	// backwards compat with command_timeout and default support for exec.timeout
	switch {
	case c.Exec.Timeout == nil && c.CommandTimeout != nil:
		c.Exec.Timeout = c.CommandTimeout
	case c.Exec.Timeout == nil:
		c.Exec.Timeout = TimeDuration(DefaultTemplateCommandTimeout)
	}
	c.Exec.Finalize()

	if c.Perms == nil {
		c.Perms = FileMode(0)
	}

	if c.Source == nil {
		c.Source = String("")
	}

	if c.Wait == nil {
		c.Wait = DefaultWaitConfig()
	}
	c.Wait.Finalize()

	if c.LeftDelim == nil {
		c.LeftDelim = String("")
	}

	if c.RightDelim == nil {
		c.RightDelim = String("")
	}

	if c.SandboxPath == nil {
		c.SandboxPath = String("")
	}

	if c.ExtFuncMap == nil {
		c.ExtFuncMap = make(template.FuncMap, 0)
	}

	if c.FunctionDenylist == nil && c.FunctionDenylistDeprecated == nil {
		c.FunctionDenylist = []string{}
		c.FunctionDenylistDeprecated = []string{}
	} else {
		c.FunctionDenylist = combineLists(c.FunctionDenylist, c.FunctionDenylistDeprecated)
	}

	if c.MapToEnvironmentVariable == nil {
		c.MapToEnvironmentVariable = String("")
	}
}

// GoString defines the printable version of this struct.
func (c *TemplateConfig) GoString() string {
	if c == nil {
		return "(*TemplateConfig)(nil)"
	}

	return fmt.Sprintf("&TemplateConfig{"+
		"Backup:%s, "+
		"Command:%s, "+
		"CommandTimeout:%s, "+
		"Contents:%s, "+
		"CreateDestDirs:%s, "+
		"Destination:%s, "+
		"ErrMissingKey:%s, "+
		"ErrFatal:%s, "+
		"Exec:%#v, "+
		"Perms:%s, "+
		"Source:%s, "+
		"Wait:%#v, "+
		"LeftDelim:%s, "+
		"RightDelim:%s, "+
		"ExtFuncMap:%s, "+
		"FunctionDenylist:%s, "+
		"SandboxPath:%s "+
		"MapToEnvironmentVariable:%s"+
		"}",
		BoolGoString(c.Backup),
		c.Command,
		TimeDurationGoString(c.CommandTimeout),
		StringGoString(c.Contents),
		BoolGoString(c.CreateDestDirs),
		StringGoString(c.Destination),
		BoolGoString(c.ErrMissingKey),
		BoolGoString(c.ErrFatal),
		c.Exec,
		FileModeGoString(c.Perms),
		StringGoString(c.Source),
		c.Wait,
		StringGoString(c.LeftDelim),
		StringGoString(c.RightDelim),
		maps.Keys(c.ExtFuncMap),
		combineLists(c.FunctionDenylist, c.FunctionDenylistDeprecated),
		StringGoString(c.SandboxPath),
		StringGoString(c.MapToEnvironmentVariable),
	)
}

// Display is the human-friendly form of this configuration. It tries to
// describe this template in as much detail as possible in a single line, so
// log consumers can uniquely identify it.
func (c *TemplateConfig) Display() string {
	if c == nil {
		return ""
	}

	source := c.Source
	if StringPresent(c.Contents) {
		source = String("(dynamic)")
	}

	destination := c.Destination
	if StringPresent(c.MapToEnvironmentVariable) {
		destination = c.MapToEnvironmentVariable
	}

	return fmt.Sprintf("%q => %q",
		StringVal(source),
		StringVal(destination),
	)
}

// TemplateConfigs is a collection of TemplateConfigs
type TemplateConfigs []*TemplateConfig

// DefaultTemplateConfigs returns a configuration that is populated with the
// default values.
func DefaultTemplateConfigs() *TemplateConfigs {
	return &TemplateConfigs{}
}

// Copy returns a deep copy of this configuration.
func (c *TemplateConfigs) Copy() *TemplateConfigs {
	o := make(TemplateConfigs, len(*c))
	for i, t := range *c {
		o[i] = t.Copy()
	}
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *TemplateConfigs) Merge(o *TemplateConfigs) *TemplateConfigs {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	*r = append(*r, *o...)

	return r
}

// Finalize ensures the configuration has no nil pointers and sets default
// values.
func (c *TemplateConfigs) Finalize() {
	if c == nil {
		*c = *DefaultTemplateConfigs()
	}

	for _, t := range *c {
		t.Finalize()
	}
}

// GoString defines the printable version of this struct.
func (c *TemplateConfigs) GoString() string {
	if c == nil {
		return "(*TemplateConfigs)(nil)"
	}

	s := make([]string, len(*c))
	for i, t := range *c {
		s[i] = t.GoString()
	}

	return "{" + strings.Join(s, ", ") + "}"
}

// ParseTemplateConfig parses a string in the form source:destination:command
// into a TemplateConfig.
func ParseTemplateConfig(s string) (*TemplateConfig, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, ErrTemplateStringEmpty
	}

	var source, destination, command string
	parts := configTemplateRe.FindAllString(s, -1)

	switch len(parts) {
	case 1:
		source = parts[0]
	case 2:
		source, destination = parts[0], parts[1]
	case 3:
		source, destination, command = parts[0], parts[1], parts[2]
	default:
		source, destination = parts[0], parts[1]
		command = strings.Join(parts[2:], ":")
	}

	var sourcePtr, destinationPtr *string
	var commandL commandList
	if source != "" {
		sourcePtr = String(source)
	}
	if destination != "" {
		destinationPtr = String(destination)
	}
	if command != "" {
		commandL = []string{command}
	}

	return &TemplateConfig{
		Source:      sourcePtr,
		Destination: destinationPtr,
		Command:     commandL,
	}, nil
}
