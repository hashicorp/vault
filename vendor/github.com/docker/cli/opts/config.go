package opts

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	swarmtypes "github.com/docker/docker/api/types/swarm"
)

// ConfigOpt is a Value type for parsing configs
type ConfigOpt struct {
	values []*swarmtypes.ConfigReference
}

// Set a new config value
func (o *ConfigOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}

	options := &swarmtypes.ConfigReference{
		File: &swarmtypes.ConfigReferenceFileTarget{
			UID:  "0",
			GID:  "0",
			Mode: 0o444,
		},
	}

	// support a simple syntax of --config foo
	if len(fields) == 1 && !strings.Contains(fields[0], "=") {
		options.File.Name = fields[0]
		options.ConfigName = fields[0]
		o.values = append(o.values, options)
		return nil
	}

	for _, field := range fields {
		key, val, ok := strings.Cut(field, "=")
		if !ok || key == "" {
			return fmt.Errorf("invalid field '%s' must be a key=value pair", field)
		}

		// TODO(thaJeztah): these options should not be case-insensitive.
		switch strings.ToLower(key) {
		case "source", "src":
			options.ConfigName = val
		case "target":
			options.File.Name = val
		case "uid":
			options.File.UID = val
		case "gid":
			options.File.GID = val
		case "mode":
			m, err := strconv.ParseUint(val, 0, 32)
			if err != nil {
				return fmt.Errorf("invalid mode specified: %v", err)
			}

			options.File.Mode = os.FileMode(m)
		default:
			return fmt.Errorf("invalid field in config request: %s", key)
		}
	}

	if options.ConfigName == "" {
		return errors.New("source is required")
	}
	if options.File.Name == "" {
		options.File.Name = options.ConfigName
	}

	o.values = append(o.values, options)
	return nil
}

// Type returns the type of this option
func (o *ConfigOpt) Type() string {
	return "config"
}

// String returns a string repr of this option
func (o *ConfigOpt) String() string {
	configs := []string{}
	for _, config := range o.values {
		repr := fmt.Sprintf("%s -> %s", config.ConfigName, config.File.Name)
		configs = append(configs, repr)
	}
	return strings.Join(configs, ", ")
}

// Value returns the config requests
func (o *ConfigOpt) Value() []*swarmtypes.ConfigReference {
	return o.values
}
