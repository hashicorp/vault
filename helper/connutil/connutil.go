package connutil

import (
	"fmt"
	"strings"
)

// Format defines the desired output format of the connection string
type Format int

const (
	NOTIMPLEMENTED Format = iota
	ADO
)

// BuildConnectionString constructs a connection string in the specified format from a key/value map.
func BuildConnectionString(conf map[string]string, format Format) (string, bool) {
	switch format {
	case ADO:
		var connectionParams []string
		for k, v := range conf {
			connectionParams = append(connectionParams, fmt.Sprintf("%s=%s", k, v))
		}
		return strings.Join(connectionParams, ";"), true
	default:
		return "", false
	}
}

// InjectDefaultsIntoConnectionString adds defaults to a configuration if they are not already
// specified in the configuration and returns a new configuration which includes the defaults
func InjectDefaultsIntoConnectionString(defaults map[string]string, conf map[string]string) map[string]string {
	for k, v := range defaults {
		if _, isSet := conf[k]; !isSet {
			conf[k] = v
		}
	}
	return conf
}

// UpgradeParameterKeysInConnectionString uses a mapping backcomp of HCL connection string properties
// to those understood by the driver to upgrade key names. If an entry in conf uses both a legacy key and a new key
// then the final key in conf will be the new key (the legacy key will be deleted) and the value associated with this
// key will be the legacy value (ie. adding a new key won't clobber an existing working configuration).
// UpgradeParameterKeysInConnectionString returns a new configuration map
func UpgradeParameterKeysInConnectionString(backcomp map[string]string, conf map[string]string) map[string]string {
	for legacyKey, newKey := range backcomp {
		value, isSet := conf[legacyKey]
		if isSet {
			conf[newKey] = value
			delete(conf, legacyKey)
		}
	}
	return conf
}
