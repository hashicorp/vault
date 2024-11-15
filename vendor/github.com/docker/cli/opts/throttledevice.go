package opts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/blkiodev"
	"github.com/docker/go-units"
)

// ValidatorThrottleFctType defines a validator function that returns a validated struct and/or an error.
type ValidatorThrottleFctType func(val string) (*blkiodev.ThrottleDevice, error)

// ValidateThrottleBpsDevice validates that the specified string has a valid device-rate format.
func ValidateThrottleBpsDevice(val string) (*blkiodev.ThrottleDevice, error) {
	k, v, ok := strings.Cut(val, ":")
	if !ok || k == "" {
		return nil, fmt.Errorf("bad format: %s", val)
	}
	// TODO(thaJeztah): should we really validate this on the client?
	if !strings.HasPrefix(k, "/dev/") {
		return nil, fmt.Errorf("bad format for device path: %s", val)
	}
	rate, err := units.RAMInBytes(v)
	if err != nil {
		return nil, fmt.Errorf("invalid rate for device: %s. The correct format is <device-path>:<number>[<unit>]. Number must be a positive integer. Unit is optional and can be kb, mb, or gb", val)
	}
	if rate < 0 {
		return nil, fmt.Errorf("invalid rate for device: %s. The correct format is <device-path>:<number>[<unit>]. Number must be a positive integer. Unit is optional and can be kb, mb, or gb", val)
	}

	return &blkiodev.ThrottleDevice{
		Path: k,
		Rate: uint64(rate),
	}, nil
}

// ValidateThrottleIOpsDevice validates that the specified string has a valid device-rate format.
func ValidateThrottleIOpsDevice(val string) (*blkiodev.ThrottleDevice, error) {
	k, v, ok := strings.Cut(val, ":")
	if !ok || k == "" {
		return nil, fmt.Errorf("bad format: %s", val)
	}
	// TODO(thaJeztah): should we really validate this on the client?
	if !strings.HasPrefix(k, "/dev/") {
		return nil, fmt.Errorf("bad format for device path: %s", val)
	}
	rate, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid rate for device: %s. The correct format is <device-path>:<number>. Number must be a positive integer", val)
	}

	return &blkiodev.ThrottleDevice{Path: k, Rate: rate}, nil
}

// ThrottledeviceOpt defines a map of ThrottleDevices
type ThrottledeviceOpt struct {
	values    []*blkiodev.ThrottleDevice
	validator ValidatorThrottleFctType
}

// NewThrottledeviceOpt creates a new ThrottledeviceOpt
func NewThrottledeviceOpt(validator ValidatorThrottleFctType) ThrottledeviceOpt {
	return ThrottledeviceOpt{
		values:    []*blkiodev.ThrottleDevice{},
		validator: validator,
	}
}

// Set validates a ThrottleDevice and sets its name as a key in ThrottledeviceOpt
func (opt *ThrottledeviceOpt) Set(val string) error {
	var value *blkiodev.ThrottleDevice
	if opt.validator != nil {
		v, err := opt.validator(val)
		if err != nil {
			return err
		}
		value = v
	}
	opt.values = append(opt.values, value)
	return nil
}

// String returns ThrottledeviceOpt values as a string.
func (opt *ThrottledeviceOpt) String() string {
	out := make([]string, 0, len(opt.values))
	for _, v := range opt.values {
		out = append(out, v.String())
	}

	return fmt.Sprintf("%v", out)
}

// GetList returns a slice of pointers to ThrottleDevices.
func (opt *ThrottledeviceOpt) GetList() []*blkiodev.ThrottleDevice {
	out := make([]*blkiodev.ThrottleDevice, len(opt.values))
	copy(out, opt.values)
	return out
}

// Type returns the option type
func (opt *ThrottledeviceOpt) Type() string {
	return "list"
}
