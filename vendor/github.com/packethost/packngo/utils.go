package packngo

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
)

var (
	timestampType = reflect.TypeOf(Timestamp{})

	// Facilities DEPRECATED Use Facilities.List
	Facilities = []string{
		"yyz1", "nrt1", "atl1", "mrs1", "hkg1", "ams1",
		"ewr1", "sin1", "dfw1", "lax1", "syd1", "sjc1",
		"ord1", "iad1", "fra1", "sea1", "dfw2"}

	// FacilityFeatures DEPRECATED Use Facilities.List
	FacilityFeatures = []string{
		"baremetal", "layer_2", "backend_transfer", "storage", "global_ipv4"}

	// UtilizationLevels DEPRECATED
	UtilizationLevels = []string{"unavailable", "critical", "limited", "normal"}

	// DevicePlans DEPRECATED Use Plans.List
	DevicePlans = []string{"c2.medium.x86", "g2.large.x86",
		"m2.xlarge.x86", "x2.xlarge.x86", "baremetal_2a", "baremetal_2a2",
		"baremetal_1", "baremetal_3", "baremetal_2", "baremetal_s",
		"baremetal_0", "baremetal_1e",
	}
)

// Stringify creates a string representation of the provided message
// DEPRECATED This is used internally and should not be exported by packngo
func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	// TODO(displague) errors here are not reported
	_ = stringifyValue(&buf, v)
	return buf.String()
}

// StreamToString converts a reader to a string
// DEPRECATED This is unused and should not be exported by packngo
func StreamToString(stream io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(stream); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// stringifyValue was graciously cargoculted from the goprotubuf library
func stringifyValue(w io.Writer, val reflect.Value) error {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		_, err := w.Write([]byte("<nil>"))
		return err
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		if _, err := fmt.Fprintf(w, `"%s"`, v); err != nil {
			return err
		}
	case reflect.Slice:
		if _, err := w.Write([]byte{'['}); err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				if _, err := w.Write([]byte{' '}); err != nil {
					return err
				}
			}

			if err := stringifyValue(w, v.Index(i)); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte{']'}); err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		if v.Type().Name() != "" {
			if _, err := w.Write([]byte(v.Type().String())); err != nil {
				return err
			}
		}

		// special handling of Timestamp values
		if v.Type() == timestampType {
			_, err := fmt.Fprintf(w, "{%s}", v.Interface())
			return err
		}

		if _, err := w.Write([]byte{'{'}); err != nil {
			return err
		}

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				if _, err := w.Write([]byte(", ")); err != nil {
					return err
				}
			} else {
				sep = true
			}

			if _, err := w.Write([]byte(v.Type().Field(i).Name)); err != nil {
				return err
			}
			if _, err := w.Write([]byte{':'}); err != nil {
				return err
			}

			if err := stringifyValue(w, fv); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte{'}'}); err != nil {
			return err
		}
	default:
		if v.CanInterface() {
			if _, err := fmt.Fprint(w, v.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

// validate UUID
func ValidateUUID(uuid string) error {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	if !r.MatchString(uuid) {
		return fmt.Errorf("%s is not a valid UUID", uuid)
	}
	return nil
}
