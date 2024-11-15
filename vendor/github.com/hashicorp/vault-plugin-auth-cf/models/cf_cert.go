// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strings"
)

// NewCFCertificateFromx509 converts a x509 certificate to a valid, well-formed CF certificate,
// erroring if this isn't possible.
func NewCFCertificateFromx509(certificate *x509.Certificate) (*CFCertificate, error) {
	if len(certificate.IPAddresses) != 1 {
		return nil, fmt.Errorf("valid CF certs have one IP address, but this has %s", certificate.IPAddresses)
	}

	cfCert := &CFCertificate{
		InstanceID: certificate.Subject.CommonName,
		IPAddress:  certificate.IPAddresses[0].String(),
	}

	spaces := 0
	orgs := 0
	apps := 0
	for _, ou := range certificate.Subject.OrganizationalUnit {
		if strings.HasPrefix(ou, "space:") {
			cfCert.SpaceID = strings.Split(ou, "space:")[1]
			spaces++
			continue
		}
		if strings.HasPrefix(ou, "organization:") {
			cfCert.OrgID = strings.Split(ou, "organization:")[1]
			orgs++
			continue
		}
		if strings.HasPrefix(ou, "app:") {
			cfCert.AppID = strings.Split(ou, "app:")[1]
			apps++
			continue
		}
	}
	if spaces > 1 {
		return nil, fmt.Errorf("expected 1 space but received %d", spaces)
	}
	if orgs > 1 {
		return nil, fmt.Errorf("expected 1 org but received %d", orgs)
	}
	if apps > 1 {
		return nil, fmt.Errorf("expected 1 app but received %d", apps)
	}
	if err := cfCert.validate(); err != nil {
		return nil, err
	}
	return cfCert, nil
}

// NewCFCertificateFromx509 converts the given fields to a valid, well-formed CF certificate,
// erroring if this isn't possible.
func NewCFCertificate(instanceID, orgID, spaceID, appID, ipAddress string) (*CFCertificate, error) {
	cfCert := &CFCertificate{
		InstanceID: instanceID,
		OrgID:      orgID,
		SpaceID:    spaceID,
		AppID:      appID,
		IPAddress:  ipAddress,
	}
	if err := cfCert.validate(); err != nil {
		return nil, err
	}
	return cfCert, nil
}

// CFCertificate isn't intended to be instantiated directly; but rather through one of the New
// methods, which contain logic validating that the expected fields exist.
type CFCertificate struct {
	InstanceID, OrgID, SpaceID, AppID, IPAddress string
}

func (c *CFCertificate) validate() error {
	if c.InstanceID == "" {
		return errors.New("no instance ID on given certificate")
	}
	if c.AppID == "" {
		return errors.New("no app ID on given certificate")
	}
	if c.OrgID == "" {
		return errors.New("no org ID on given certificate")
	}
	if c.SpaceID == "" {
		return errors.New("no space ID on given certificate")
	}
	if c.IPAddress == "" {
		return errors.New("ip address is unspecified")
	}
	if net.ParseIP(c.IPAddress) == nil {
		return fmt.Errorf("%q could not be parsed as a valid IP address", c.IPAddress)
	}
	return nil
}
