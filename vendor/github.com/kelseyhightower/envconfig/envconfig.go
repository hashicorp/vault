// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidSpecification indicates that a specification is of the wrong type.
var ErrInvalidSpecification = errors.New("specification must be a struct pointer")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct {
	KeyName   string
	FieldName string
	TypeName  string
	Value     string
	Err       error
}

// Decoder has the same semantics as Setter, but takes higher precedence.
// It is provided for historical compatibility.
type Decoder interface {
	Decode(value string) error
}

// Setter is implemented by types can self-deserialize values.
// Any type that implements flag.Value also implements Setter.
type Setter interface {
	Set(value string) error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("envconfig.Process: assigning %[1]s to %[2]s: converting '%[3]s' to type %[4]s. details: %[5]s", e.KeyName, e.FieldName, e.Value, e.TypeName, e.Err)
}

// varInfo maintains information about the configuration variable
type varInfo struct {
	Name  string
	Alt   string
	Key   string
	Field reflect.Value
	Tags  reflect.StructTag
}

// GatherInfo gathers information about the specified struct
func gatherInfo(prefix string, spec interface{}) ([]varInfo, error) {
	expr := regexp.MustCompile("([^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")
	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr {
		return nil, ErrInvalidSpecification
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, ErrInvalidSpecification
	}
	typeOfSpec := s.Type()

	// over allocate an info array, we will extend if needed later
	infos := make([]varInfo, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ftype := typeOfSpec.Field(i)
		if !f.CanSet() || ftype.Tag.Get("ignored") == "true" {
			continue
		}

		for f.Kind() == reflect.Ptr {
			if f.IsNil() {
				if f.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone
					break
				}
				// nil pointer to struct: create a zero instance
				f.Set(reflect.New(f.Type().Elem()))
			}
			f = f.Elem()
		}

		// Capture information about the config variable
		info := varInfo{
			Name:  ftype.Name,
			Field: f,
			Tags:  ftype.Tag,
			Alt:   strings.ToUpper(ftype.Tag.Get("envconfig")),
		}

		// Default to the field name as the env var name (will be upcased)
		info.Key = info.Name

		// Best effort to un-pick camel casing as separate words
		if ftype.Tag.Get("split_words") == "true" {
			words := expr.FindAllStringSubmatch(ftype.Name, -1)
			if len(words) > 0 {
				var name []string
				for _, words := range words {
					name = append(name, words[0])
				}

				info.Key = strings.Join(name, "_")
			}
		}
		if info.Alt != "" {
			info.Key = info.Alt
		}
		if prefix != "" {
			info.Key = fmt.Sprintf("%s_%s", prefix, info.Key)
		}
		info.Key = strings.ToUpper(info.Key)
		infos = append(infos, info)

		if f.Kind() == reflect.Struct {
			// honor Decode if present
			if decoderFrom(f) == nil && setterFrom(f) == nil && textUnmarshaler(f) == nil {
				innerPrefix := prefix
				if !ftype.Anonymous {
					innerPrefix = info.Key
				}

				embeddedPtr := f.Addr().Interface()
				embeddedInfos, err := gatherInfo(innerPrefix, embeddedPtr)
				if err != nil {
					return nil, err
				}
				infos = append(infos[:len(infos)-1], embeddedInfos...)

				continue
			}
		}
	}
	return infos, nil
}

// Process populates the specified struct based on environment variables
func Process(prefix string, spec interface{}) error {
	infos, err := gatherInfo(prefix, spec)

	for _, info := range infos {

		// `os.Getenv` cannot differentiate between an explicitly set empty value
		// and an unset value. `os.LookupEnv` is preferred to `syscall.Getenv`,
		// but it is only available in go1.5 or newer. We're using Go build tags
		// here to use os.LookupEnv for >=go1.5
		value, ok := lookupEnv(info.Key)
		if !ok && info.Alt != "" {
			value, ok = lookupEnv(info.Alt)
		}

		def := info.Tags.Get("default")
		if def != "" && !ok {
			value = def
		}

		req := info.Tags.Get("required")
		if !ok && def == "" {
			if req == "true" {
				return fmt.Errorf("required key %s missing value", info.Key)
			}
			continue
		}

		err := processField(value, info.Field)
		if err != nil {
			return &ParseError{
				KeyName:   info.Key,
				FieldName: info.Name,
				TypeName:  info.Field.Type().String(),
				Value:     value,
				Err:       err,
			}
		}
	}

	return err
}

// MustProcess is the same as Process but panics if an error occurs
func MustProcess(prefix string, spec interface{}) {
	if err := Process(prefix, spec); err != nil {
		panic(err)
	}
}

func processField(value string, field reflect.Value) error {
	typ := field.Type()

	decoder := decoderFrom(field)
	if decoder != nil {
		return decoder.Decode(value)
	}
	// look for Set method if Decode not defined
	setter := setterFrom(field)
	if setter != nil {
		return setter.Set(value)
	}

	if t := textUnmarshaler(field); t != nil {
		return t.UnmarshalText([]byte(value))
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		}
		if err != nil {
			return err
		}

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, typ.Bits())
		if err != nil {
			return err
		}
		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, typ.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(val)
	case reflect.Slice:
		vals := strings.Split(value, ",")
		sl := reflect.MakeSlice(typ, len(vals), len(vals))
		for i, val := range vals {
			err := processField(val, sl.Index(i))
			if err != nil {
				return err
			}
		}
		field.Set(sl)
	case reflect.Map:
		pairs := strings.Split(value, ",")
		mp := reflect.MakeMap(typ)
		for _, pair := range pairs {
			kvpair := strings.Split(pair, ":")
			if len(kvpair) != 2 {
				return fmt.Errorf("invalid map item: %q", pair)
			}
			k := reflect.New(typ.Key()).Elem()
			err := processField(kvpair[0], k)
			if err != nil {
				return err
			}
			v := reflect.New(typ.Elem()).Elem()
			err = processField(kvpair[1], v)
			if err != nil {
				return err
			}
			mp.SetMapIndex(k, v)
		}
		field.Set(mp)
	}

	return nil
}

func interfaceFrom(field reflect.Value, fn func(interface{}, *bool)) {
	// it may be impossible for a struct field to fail this check
	if !field.CanInterface() {
		return
	}
	var ok bool
	fn(field.Interface(), &ok)
	if !ok && field.CanAddr() {
		fn(field.Addr().Interface(), &ok)
	}
}

func decoderFrom(field reflect.Value) (d Decoder) {
	interfaceFrom(field, func(v interface{}, ok *bool) { d, *ok = v.(Decoder) })
	return d
}

func setterFrom(field reflect.Value) (s Setter) {
	interfaceFrom(field, func(v interface{}, ok *bool) { s, *ok = v.(Setter) })
	return s
}

func textUnmarshaler(field reflect.Value) (t encoding.TextUnmarshaler) {
	interfaceFrom(field, func(v interface{}, ok *bool) { t, *ok = v.(encoding.TextUnmarshaler) })
	return t
}
