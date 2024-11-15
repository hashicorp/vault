package opts

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/go-units"
	"github.com/sirupsen/logrus"
)

// MountOpt is a Value type for parsing mounts
type MountOpt struct {
	values []mounttypes.Mount
}

// Set a new mount value
//
//nolint:gocyclo
func (m *MountOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}

	mount := mounttypes.Mount{}

	volumeOptions := func() *mounttypes.VolumeOptions {
		if mount.VolumeOptions == nil {
			mount.VolumeOptions = &mounttypes.VolumeOptions{
				Labels: make(map[string]string),
			}
		}
		if mount.VolumeOptions.DriverConfig == nil {
			mount.VolumeOptions.DriverConfig = &mounttypes.Driver{}
		}
		return mount.VolumeOptions
	}

	bindOptions := func() *mounttypes.BindOptions {
		if mount.BindOptions == nil {
			mount.BindOptions = new(mounttypes.BindOptions)
		}
		return mount.BindOptions
	}

	tmpfsOptions := func() *mounttypes.TmpfsOptions {
		if mount.TmpfsOptions == nil {
			mount.TmpfsOptions = new(mounttypes.TmpfsOptions)
		}
		return mount.TmpfsOptions
	}

	setValueOnMap := func(target map[string]string, value string) {
		k, v, _ := strings.Cut(value, "=")
		if k != "" {
			target[k] = v
		}
	}

	mount.Type = mounttypes.TypeVolume // default to volume mounts
	// Set writable as the default
	for _, field := range fields {
		key, val, ok := strings.Cut(field, "=")

		// TODO(thaJeztah): these options should not be case-insensitive.
		key = strings.ToLower(key)

		if !ok {
			switch key {
			case "readonly", "ro":
				mount.ReadOnly = true
				continue
			case "volume-nocopy":
				volumeOptions().NoCopy = true
				continue
			case "bind-nonrecursive":
				bindOptions().NonRecursive = true
				continue
			default:
				return fmt.Errorf("invalid field '%s' must be a key=value pair", field)
			}
		}

		switch key {
		case "type":
			mount.Type = mounttypes.Type(strings.ToLower(val))
		case "source", "src":
			mount.Source = val
			if strings.HasPrefix(val, "."+string(filepath.Separator)) || val == "." {
				if abs, err := filepath.Abs(val); err == nil {
					mount.Source = abs
				}
			}
		case "target", "dst", "destination":
			mount.Target = val
		case "readonly", "ro":
			mount.ReadOnly, err = strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, val)
			}
		case "consistency":
			mount.Consistency = mounttypes.Consistency(strings.ToLower(val))
		case "bind-propagation":
			bindOptions().Propagation = mounttypes.Propagation(strings.ToLower(val))
		case "bind-nonrecursive":
			bindOptions().NonRecursive, err = strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, val)
			}
			logrus.Warn("bind-nonrecursive is deprecated, use bind-recursive=disabled instead")
		case "bind-recursive":
			switch val {
			case "enabled": // read-only mounts are recursively read-only if Engine >= v25 && kernel >= v5.12, otherwise writable
				// NOP
			case "disabled": // alias of bind-nonrecursive=true
				bindOptions().NonRecursive = true
			case "writable": // conforms to the default read-only bind-mount of Docker v24; read-only mounts are recursively mounted but not recursively read-only
				bindOptions().ReadOnlyNonRecursive = true
			case "readonly": // force recursively read-only, or raise an error
				bindOptions().ReadOnlyForceRecursive = true
				// TODO: implicitly set propagation and error if the user specifies a propagation in a future refactor/UX polish pass
				// https://github.com/docker/cli/pull/4316#discussion_r1341974730
			default:
				return fmt.Errorf("invalid value for %s: %s (must be \"enabled\", \"disabled\", \"writable\", or \"readonly\")",
					key, val)
			}
		case "volume-subpath":
			volumeOptions().Subpath = val
		case "volume-nocopy":
			volumeOptions().NoCopy, err = strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid value for volume-nocopy: %s", val)
			}
		case "volume-label":
			setValueOnMap(volumeOptions().Labels, val)
		case "volume-driver":
			volumeOptions().DriverConfig.Name = val
		case "volume-opt":
			if volumeOptions().DriverConfig.Options == nil {
				volumeOptions().DriverConfig.Options = make(map[string]string)
			}
			setValueOnMap(volumeOptions().DriverConfig.Options, val)
		case "tmpfs-size":
			sizeBytes, err := units.RAMInBytes(val)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, val)
			}
			tmpfsOptions().SizeBytes = sizeBytes
		case "tmpfs-mode":
			ui64, err := strconv.ParseUint(val, 8, 32)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, val)
			}
			tmpfsOptions().Mode = os.FileMode(ui64)
		default:
			return fmt.Errorf("unexpected key '%s' in '%s'", key, field)
		}
	}

	if mount.Type == "" {
		return errors.New("type is required")
	}

	if mount.Target == "" {
		return errors.New("target is required")
	}

	if mount.VolumeOptions != nil && mount.Type != mounttypes.TypeVolume {
		return fmt.Errorf("cannot mix 'volume-*' options with mount type '%s'", mount.Type)
	}
	if mount.BindOptions != nil && mount.Type != mounttypes.TypeBind {
		return fmt.Errorf("cannot mix 'bind-*' options with mount type '%s'", mount.Type)
	}
	if mount.TmpfsOptions != nil && mount.Type != mounttypes.TypeTmpfs {
		return fmt.Errorf("cannot mix 'tmpfs-*' options with mount type '%s'", mount.Type)
	}

	if mount.BindOptions != nil {
		if mount.BindOptions.ReadOnlyNonRecursive {
			if !mount.ReadOnly {
				return errors.New("option 'bind-recursive=writable' requires 'readonly' to be specified in conjunction")
			}
		}
		if mount.BindOptions.ReadOnlyForceRecursive {
			if !mount.ReadOnly {
				return errors.New("option 'bind-recursive=readonly' requires 'readonly' to be specified in conjunction")
			}
			if mount.BindOptions.Propagation != mounttypes.PropagationRPrivate {
				return errors.New("option 'bind-recursive=readonly' requires 'bind-propagation=rprivate' to be specified in conjunction")
			}
		}
	}

	m.values = append(m.values, mount)
	return nil
}

// Type returns the type of this option
func (m *MountOpt) Type() string {
	return "mount"
}

// String returns a string repr of this option
func (m *MountOpt) String() string {
	mounts := []string{}
	for _, mount := range m.values {
		repr := fmt.Sprintf("%s %s %s", mount.Type, mount.Source, mount.Target)
		mounts = append(mounts, repr)
	}
	return strings.Join(mounts, ", ")
}

// Value returns the mounts
func (m *MountOpt) Value() []mounttypes.Mount {
	return m.values
}
