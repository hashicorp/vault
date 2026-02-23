// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	slogctx "github.com/veqryn/slog-context"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

// Config represents the pipeline utilities configuration file.
type Config struct {
	ChangedFiles *changed.Config `hcl:"changed_files,block"`
}

// DecodeReq contains the parameters for decoding a pipeline configuration.
type DecodeReq struct {
	// Path is the explicit path to the pipeline.hcl file. If empty, RepoRoot
	// will be used to construct the default path.
	Path string `json:"path,omitempty"`
}

// DecodeRes contains the result of decoding a pipeline configuration.
type DecodeRes struct {
	// Path is the absolute path to the configuration file that was loaded.
	Path string `json:"path,omitempty"`

	// Config is the decoded pipeline configuration.
	Config *Config `json:"config,omitempty"`

	// Err contains the error if decoding failed.
	Err error

	// ErrStr contains the error message if decoding failed.
	ErrStr string `json:"error,omitempty"`
}

// Decode decodes a pipeline configuration from the specified path.
// It returns a DecodeRes containing the absolute path used, the decoded
// config, and any error that occurred. If you do not need the surrounding
// DecodeRes struct you can use Load() instead.
func Decode(ctx context.Context, req *DecodeReq) *DecodeRes {
	res := &DecodeRes{}

	// Determine the path to use
	if req.Path == "" {
		res.Err = errors.New("no path or repo root provided")
		res.ErrStr = res.Err.Error()
	}

	// Get absolute path
	absPath, err := filepath.Abs(req.Path)
	if err != nil {
		res.Err = err
		res.ErrStr = err.Error()

		return res
	}
	res.Path = absPath

	// Load the configuration
	cfg, err := Load(ctx, absPath)
	if err != nil {
		res.Err = err
		res.ErrStr = err.Error()

		return res
	}
	res.Config = cfg

	return res
}

// Load takes a context and variable opts and locates the pipeline configuration
// file, parses it, and returns the Config.
func Load(ctx context.Context, path string) (*Config, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ctx = slogctx.Append(ctx,
		slog.String("path", path),
	)

	slog.Default().DebugContext(ctx, "loading pipeline configuration")
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return DecodeFile(path)
}

// DecodeFile decodes the configuration file at the given path into a Config.
func DecodeFile(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return DecodeBytes(bytes)
}

// DecodeBytes decodes the configuration from bytes into a Config.
func DecodeBytes(body []byte) (*Config, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(body, "pipeline.hcl")
	if diags.HasErrors() {
		return nil, errors.Join(diags.Errs()...)
	}

	spec := &Config{}
	moreDiags := gohcl.DecodeBody(file.Body, evalContext(), spec)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, errors.Join(diags.Errs()...)
	}

	return spec, nil
}

// Validate validates that the response is initialized and does not have
// an embedded error.
func (d *DecodeRes) Validate(ctx context.Context) error {
	if d == nil {
		return errors.New("pipeline.hcl decode response is unitialized")
	}

	if d.Err != nil {
		return d.Err
	}

	return nil
}

