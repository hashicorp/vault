package cert

import (
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCRLs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "crls/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the certificate",
			},

			"crl": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The public certificate that should be trusted.
May be DER or PEM encoded. Note: the expiration time
is ignored; if the CRL is no longer valid, delete it
using the same name as specified here.`,
			},

			"serial": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `If specified, for a read, information for this
serial will be returned rather than the named CRL.
CRL. For a delete, only this serial will be removed
from the named CRL entry. This can be a hex-formatted
string separated by : or -, or an integer string;
this will be assumed to be base 10 unless prefixed
by "0x" for base 16 or "0" for base 8.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathCRLDelete,
			logical.ReadOperation:   b.pathCRLRead,
			logical.WriteOperation:  b.pathCRLWrite,
		},

		HelpSynopsis:    pathCRLsHelpSyn,
		HelpDescription: pathCRLsHelpDesc,
	}
}

func parseSerialString(input string) (*big.Int, error) {
	ret := &big.Int{}

	switch {
	case strings.Count(input, ":") > 0:
		serialBytes := certutil.ParseHexFormatted(input, ":")
		if serialBytes == nil {
			return nil, fmt.Errorf("error parsing serial %s", input)
		}
		ret.SetBytes(serialBytes)
	case strings.Count(input, "-") > 0:
		serialBytes := certutil.ParseHexFormatted(input, "-")
		if serialBytes == nil {
			return nil, fmt.Errorf("error parsing serial %s", input)
		}
		ret.SetBytes(serialBytes)
	default:
		var success bool
		ret, success = ret.SetString(input, 0)
		if !success {
			return nil, fmt.Errorf("error parsing serial %s", input)
		}
	}

	return ret, nil
}

func (b *backend) pathCRLDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(`"name" parameter cannot be empty`), nil
	}

	serialStr := d.Get("serial").(string)

	if serialStr != "" {
		serial, err := parseSerialString(serialStr)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		err = b.deleteSerial(req.Storage, name, serial)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"error deleting serial %s from CRL %s: %s", serial, name, err),
			), nil
		}
		return nil, nil
	}

	// deleteIndex is best effort to ensure it removes as much as possible if there is
	// a problem, so it does not currently return an error
	b.deleteIndex(req.Storage, name)

	err := req.Storage.Delete("crls/set/" + name)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"error deleting set %s: %s", name, err),
		), nil
	}

	return nil, nil
}

func (b *backend) pathCRLRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(`"name" parameter cannot be empty`), nil
	}

	return nil, nil
	/*
		return &logical.Response{
			Data: map[string]interface{}{
				"certificate":  cert.Certificate,
				"display_name": cert.DisplayName,
				"policies":     strings.Join(cert.Policies, ","),
				"ttl":          duration / time.Second,
			},
		}, nil
	*/
}

func (b *backend) pathCRLWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(`"name" parameter cannot be empty`), nil
	}
	crl := d.Get("crl").(string)

	certList, err := x509.ParseCRL([]byte(crl))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse CRL: %v", err)), nil
	}
	if certList == nil {
		return logical.ErrorResponse("parsed CRL is nil"), nil
	}

	entry, err := logical.StorageEntryJSON("crls/set/"+name, crl)
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(entry); err != nil {
		return nil, err
	}

	// Clear out old entries if this is replacing a previously-set CRL.
	// In practice this is what the index is for; it lets us store
	// the certs individually for storage efficiency but ensure we can
	// clean up properly. So use it to clean up.
	// N.B.: This is best-effort. The worst thing that can happen is
	// some wasted storage.
	b.deleteIndex(req.Storage, name)

	crlSetIndex := []*big.Int{}
	for _, revokedCert := range certList.TBSCertList.RevokedCertificates {
		crlSetIndex = append(crlSetIndex, revokedCert.SerialNumber)
	}

	entry, err = logical.StorageEntryJSON("crls/index/"+name, crlSetIndex)
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(entry); err != nil {
		return nil, err
	}

	for _, revokedSerial := range crlSetIndex {
		entry, err = logical.StorageEntryJSON("crls/serial/"+revokedSerial.String(),
			&RevokedSerialInfo{
				name: nil,
			})
		if err != nil {
			return nil, err
		}
		if err = req.Storage.Put(entry); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (b *backend) deleteIndex(storage logical.Storage, name string) {
	entry, err := storage.Get("crls/index/" + name)
	if err != nil {
		return
	}
	if entry == nil {
		return
	}

	crlSetIndex := []*big.Int{}
	var revokedSerialInfo RevokedSerialInfo

	err = entry.DecodeJSON(&crlSetIndex)
	if err != nil {
		goto destroyIndex
	}

	for _, serial := range crlSetIndex {
		entry, err := storage.Get("crls/serial/" + serial.String())
		if err != nil {
			goto deleteSerial
		}
		err = entry.DecodeJSON(&revokedSerialInfo)
		// In theory we could now delete the serial when it still exists in a CRL set,
		// but if we can't decode it something is wrong with the entry anyways, and
		// we're not going to be able to decode it when checking revocation either.
		if err != nil {
			goto deleteSerial
		}

		delete(revokedSerialInfo, name)
		if len(revokedSerialInfo) > 0 {
			entry, err = logical.StorageEntryJSON("crls/serial/"+serial.String(), revokedSerialInfo)
			if err != nil {
				continue
			}
			storage.Put(entry)
			continue
		}
	deleteSerial:
		storage.Delete("crls/serial/" + serial.String())
	}

destroyIndex:
	storage.Delete("crls/index/" + name)

	return
}

func (b *backend) deleteSerial(storage logical.Storage, name string, serial *big.Int) error {
	var revokedSerialInfo RevokedSerialInfo
	entry, err := storage.Get("crls/serial/" + serial.String())
	if err != nil {
		return fmt.Errorf("error retrieving entry for serial %s: %s", serial, err)
	}
	if entry == nil {
		return nil
	}

	err = entry.DecodeJSON(&revokedSerialInfo)
	// In theory we could now delete the serial when it still exists in a CRL set,
	// but if we can't decode it something is wrong with the entry anyways, and
	// we're not going to be able to decode it when checking revocation either.
	if err != nil {
		goto deleteSerial
	}

	delete(revokedSerialInfo, name)

	if len(revokedSerialInfo) > 0 {
		entry, err := logical.StorageEntryJSON("crls/serial/"+serial.String(), revokedSerialInfo)
		if err != nil {
			return fmt.Errorf("error creating updated storage entry for serial %s: %s", serial, err)
		}
		err = storage.Put(entry)
		if err != nil {
			return fmt.Errorf("error storing updated entry for serial %s: %s", serial, err)
		}

		return nil

	}

deleteSerial:
	err = storage.Delete("crls/serial/" + serial.String())
	if err != nil {
		return fmt.Errorf("error deleting serial entry for serial %s: %s", serial, err)
	}

	return nil
}

type RevokedSerialInfo map[string]interface{}

//FIXME
const pathCRLsHelpSyn = `
Manage trusted certificates used for authentication.
`

const pathCRLsHelpDesc = `
This endpoint allows you to create, read, update, and delete trusted certificates
that are allowed to authenticate.

Deleting a certificate will not revoke auth for prior authenticated connections.
To do this, do a revoke on "login". If you don't need to revoke login immediately,
then the next renew will cause the lease to expire.
`
