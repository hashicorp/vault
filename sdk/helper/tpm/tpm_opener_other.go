// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package tpm

import (
	"fmt"
	"io"
	"os"

	"github.com/google/go-attestation/attest"
	"github.com/google/go-tpm/legacy/tpm2"
)

// openTPM opens a TPM device at the given path.
// If devicePath is empty, uses the default TPM device.
// If devicePath is a socket, dials it via unix socket.
// If devicePath is a regular file, opens it as a character device.
// Returns a tpmHandle that closes both the TPM and the underlying connection.
func openTPM(devicePath string) (*attest.TPM, error) {
	if devicePath == "" {
		// No device path specified, use default TPM
		atpm, err := attest.OpenTPM(&attest.OpenConfig{})
		if err != nil {
			return nil, err
		}
		return atpm, nil
	}

	tpm, err := tpm2.OpenTPM(devicePath)
	if err != nil {
		return nil, err
	}

	atpm, err := attest.OpenTPM(&attest.OpenConfig{
		CommandChannel: &linuxCmdChannel{tpm},
	})
	if err != nil {
		tpm.Close()
		return nil, fmt.Errorf("failed to open TPM via device file: %w", err)
	}

	return atpm, nil
}

type linuxCmdChannel struct {
	io.ReadWriteCloser
}

// MeasurementLog implements CommandChannelTPM20.
func (cc *linuxCmdChannel) MeasurementLog() ([]byte, error) {
	return os.ReadFile("/sys/kernel/security/tpm0/binary_bios_measurements")
}
