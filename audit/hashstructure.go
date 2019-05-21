package audit

import (
	"errors"
	"reflect"
	"time"

	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
	"github.com/mitchellh/reflectwalk"
)

// HashString hashes the given opaque string and returns it
func HashString(salter *salt.Salt, data string) string {
	return salter.GetIdentifiedHMAC(data)
}

// Hash will hash the given type. This has built-in support for auth,
// requests, and responses. If it is a type that isn't recognized, then
// it will be passed through.
//
// The structure is modified in-place.
func Hash(salter *salt.Salt, raw interface{}, nonHMACDataKeys []string) error {
	fn := salter.GetIdentifiedHMAC

	switch s := raw.(type) {
	case *logical.Auth:
		if s == nil {
			return nil
		}
		if s.ClientToken != "" {
			s.ClientToken = fn(s.ClientToken)
		}
		if s.Accessor != "" {
			s.Accessor = fn(s.Accessor)
		}

	case *logical.Request:
		if s == nil {
			return nil
		}
		if s.Auth != nil {
			if err := Hash(salter, s.Auth, nil); err != nil {
				return err
			}
		}

		if s.ClientToken != "" {
			s.ClientToken = fn(s.ClientToken)
		}

		if s.ClientTokenAccessor != "" {
			s.ClientTokenAccessor = fn(s.ClientTokenAccessor)
		}

		data, err := HashStructure(s.Data, fn, nonHMACDataKeys)
		if err != nil {
			return err
		}

		s.Data = data.(map[string]interface{})

	case *logical.Response:
		if s == nil {
			return nil
		}

		if s.Auth != nil {
			if err := Hash(salter, s.Auth, nil); err != nil {
				return err
			}
		}

		if s.WrapInfo != nil {
			if err := Hash(salter, s.WrapInfo, nil); err != nil {
				return err
			}
		}

		data, err := HashStructure(s.Data, fn, nonHMACDataKeys)
		if err != nil {
			return err
		}

		s.Data = data.(map[string]interface{})

	case *wrapping.ResponseWrapInfo:
		if s == nil {
			return nil
		}

		s.Token = fn(s.Token)
		s.Accessor = fn(s.Accessor)

		if s.WrappedAccessor != "" {
			s.WrappedAccessor = fn(s.WrappedAccessor)
		}
	}

	return nil
}

// HashStructure takes an interface and hashes all the values within
// the structure. Only _values_ are hashed: keys of objects are not.
//
// For the HashCallback, see the built-in HashCallbacks below.
func HashStructure(s map[string]interface{}, cb HashCallback, ignoredKeys []string) (interface{}, error) {
	scopy, err := copystructure.Copy(s)
	if err != nil {
		return nil, err
	}

	walker := &hashWalker{Callback: cb, IgnoredKeys: ignoredKeys}
	if err := reflectwalk.Walk(scopy, walker); err != nil {
		return nil, err
	}

	return scopy, nil
}

// HashCallback is the callback called for HashStructure to hash
// a value.
type HashCallback func(string) string

// hashWalker implements interfaces for the reflectwalk package
// (github.com/mitchellh/reflectwalk) that can be used to automatically
// replace primitives with a hashed value.
type hashWalker struct {
	// Callback is the function to call with the primitive that is
	// to be hashed. If there is an error, walking will be halted
	// immediately and the error returned.
	Callback HashCallback

	// IgnoreKeys are the keys that wont have the HashCallback applied
	IgnoredKeys []string

	// MapElem appends the key itself (not the reflect.Value) to key.
	// The last element in key is the most recently entered map key.
	// Since Exit pops the last element of key, only nesting to another
	// structure increases the size of this slice.
	key       []string
	lastValue reflect.Value
	// Enter appends to loc and exit pops loc. The last element of loc is thus
	// the current location.
	loc []reflectwalk.Location
	// Map and Slice append to cs, Exit pops the last element off cs.
	// The last element in cs is the most recently entered map or slice.
	cs []reflect.Value
	// MapElem and SliceElem append to csKey. The last element in csKey is the
	// most recently entered map key or slice index. Since Exit pops the last
	// element of csKey, only nesting to another structure increases the size of
	// this slice.
	csKey []reflect.Value
}

// hashTimeType stores a pre-computed reflect.Type for a time.Time so
// we can quickly compare in hashWalker.Struct. We create an empty/invalid
// time.Time{} so we don't need to incur any additional startup cost vs.
// Now() or Unix().
var hashTimeType = reflect.TypeOf(time.Time{})

