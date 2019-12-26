package command

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/posener/complete"
)

// FlagExample is an interface which declares an example value.
type FlagExample interface {
	Example() string
}

// FlagVisibility is an interface which declares whether a flag should be
// hidden from help and completions. This is usually used for deprecations
// on "internal-only" flags.
type FlagVisibility interface {
	Hidden() bool
}

// FlagBool is an interface which boolean flags implement.
type FlagBool interface {
	IsBoolFlag() bool
}

// -- BoolVar  and boolValue
type BoolVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    bool
	Hidden     bool
	EnvVar     string
	Target     *bool
	Completion complete.Predictor
}

func (f *FlagSet) BoolVar(i *BoolVar) {
	def := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if b, err := strconv.ParseBool(v); err == nil {
			def = b
		}
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    strconv.FormatBool(i.Default),
		EnvVar:     i.EnvVar,
		Value:      newBoolValue(def, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type boolValue struct {
	hidden bool
	target *bool
}

func newBoolValue(def bool, target *bool, hidden bool) *boolValue {
	*target = def

	return &boolValue{
		hidden: hidden,
		target: target,
	}
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*b.target = v
	return nil
}

func (b *boolValue) Get() interface{} { return *b.target }
func (b *boolValue) String() string   { return strconv.FormatBool(*b.target) }
func (b *boolValue) Example() string  { return "" }
func (b *boolValue) Hidden() bool     { return b.hidden }
func (b *boolValue) IsBoolFlag() bool { return true }

// -- IntVar and intValue
type IntVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    int
	Hidden     bool
	EnvVar     string
	Target     *int
	Completion complete.Predictor
}

func (f *FlagSet) IntVar(i *IntVar) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if i, err := strconv.ParseInt(v, 0, 64); err == nil {
			initial = int(i)
		}
	}

	def := ""
	if i.Default != 0 {
		def = strconv.FormatInt(int64(i.Default), 10)
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newIntValue(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type intValue struct {
	hidden bool
	target *int
}

func newIntValue(def int, target *int, hidden bool) *intValue {
	*target = def
	return &intValue{
		hidden: hidden,
		target: target,
	}
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	*i.target = int(v)
	return nil
}

func (i *intValue) Get() interface{} { return int(*i.target) }
func (i *intValue) String() string   { return strconv.Itoa(int(*i.target)) }
func (i *intValue) Example() string  { return "int" }
func (i *intValue) Hidden() bool     { return i.hidden }

// -- Int64Var and int64Value
type Int64Var struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    int64
	Hidden     bool
	EnvVar     string
	Target     *int64
	Completion complete.Predictor
}

func (f *FlagSet) Int64Var(i *Int64Var) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if i, err := strconv.ParseInt(v, 0, 64); err == nil {
			initial = i
		}
	}

	def := ""
	if i.Default != 0 {
		def = strconv.FormatInt(int64(i.Default), 10)
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newInt64Value(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type int64Value struct {
	hidden bool
	target *int64
}

func newInt64Value(def int64, target *int64, hidden bool) *int64Value {
	*target = def
	return &int64Value{
		hidden: hidden,
		target: target,
	}
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	*i.target = v
	return nil
}

func (i *int64Value) Get() interface{} { return int64(*i.target) }
func (i *int64Value) String() string   { return strconv.FormatInt(int64(*i.target), 10) }
func (i *int64Value) Example() string  { return "int" }
func (i *int64Value) Hidden() bool     { return i.hidden }

// -- UintVar && uintValue
type UintVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    uint
	Hidden     bool
	EnvVar     string
	Target     *uint
	Completion complete.Predictor
}

func (f *FlagSet) UintVar(i *UintVar) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if i, err := strconv.ParseUint(v, 0, 64); err == nil {
			initial = uint(i)
		}
	}

	def := ""
	if i.Default != 0 {
		def = strconv.FormatUint(uint64(i.Default), 10)
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newUintValue(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type uintValue struct {
	hidden bool
	target *uint
}

func newUintValue(def uint, target *uint, hidden bool) *uintValue {
	*target = def
	return &uintValue{
		hidden: hidden,
		target: target,
	}
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}

	*i.target = uint(v)
	return nil
}

func (i *uintValue) Get() interface{} { return uint(*i.target) }
func (i *uintValue) String() string   { return strconv.FormatUint(uint64(*i.target), 10) }
func (i *uintValue) Example() string  { return "uint" }
func (i *uintValue) Hidden() bool     { return i.hidden }

// -- Uint64Var and uint64Value
type Uint64Var struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    uint64
	Hidden     bool
	EnvVar     string
	Target     *uint64
	Completion complete.Predictor
}

