package cert

import (
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCRLSets(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "crlsets/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the certificate",
			},

			"crl": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The public certificate that should be trusted.
Must be x509 PEM encoded.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathCRLSetDelete,
			logical.ReadOperation:   b.pathCRLSetRead,
			logical.WriteOperation:  b.pathCRLSetWrite,
		},

		HelpSynopsis:    pathCRLSetsHelpSyn,
		HelpDescription: pathCRLSetsHelpDesc,
	}
}

func (b *backend) pathCRLSetDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// FIXME: There should be a path, or a value to this path, to clear out
	// an individual serial in case of an index out-of-sync issue
	err := req.Storage.Delete("cert/" + strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathCRLSetRead(
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

func (b *backend) pathCRLSetWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	crl := d.Get("crl").(string)

	certList, err := x509.ParseCRL([]byte(crl))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse CRL: %v", err)), nil
	}
	if certList == nil {
		return logical.ErrorResponse("parsed CRL is nil"), nil
	}

	entry, err := logical.StorageEntryJSON("crlsets/set/"+name, crl)
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
	// some wasted storage
	b.cleanIndex(req.Storage, name)

	crlSetIndex := []*big.Int{}
	for _, revokedCert := range certList.TBSCertList.RevokedCertificates {
		crlSetIndex = append(crlSetIndex, revokedCert.SerialNumber)
	}

	entry, err = logical.StorageEntryJSON("crlsets/index/"+name, crlSetIndex)
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(entry); err != nil {
		return nil, err
	}

	for _, revokedSerial := range crlSetIndex {
		entry, err = logical.StorageEntryJSON("crlsets/serial/"+revokedSerial.String(),
			&RevokedSerial{
				CRLSet: name,
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

func (b *backend) cleanIndex(storage logical.Storage, name string) {
	entry, err := storage.Get("crlsets/index/" + name)
	if err != nil {
		return
	}
	if entry == nil {
		return
	}

	crlSetIndex := []*big.Int{}
	err = entry.DecodeJSON(&crlSetIndex)
	if err != nil {
		goto destroyIndex
	}

	for _, serial := range crlSetIndex {
		storage.Delete("crlsets/serial/" + serial.String())
	}

destroyIndex:
	storage.Delete("crlsets/index/" + name)
	return
}

type RevokedSerial struct {
	CRLSet string `json:"crlset"`
}

//FIXME
const pathCRLSetsHelpSyn = `
Manage trusted certificates used for authentication.
`

const pathCRLSetsHelpDesc = `
This endpoint allows you to create, read, update, and delete trusted certificates
that are allowed to authenticate.

Deleting a certificate will not revoke auth for prior authenticated connections.
To do this, do a revoke on "login". If you don't need to revoke login immediately,
then the next renew will cause the lease to expire.
`
