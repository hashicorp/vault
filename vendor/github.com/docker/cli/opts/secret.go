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

// SecretOpt is a Value type for parsing secrets
type SecretOpt struct {
	values []*swarmtypes.SecretReference
}

// Set a new secret value
func (o *SecretOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}

	options := &swarmtypes.SecretReference{
		File: &swarmtypes.SecretReferenceFileTarget{
			UID:  "0",
			GID:  "0",
			Mode: 0o444,
		},
	}

	// support a simple syntax of --secret foo
	if len(fields) == 1 && !strings.Contains(fields[0], "=") {
		options.File.Name = fields[0]
		options.SecretName = fields[0]
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
			options.SecretName = val
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
			return errors.New("invalid field in secret request: " + key)
		}
	}

	if options.SecretName == "" {
		return errors.New("source is required")
	}
	if options.File.Name == "" {
		options.File.Name = options.SecretName
	}

	o.values = append(o.values, options)
	return nil
}

// Type returns the type of this option
func (o *SecretOpt) Type() string {
	return "secret"
}

// String returns a string repr of this option
func (o *SecretOpt) String() string {
	secrets := []string{}
	for _, secret := range o.values {
		repr := fmt.Sprintf("%s -> %s", secret.SecretName, secret.File.Name)
		secrets = append(secrets, repr)
	}
	return strings.Join(secrets, ", ")
}

// Value returns the secret requests
func (o *SecretOpt) Value() []*swarmtypes.SecretReference {
	return o.values
}