func (f *FlagSet) Uint64Var(i *Uint64Var) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if i, err := strconv.ParseUint(v, 0, 64); err == nil {
			initial = i
		}
	}

	def := ""
	if i.Default != 0 {
		strconv.FormatUint(i.Default, 10)
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newUint64Value(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type uint64Value struct {
	hidden bool
	target *uint64
}

func newUint64Value(def uint64, target *uint64, hidden bool) *uint64Value {
	*target = def
	return &uint64Value{
		hidden: hidden,
		target: target,
	}
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}

	*i.target = v
	return nil
}

func (i *uint64Value) Get() interface{} { return uint64(*i.target) }
func (i *uint64Value) String() string   { return strconv.FormatUint(uint64(*i.target), 10) }
func (i *uint64Value) Example() string  { return "uint" }
func (i *uint64Value) Hidden() bool     { return i.hidden }

// -- StringVar and stringValue
type StringVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    string
	Hidden     bool
	EnvVar     string
	Target     *string
	Completion complete.Predictor
}

func (f *FlagSet) StringVar(i *StringVar) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		initial = v
	}

	def := ""
	if i.Default != "" {
		def = i.Default
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newStringValue(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type stringValue struct {
	hidden bool
	target *string
}

func newStringValue(def string, target *string, hidden bool) *stringValue {
	*target = def
	return &stringValue{
		hidden: hidden,
		target: target,
	}
}

func (s *stringValue) Set(val string) error {
	*s.target = val
	return nil
}

func (s *stringValue) Get() interface{} { return *s.target }
func (s *stringValue) String() string   { return *s.target }
func (s *stringValue) Example() string  { return "string" }
func (s *stringValue) Hidden() bool     { return s.hidden }

// -- Float64Var and float64Value
type Float64Var struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    float64
	Hidden     bool
	EnvVar     string
	Target     *float64
	Completion complete.Predictor
}

func (f *FlagSet) Float64Var(i *Float64Var) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if i, err := strconv.ParseFloat(v, 64); err == nil {
			initial = i
		}
	}

	def := ""
	if i.Default != 0 {
		def = strconv.FormatFloat(i.Default, 'e', -1, 64)
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newFloat64Value(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type float64Value struct {
	hidden bool
	target *float64
}

func newFloat64Value(def float64, target *float64, hidden bool) *float64Value {
	*target = def
	return &float64Value{
		hidden: hidden,
		target: target,
	}
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f.target = v
	return nil
}

func (f *float64Value) Get() interface{} { return float64(*f.target) }
func (f *float64Value) String() string   { return strconv.FormatFloat(float64(*f.target), 'g', -1, 64) }
func (f *float64Value) Example() string  { return "float" }
func (f *float64Value) Hidden() bool     { return f.hidden }

// -- DurationVar and durationValue
type DurationVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    time.Duration
	Hidden     bool
	EnvVar     string
	Target     *time.Duration
	Completion complete.Predictor
}

func (f *FlagSet) DurationVar(i *DurationVar) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if d, err := time.ParseDuration(appendDurationSuffix(v)); err == nil {
			initial = d
		}
	}

	def := ""
	if i.Default != 0 {
		def = i.Default.String()
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newDurationValue(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type durationValue struct {
	hidden bool
	target *time.Duration
}

func newDurationValue(def time.Duration, target *time.Duration, hidden bool) *durationValue {
	*target = def
	return &durationValue{
		hidden: hidden,
		target: target,
	}
}

func (d *durationValue) Set(s string) error {
	// Maintain bc for people specifying "system" as the value.
	if s == "system" {
		s = "-1"
	}

	v, err := time.ParseDuration(appendDurationSuffix(s))
	if err != nil {
		return err
	}
	*d.target = v
	return nil
}

func (d *durationValue) Get() interface{} { return time.Duration(*d.target) }
func (d *durationValue) String() string   { return (*d.target).String() }
func (d *durationValue) Example() string  { return "duration" }
func (d *durationValue) Hidden() bool     { return d.hidden }

// appendDurationSuffix is used as a backwards-compat tool for assuming users
// meant "seconds" when they do not provide a suffixed duration value.
func appendDurationSuffix(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "m") || strings.HasSuffix(s, "h") {
		return s
	}
	return s + "s"
}

// -- StringSliceVar and stringSliceValue
type StringSliceVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    []string
	Hidden     bool
	EnvVar     string
	Target     *[]string
	Completion complete.Predictor
}

