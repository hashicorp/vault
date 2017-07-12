package kvbuilder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/mitchellh/mapstructure"
)

// Builder is a struct to build a key/value mapping based on a list
// of "k=v" pairs, where the value might come from stdin, a file, etc.
type Builder struct {
	Stdin io.Reader

	result map[string]interface{}
	stdin  bool
}

// Map returns the built map.
func (b *Builder) Map() map[string]interface{} {
	return b.result
}

// Add adds to the mapping with the given args.
func (b *Builder) Add(args ...string) error {
	for _, a := range args {
		if err := b.add(a); err != nil {
			return fmt.Errorf("Invalid key/value pair '%s': %s", a, err)
		}
	}

	return nil
}

func (b *Builder) add(raw string) error {
	// Regardless of validity, make sure we make our result
	if b.result == nil {
		b.result = make(map[string]interface{})
	}

	// Empty strings are fine, just ignored
	if raw == "" {
		return nil
	}

	// Split into key/value
	parts := strings.SplitN(raw, "=", 2)

	// If the arg is exactly "-", then we need to read from stdin
	// and merge the results into the resulting structure.
	if len(parts) == 1 {
		if raw == "-" {
			if b.Stdin == nil {
				return fmt.Errorf("stdin is not supported")
			}
			if b.stdin {
				return fmt.Errorf("stdin already consumed")
			}

			b.stdin = true
			return b.addReader(b.Stdin)
		}

		// If the arg begins with "@" then we need to read a file directly
		if raw[0] == '@' {
			f, err := os.Open(raw[1:])
			if err != nil {
				return err
			}
			defer f.Close()

			return b.addReader(f)
		}
	}

	if len(parts) != 2 {
		return fmt.Errorf("format must be key=value")
	}
	key, value := parts[0], parts[1]

	if len(value) > 0 {
		if value[0] == '@' {
			contents, err := ioutil.ReadFile(value[1:])
			if err != nil {
				return fmt.Errorf("error reading file: %s", err)
			}

			value = string(contents)
		} else if value[0] == '\\' && value[1] == '@' {
			value = value[1:]
		} else if value == "-" {
			if b.Stdin == nil {
				return fmt.Errorf("stdin is not supported")
			}
			if b.stdin {
				return fmt.Errorf("stdin already consumed")
			}
			b.stdin = true

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, b.Stdin); err != nil {
				return err
			}

			value = buf.String()
		}
	}

	// Repeated keys will be converted into a slice
	if existingValue, ok := b.result[key]; ok {
		var sliceValue []interface{}
		if err := mapstructure.WeakDecode(existingValue, &sliceValue); err != nil {
			return err
		}
		sliceValue = append(sliceValue, value)
		b.result[key] = sliceValue
		return nil
	}

	b.result[key] = value
	return nil
}

func (b *Builder) addReader(r io.Reader) error {
	return jsonutil.DecodeJSONFromReader(r, &b.result)
}
