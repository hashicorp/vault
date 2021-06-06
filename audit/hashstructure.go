package audit

import (
	"encoding/json"
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

// HashAuth returns a hashed copy of the logical.Auth input.
func HashAuth(salter *salt.Salt, in *logical.Auth, HMACAccessor bool) (*logical.Auth, error) {
	if in == nil {
		return nil, nil
	}

	fn := salter.GetIdentifiedHMAC
	auth := *in

	if auth.ClientToken != "" {
		auth.ClientToken = fn(auth.ClientToken)
	}
	if HMACAccessor && auth.Accessor != "" {
		auth.Accessor = fn(auth.Accessor)
	}
	return &auth, nil
}

// HashRequest returns a hashed copy of the logical.Request input.
func HashRequest(salter *salt.Salt, in *logical.Request, HMACAccessor bool, nonHMACDataKeys []string) (*logical.Request, error) {
	if in == nil {
		return nil, nil
	}

	fn := salter.GetIdentifiedHMAC
	req := *in

	if req.Auth != nil {
		cp, err := copystructure.Copy(req.Auth)
		if err != nil {
			return nil, err
		}

		req.Auth, err = HashAuth(salter, cp.(*logical.Auth), HMACAccessor)
		if err != nil {
			return nil, err
		}
	}

	if req.ClientToken != "" {
		req.ClientToken = fn(req.ClientToken)
	}
	if HMACAccessor && req.ClientTokenAccessor != "" {
		req.ClientTokenAccessor = fn(req.ClientTokenAccessor)
	}

	if req.Data != nil {
		copy, err := copystructure.Copy(req.Data)
		if err != nil {
			return nil, err
		}

		err = hashMap(fn, copy.(map[string]interface{}), nonHMACDataKeys)
		if err != nil {
			return nil, err
		}
		req.Data = copy.(map[string]interface{})
	}

	return &req, nil
}

func hashMap(fn func(string) string, data map[string]interface{}, nonHMACDataKeys []string) error {
	for k, v := range data {
		if o, ok := v.(logical.OptMarshaler); ok {
			marshaled, err := o.MarshalJSONWithOptions(&logical.MarshalOptions{
				ValueHasher: fn,
			})
			if err != nil {
				return err
			}
			data[k] = json.RawMessage(marshaled)
		}
	}

	return HashStructure(data, fn, nonHMACDataKeys)
}

// HashResponse returns a hashed copy of the logical.Request input.
func HashResponse(salter *salt.Salt, in *logical.Response, HMACAccessor bool, nonHMACDataKeys []string) (*logical.Response, error) {
	if in == nil {
		return nil, nil
	}

	fn := salter.GetIdentifiedHMAC
	resp := *in

	if resp.Auth != nil {
		cp, err := copystructure.Copy(resp.Auth)
		if err != nil {
			return nil, err
		}

		resp.Auth, err = HashAuth(salter, cp.(*logical.Auth), HMACAccessor)
		if err != nil {
			return nil, err
		}
	}

	if resp.Data != nil {
		copy, err := copystructure.Copy(resp.Data)
		if err != nil {
			return nil, err
		}

		mapCopy := copy.(map[string]interface{})
		if b, ok := mapCopy[logical.HTTPRawBody].([]byte); ok {
			mapCopy[logical.HTTPRawBody] = string(b)
		}

		err = hashMap(fn, mapCopy, nonHMACDataKeys)
		if err != nil {
			return nil, err
		}
		resp.Data = mapCopy
	}
	if resp.WrapInfo != nil {
		var err error
		resp.WrapInfo, err = HashWrapInfo(salter, resp.WrapInfo, HMACAccessor)
		if err != nil {
			return nil, err
		}
	}

	return &resp, nil
}

// HashWrapInfo returns a hashed copy of the wrapping.ResponseWrapInfo input.
func HashWrapInfo(salter *salt.Salt, in *wrapping.ResponseWrapInfo, HMACAccessor bool) (*wrapping.ResponseWrapInfo, error) {
	if in == nil {
		return nil, nil
	}

	fn := salter.GetIdentifiedHMAC
	wrapinfo := *in

	wrapinfo.Token = fn(wrapinfo.Token)

	if HMACAccessor {
		wrapinfo.Accessor = fn(wrapinfo.Accessor)

		if wrapinfo.WrappedAccessor != "" {
			wrapinfo.WrappedAccessor = fn(wrapinfo.WrappedAccessor)
		}
	}

	return &wrapinfo, nil
}

// HashStructure takes an interface and hashes all the values within
// the structure. Only _values_ are hashed: keys of objects are not.
//
// For the HashCallback, see the built-in HashCallbacks below.
func HashStructure(s interface{}, cb HashCallback, ignoredKeys []string) error {
	walker := &hashWalker{Callback: cb, IgnoredKeys: ignoredKeys}
	return reflectwalk.Walk(s, walker)
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

	if len(w.loc) < 3 {
		// The last element of w.loc is reflectwalk.Struct, by definition.
		// If len(w.loc) < 3 that means hashWalker.Walk was given a struct
		// value and this is the very first step in the walk, and we don't
		// currently support structs as inputs,
		return errors.New("structs as direct inputs not supported")
	}

	// Second to last element of w.loc is location that contains this struct.
	switch w.loc[len(w.loc)-2] {
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
	case reflectwalk.SliceElem:
		s := w.cs[len(w.cs)-1]
		si := int(w.csKey[len(w.cs)-1].Int())
		s.Slice(si, si+1).Index(0).Set(resultVal)
	default:
		// Otherwise, we should be addressable
		setV.Set(resultVal)
	}

	return nil
}
