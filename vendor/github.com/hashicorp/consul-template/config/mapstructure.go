package config

import (
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

// StringToFileModeFunc returns a function that converts strings to os.FileMode
// value. This is designed to be used with mapstructure for parsing out a
// filemode value.
func StringToFileModeFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(os.FileMode(0)) {
			return data, nil
		}

		// Convert it by parsing
		v, err := strconv.ParseUint(data.(string), 8, 12)
		if err != nil {
			return data, err
		}
		return os.FileMode(v), nil
	}
}

// StringToWaitDurationHookFunc returns a function that converts strings to wait
// value. This is designed to be used with mapstructure for parsing out a wait
// value.
func StringToWaitDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(WaitConfig{}) {
			return data, nil
		}

		// Convert it by parsing
		return ParseWaitConfig(data.(string))
	}
}

// ConsulStringToStructFunc checks if the value set for the key should actually
// be a struct and sets the appropriate value in the struct. This is for
// backwards-compatability with older versions of Consul Template.
func ConsulStringToStructFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(ConsulConfig{}) && f.Kind() == reflect.String {
			log.Println("[WARN] consul now accepts a stanza instead of a string. " +
				"Update your configuration files and change consul = \"\" to " +
				"consul { } instead.")
			return &ConsulConfig{
				Address: String(data.(string)),
			}, nil
		}

		return data, nil
	}
}