// evalContext is the default eval context for the decoder. We add many stdlib
// functions and our own that handles joining paths.
func evalContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Functions: map[string]function.Function{
			"abs":                    stdlib.AbsoluteFunc,
			"absolute":               stdlib.AbsoluteFunc,
			"add":                    stdlib.AddFunc,
			"and":                    stdlib.AndFunc,
			"byteslen":               stdlib.BytesLenFunc,
			"bytessclice":            stdlib.BytesSliceFunc,
			"ceil":                   stdlib.CeilFunc,
			"chomp":                  stdlib.ChompFunc,
			"chunklist":              stdlib.ChunklistFunc,
			"coalesce":               stdlib.CoalesceFunc,
			"coalescelist":           stdlib.CoalesceListFunc,
			"compact":                stdlib.CompactFunc,
			"concat":                 stdlib.ConcatFunc,
			"contains":               stdlib.ContainsFunc,
			"csvdecode":              stdlib.CSVDecodeFunc,
			"distinct":               stdlib.DistinctFunc,
			"divide":                 stdlib.DivideFunc,
			"element":                stdlib.ElementFunc,
			"equal":                  stdlib.EqualFunc,
			"flatten":                stdlib.FlattenFunc,
			"floor":                  stdlib.FloorFunc,
			"format":                 stdlib.FormatFunc,
			"formatdate":             stdlib.FormatDateFunc,
			"formatlist":             stdlib.FormatListFunc,
			"greaterthan":            stdlib.GreaterThanFunc,
			"greaterthanorequalto":   stdlib.GreaterThanOrEqualToFunc,
			"hasindex":               stdlib.HasIndexFunc,
			"indent":                 stdlib.IndentFunc,
			"index":                  stdlib.IndexFunc,
			"int":                    stdlib.IntFunc,
			"jsondecode":             stdlib.JSONDecodeFunc,
			"jsonencode":             stdlib.JSONEncodeFunc,
			"join":                   stdlib.JoinFunc,
			"joinpath":               joinPath,
			"keys":                   stdlib.KeysFunc,
			"length":                 stdlib.LengthFunc,
			"lessthan":               stdlib.LessThanFunc,
			"lessthanorequalto":      stdlib.LessThanOrEqualToFunc,
			"log":                    stdlib.LogFunc,
			"lookup":                 stdlib.LookupFunc,
			"lower":                  stdlib.LowerFunc,
			"max":                    stdlib.MaxFunc,
			"merge":                  stdlib.MergeFunc,
			"min":                    stdlib.MinFunc,
			"modulo":                 stdlib.ModuloFunc,
			"multiply":               stdlib.MultiplyFunc,
			"negate":                 stdlib.NegateFunc,
			"not":                    stdlib.NotFunc,
			"notequal":               stdlib.NotEqualFunc,
			"or":                     stdlib.OrFunc,
			"parseint":               stdlib.ParseIntFunc,
			"pow":                    stdlib.PowFunc,
			"range":                  stdlib.RangeFunc,
			"regex":                  stdlib.RegexFunc,
			"regexall":               stdlib.RegexAllFunc,
			"regexreplace":           stdlib.RegexReplaceFunc,
			"replace":                stdlib.ReplaceFunc,
			"reverse":                stdlib.ReverseFunc,
			"reverselist":            stdlib.ReverseListFunc,
			"sethaselement":          stdlib.SetHasElementFunc,
			"setintersection":        stdlib.SetIntersectionFunc,
			"setproduct":             stdlib.SetProductFunc,
			"setsubtract":            stdlib.SetSubtractFunc,
			"setsymmetricdifference": stdlib.SetSymmetricDifferenceFunc,
			"setunion":               stdlib.SetUnionFunc,
			"signum":                 stdlib.SignumFunc,
			"slice":                  stdlib.SliceFunc,
			"sort":                   stdlib.SortFunc,
			"split":                  stdlib.SplitFunc,
			"strlen":                 stdlib.StrlenFunc,
			"substr":                 stdlib.SubstrFunc,
			"subtract":               stdlib.SubtractFunc,
			"timeadd":                stdlib.TimeAddFunc,
			"title":                  stdlib.TitleFunc,
			"trim":                   stdlib.TrimFunc,
			"trimprefix":             stdlib.TrimPrefixFunc,
			"trimspace":              stdlib.TrimSpaceFunc,
			"trimsuffix":             stdlib.TrimSuffixFunc,
			"upper":                  stdlib.UpperFunc,
			"values":                 stdlib.ValuesFunc,
			"zipmap":                 stdlib.ZipmapFunc,
		},
	}
}

// joinPath joins variable string arguments to together as a path with the OS
// path separator.
var joinPath = function.New(&function.Spec{
	Description: "Joins the paths with the OS specific separator",
	VarParam: &function.Parameter{
		Name:        "paths",
		Description: "One or more lists of strings to join.",
		Type:        cty.String,
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		items := make([]string, 0, len(args))
		for _, arg := range args {
			if !arg.IsWhollyKnown() {
				return cty.UnknownVal(cty.String), nil
			}
			items = append(items, arg.AsString())
		}

		return cty.StringVal(filepath.ToSlash(strings.Join(items, string(os.PathSeparator)))), nil
	},
})