func (w *hashWalker) Enter(loc reflectwalk.Location) error {
	w.loc = append(w.loc, loc)
	return nil
}

func (w *hashWalker) Exit(loc reflectwalk.Location) error {
	w.loc = w.loc[:len(w.loc)-1]

	switch loc {
	case reflectwalk.Map:
		w.cs = w.cs[:len(w.cs)-1]
	case reflectwalk.MapValue:
		w.key = w.key[:len(w.key)-1]
		w.csKey = w.csKey[:len(w.csKey)-1]
	case reflectwalk.Slice:
		w.cs = w.cs[:len(w.cs)-1]
	case reflectwalk.SliceElem:
		w.csKey = w.csKey[:len(w.csKey)-1]
	}

	return nil
}

func (w *hashWalker) Map(m reflect.Value) error {
	w.cs = append(w.cs, m)
	return nil
}

func (w *hashWalker) MapElem(m, k, v reflect.Value) error {
	w.csKey = append(w.csKey, k)
	w.key = append(w.key, k.String())
	w.lastValue = v
	return nil
}

func (w *hashWalker) Slice(s reflect.Value) error {
	w.cs = append(w.cs, s)
	return nil
}

func (w *hashWalker) SliceElem(i int, elem reflect.Value) error {
	w.csKey = append(w.csKey, reflect.ValueOf(i))
	return nil
}

func (w *hashWalker) Struct(v reflect.Value) error {
	// We are looking for time values. If it isn't one, ignore it.
	if v.Type() != hashTimeType {
		return nil
	}

	if len(w.loc) < 2 {
		// The last element of w.loc is reflectwalk.Struct, by definition.
		// If len(w.loc) == 1 that means hashWalker.Walk was given a struct
		// value and this is the very first step in the walk, and we don't
		// currently support structs as inputs,
		return errors.New("structs as direct inputs not supported")
	}

	// Second to last element of w.loc is location that contains this struct.
	switch w.loc[len(w.loc)-2] {
	case reflectwalk.MapKey:
		return errors.New("time.Time value in a map key cannot be hashed for audits")
	case reflectwalk.MapValue:
		// Create a string value of the time. IMPORTANT: this must never change
		// across Vault versions or the hash value of equivalent time.Time will
		// change.
		strVal := v.Interface().(time.Time).Format(time.RFC3339Nano)

		// Set the map value to the string instead of the time.Time object
		m := w.cs[len(w.cs)-1]
		mk := w.csKey[len(w.cs)-1]
		m.SetMapIndex(mk, reflect.ValueOf(strVal))
	case reflectwalk.SliceElem:
		// Create a string value of the time. IMPORTANT: this must never change
		// across Vault versions or the hash value of equivalent time.Time will
		// change.
		strVal := v.Interface().(time.Time).Format(time.RFC3339Nano)

		// Set the map value to the string instead of the time.Time object
		s := w.cs[len(w.cs)-1]
		si := int(w.csKey[len(w.cs)-1].Int())
		s.Slice(si, si+1).Index(0).Set(reflect.ValueOf(strVal))
	}

	// Skip this entry so that we don't walk the struct.
	return reflectwalk.SkipEntry
}

func (w *hashWalker) StructField(reflect.StructField, reflect.Value) error {
	return nil
}

// Primitive calls Callback to transform strings in-place, except for map keys.
// Strings hiding within interfaces are also transformed.
func (w *hashWalker) Primitive(v reflect.Value) error {
	if w.Callback == nil {
		return nil
	}

	// We don't touch map keys
	if w.loc[len(w.loc)-1] == reflectwalk.MapKey {
		return nil
	}

	setV := v

	// We only care about strings
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.String {
		return nil
	}

	// See if the current key is part of the ignored keys
	currentKey := w.key[len(w.key)-1]
	if strutil.StrListContains(w.IgnoredKeys, currentKey) {
		return nil
	}

	replaceVal := w.Callback(v.String())

	resultVal := reflect.ValueOf(replaceVal)
	switch w.loc[len(w.loc)-1] {
	case reflectwalk.MapValue:
		// If we're in a map, then the only way to set a map value is
		// to set it directly.
		m := w.cs[len(w.cs)-1]
		mk := w.csKey[len(w.cs)-1]
		m.SetMapIndex(mk, resultVal)
	default:
		// Otherwise, we should be addressable
		setV.Set(resultVal)
	}

	return nil
}