func (f *FlagSet) StringSliceVar(i *StringSliceVar) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		initial = parts
	}

	def := ""
	if i.Default != nil {
		def = strings.Join(i.Default, ",")
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newStringSliceValue(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type stringSliceValue struct {
	hidden bool
	target *[]string
}

func newStringSliceValue(def []string, target *[]string, hidden bool) *stringSliceValue {
	*target = def
	return &stringSliceValue{
		hidden: hidden,
		target: target,
	}
}

func (s *stringSliceValue) Set(val string) error {
	*s.target = append(*s.target, strings.TrimSpace(val))
	return nil
}

func (s *stringSliceValue) Get() interface{} { return *s.target }
func (s *stringSliceValue) String() string   { return strings.Join(*s.target, ",") }
func (s *stringSliceValue) Example() string  { return "string" }
func (s *stringSliceValue) Hidden() bool     { return s.hidden }

// -- StringMapVar and stringMapValue
type StringMapVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    map[string]string
	Hidden     bool
	Target     *map[string]string
	Completion complete.Predictor
}

func (f *FlagSet) StringMapVar(i *StringMapVar) {
	def := ""
	if i.Default != nil {
		def = mapToKV(i.Default)
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		Value:      newStringMapValue(i.Default, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type stringMapValue struct {
	hidden bool
	target *map[string]string
}

func newStringMapValue(def map[string]string, target *map[string]string, hidden bool) *stringMapValue {
	*target = def
	return &stringMapValue{
		hidden: hidden,
		target: target,
	}
}

func (s *stringMapValue) Set(val string) error {
	idx := strings.Index(val, "=")
	if idx == -1 {
		return fmt.Errorf("missing = in KV pair: %q", val)
	}

	if *s.target == nil {
		*s.target = make(map[string]string)
	}

	k, v := val[0:idx], val[idx+1:]
	(*s.target)[k] = v
	return nil
}

func (s *stringMapValue) Get() interface{} { return *s.target }
func (s *stringMapValue) String() string   { return mapToKV(*s.target) }
func (s *stringMapValue) Example() string  { return "key=value" }
func (s *stringMapValue) Hidden() bool     { return s.hidden }

func mapToKV(m map[string]string) string {
	list := make([]string, 0, len(m))
	for k, _ := range m {
		list = append(list, k)
	}
	sort.Strings(list)

	for i, k := range list {
		list[i] = k + "=" + m[k]
	}

	return strings.Join(list, ",")
}

// -- VarFlag
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
	// If the flag is marked as hidden, just add it to the set and return to
	// avoid unnecessary computations here. We do not want to add completions or
	// generate help output for hidden flags.
	if v, ok := i.Value.(FlagVisibility); ok && v.Hidden() {
		f.Var(i.Value, i.Name, "")
		return
	}

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

	// Add aliases to the main set
	for _, a := range i.Aliases {
		f.mainSet.Var(i.Value, a, "")
	}

	f.Var(i.Value, i.Name, usage)
	f.completions["-"+i.Name] = i.Completion
}

// Var is a lower-level API for adding something to the flags. It should be used
// with caution, since it bypasses all validation. Consider VarFlag instead.
func (f *FlagSet) Var(value flag.Value, name, usage string) {
	f.mainSet.Var(value, name, usage)
	f.flagSet.Var(value, name, usage)
}

// -- helpers
func envDefault(key, def string) string {
	if v, exist := os.LookupEnv(key); exist {
		return v
	}
	return def
}

func envBoolDefault(key string, def bool) bool {
	if v, exist := os.LookupEnv(key); exist {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		return b
	}
	return def
}

func envDurationDefault(key string, def time.Duration) time.Duration {
	if v, exist := os.LookupEnv(key); exist {
		d, err := time.ParseDuration(v)
		if err != nil {
			panic(err)
		}
		return d
	}
	return def
}
