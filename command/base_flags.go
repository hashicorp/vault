package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/posener/complete"
)

// FlagExample is an interface which declares an example value.
type FlagExample interface {
	Example() string
}

type BoolVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    bool
	EnvVar     string
	Target     *bool
	Completion complete.Predictor
}

func (f *FlagSet) BoolVar(i *BoolVar) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if b, err := strconv.ParseBool(v); err != nil {
			def = b
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatBool(i.Default),
		EnvVar:     i.EnvVar,
		Value:      newBoolValue(def, i.Target),
		Completion: i.Completion,
	})
}

type IntVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    int
	EnvVar     string
	Target     *int
	Completion complete.Predictor
}

func (f *FlagSet) IntVar(i *IntVar) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if i, err := strconv.ParseInt(v, 0, 64); err != nil {
			def = int(i)
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatInt(int64(i.Default), 10),
		EnvVar:     i.EnvVar,
		Value:      newIntValue(def, i.Target),
		Completion: i.Completion,
	})
}

type Int64Var struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    int64
	EnvVar     string
	Target     *int64
	Completion complete.Predictor
}

func (f *FlagSet) Int64Var(i *Int64Var) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if i, err := strconv.ParseInt(v, 0, 64); err != nil {
			def = i
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatInt(i.Default, 10),
		EnvVar:     i.EnvVar,
		Value:      newInt64Value(def, i.Target),
		Completion: i.Completion,
	})
}

type UintVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    uint
	EnvVar     string
	Target     *uint
	Completion complete.Predictor
}

func (f *FlagSet) UintVar(i *UintVar) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if i, err := strconv.ParseUint(v, 0, 64); err != nil {
			def = uint(i)
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatUint(uint64(i.Default), 10),
		EnvVar:     i.EnvVar,
		Value:      newUintValue(def, i.Target),
		Completion: i.Completion,
	})
}

type Uint64Var struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    uint64
	EnvVar     string
	Target     *uint64
	Completion complete.Predictor
}

func (f *FlagSet) Uint64Var(i *Uint64Var) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if i, err := strconv.ParseUint(v, 0, 64); err != nil {
			def = i
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatUint(i.Default, 10),
		EnvVar:     i.EnvVar,
		Value:      newUint64Value(def, i.Target),
		Completion: i.Completion,
	})
}

type StringVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    string
	EnvVar     string
	Target     *string
	Completion complete.Predictor
}

func (f *FlagSet) StringVar(i *StringVar) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		def = v
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    i.Default,
		EnvVar:     i.EnvVar,
		Value:      newStringValue(def, i.Target),
		Completion: i.Completion,
	})
}

type Float64Var struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    float64
	EnvVar     string
	Target     *float64
	Completion complete.Predictor
}

func (f *FlagSet) Float64Var(i *Float64Var) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if i, err := strconv.ParseFloat(v, 64); err != nil {
			def = i
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatFloat(i.Default, 'e', -1, 64),
		EnvVar:     i.EnvVar,
		Value:      newFloat64Value(def, i.Target),
		Completion: i.Completion,
	})
}

type DurationVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    time.Duration
	EnvVar     string
	Target     *time.Duration
	Completion complete.Predictor
}

func (f *FlagSet) DurationVar(i *DurationVar) {
	def := i.Default
	if v := os.Getenv(i.EnvVar); v != "" {
		if d, err := time.ParseDuration(v); err != nil {
			def = d
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    i.Default.String(),
		EnvVar:     i.EnvVar,
		Value:      newDurationValue(def, i.Target),
		Completion: i.Completion,
	})
}

type VarFlag struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    string
	EnvVar     string
	Value      flag.Value
	Completion complete.Predictor
}

func (f *FlagSet) VarFlag(i *VarFlag) {
	// Calculate the full usage
	usage := i.Usage

	if len(i.Aliases) > 0 {
		sentence := make([]string, len(i.Aliases))
		for i, a := range i.Aliases {
			sentence[i] = fmt.Sprintf(`"-%s"`, a)
		}

		aliases := ""
		switch len(sentence) {
		case 0:
			// impossible...
		case 1:
			aliases = sentence[0]
		case 2:
			aliases = sentence[0] + " and " + sentence[1]
		default:
			sentence[len(sentence)-1] = "and " + sentence[len(sentence)-1]
			aliases = strings.Join(sentence, ", ")
		}

		usage += fmt.Sprintf(" This is aliased as %s.", aliases)
	}

	if i.Default != "" {
		usage += fmt.Sprintf(" The default is %s.", i.Default)
	}

	if i.EnvVar != "" {
		usage += fmt.Sprintf(" This can also be specified via the %s "+
			"environment variable.", i.EnvVar)
	}

	f.mainSet.Var(i.Value, i.Name, "") // No point in passing along usage here

	// Add aliases to the main set
	for _, a := range i.Aliases {
		f.mainSet.Var(i.Value, a, "")
	}

	f.flagSet.Var(i.Value, i.Name, usage)
	f.completions["-"+i.Name] = i.Completion
}

func (f *FlagSet) Var(value flag.Value, name, usage string) {
	f.mainSet.Var(value, name, usage)
	f.flagSet.Var(value, name, usage)
}

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) Example() string { return "" }

func (b *boolValue) IsBoolFlag() bool { return true }

// optional interface to indicate boolean flags that can be
// supplied without "=value" text
type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func (i *intValue) Example() string { return "int" }

// -- int64 Value
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() interface{} { return int64(*i) }

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func (i *int64Value) Example() string { return "int" }

// -- uint Value
type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() interface{} { return uint(*i) }

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

func (i *uintValue) Example() string { return "uint" }

// -- uint64 Value
type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() interface{} { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func (i *uint64Value) Example() string { return "uint" }

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

func (s *stringValue) Example() string { return "string" }

// -- float64 Value
type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() interface{} { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

func (f *float64Value) Example() string { return "float" }

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

func (d *durationValue) Example() string { return "duration" }

// -- helpers
func envDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBoolDefault(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		return b
	}
	return def
}

func envDurationDefault(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			panic(err)
		}
		return d
	}
	return def
}
