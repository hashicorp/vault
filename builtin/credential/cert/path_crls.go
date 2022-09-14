package cert

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	url2 "net/url"
	"strings"
	"sync"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCRLs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "crls/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the certificate",
			},

			"crl": {
				Type: framework.TypeString,
				Description: `The public certificate that should be trusted.
May be DER or PEM encoded. Note: the expiration time
is ignored; if the CRL is no longer valid, delete it
using the same name as specified here.`,
			},
			"cdp": {
				Type:        framework.TypeString,
				Description: `The URL of a CRL distribution point.  Only one of crl or cdp parameters should be specified.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathCRLDelete,
			logical.ReadOperation:   b.pathCRLRead,
			logical.UpdateOperation: b.pathCRLWrite,
		},

		HelpSynopsis:    pathCRLsHelpSyn,
		HelpDescription: pathCRLsHelpDesc,
	}
}

func (b *backend) populateCRLs(ctx context.Context, storage logical.Storage) error {
	b.crlUpdateMutex.Lock()
	defer b.crlUpdateMutex.Unlock()

	if b.crls != nil {
		return nil
	}

	b.crls = map[string]CRLInfo{}

	keys, err := storage.List(ctx, "crls/")
	if err != nil {
		return fmt.Errorf("error listing CRLs: %w", err)
	}
	if keys == nil || len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		entry, err := storage.Get(ctx, "crls/"+key)
		if err != nil {
			b.crls = nil
			return fmt.Errorf("error loading CRL %q: %w", key, err)
		}
		if entry == nil {
			continue
		}
		var crlInfo CRLInfo
		err = entry.DecodeJSON(&crlInfo)
		if err != nil {
			b.crls = nil
			return fmt.Errorf("error decoding CRL %q: %w", key, err)
		}
		if crlInfo.CDP != nil {
			crlInfo.CDP.fetchOnce = &sync.Once{}
		}
		b.crls[key] = crlInfo
	}

	return nil
}

func (b *backend) findSerialInCRLs(serial *big.Int) map[string]RevokedSerialInfo {
	b.crlUpdateMutex.RLock()
	defer b.crlUpdateMutex.RUnlock()
	ret := map[string]RevokedSerialInfo{}
	for key, crl := range b.crls {
		if crl.Serials == nil {
			continue
		}
		if info, ok := crl.Serials[serial.String()]; ok {
			ret[key] = info
		}
	}
	return ret
}

func parseSerialString(input string) (*big.Int, error) {
	ret := &big.Int{}

	switch {
	case strings.Count(input, ":") > 0:
		serialBytes := certutil.ParseHexFormatted(input, ":")
		if serialBytes == nil {
			return nil, fmt.Errorf("error parsing serial %q", input)
		}
		ret.SetBytes(serialBytes)
	case strings.Count(input, "-") > 0:
		serialBytes := certutil.ParseHexFormatted(input, "-")
		if serialBytes == nil {
			return nil, fmt.Errorf("error parsing serial %q", input)
		}
		ret.SetBytes(serialBytes)
	default:
		var success bool
		ret, success = ret.SetString(input, 0)
		if !success {
			return nil, fmt.Errorf("error parsing serial %q", input)
		}
	}

	return ret, nil
}

func (b *backend) pathCRLDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(`"name" parameter cannot be empty`), nil
	}

	if err := b.populateCRLs(ctx, req.Storage); err != nil {
		return nil, err
	}

	b.crlUpdateMutex.Lock()
	defer b.crlUpdateMutex.Unlock()

	_, ok := b.crls[name]
	if !ok {
		return logical.ErrorResponse(fmt.Sprintf(
			"no such CRL %s", name,
		)), nil
	}

	if err := req.Storage.Delete(ctx, "crls/"+name); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"error deleting crl %s: %v", name, err),
		), nil
	}

	delete(b.crls, name)

	return nil, nil
}

func (b *backend) pathCRLRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(`"name" parameter must be set`), nil
	}

	if err := b.populateCRLs(ctx, req.Storage); err != nil {
		return nil, err
	}

	b.crlUpdateMutex.RLock()
	defer b.crlUpdateMutex.RUnlock()

	var retData map[string]interface{}

	crl, ok := b.crls[name]
	if !ok {
		return logical.ErrorResponse(fmt.Sprintf(
			"no such CRL %s", name,
		)), nil
	}

	retData = structs.New(&crl).Map()

	return &logical.Response{
		Data: retData,
	}, nil
}

func (b *backend) pathCRLWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(`"name" parameter cannot be empty`), nil
	}
	if crlRaw, ok := d.GetOk("crl"); ok {
		crl := crlRaw.(string)
		certList, err := x509.ParseCRL([]byte(crl))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse CRL: %v", err)), nil
		}
		if certList == nil {
			return logical.ErrorResponse("parsed CRL is nil"), nil
		}

		err = b.setCRL(ctx, req.Storage, certList, name, nil)
		if err != nil {
			return nil, err
		}
	} else if cdpRaw, ok := d.GetOk("cdp"); ok {
		cdl := cdpRaw.(string)
		if cdl == "" {
			return logical.ErrorResponse("empty CDP url"), nil
		}
		_, err := url2.Parse(cdl)
		if err != nil {
			return logical.ErrorResponse("invalid CDP url: %v", err), nil
		}
		crl := &CRLInfo{
			CDP: &CDPInfo{
				Url:       cdl,
				fetchOnce: &sync.Once{},
			},
		}
		err = b.fetchCRL(ctx, req.Storage, name, crl)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (b *backend) setCRL(ctx context.Context, storage logical.Storage, certList *pkix.CertificateList, name string, cdp *CDPInfo) error {
	if err := b.populateCRLs(ctx, storage); err != nil {
		return err
	}

	b.crlUpdateMutex.Lock()
	defer b.crlUpdateMutex.Unlock()

	crlInfo := CRLInfo{
		CDP:     cdp,
		Serials: map[string]RevokedSerialInfo{},
	}
	for _, revokedCert := range certList.TBSCertList.RevokedCertificates {
		crlInfo.Serials[revokedCert.SerialNumber.String()] = RevokedSerialInfo{}
	}

	entry, err := logical.StorageEntryJSON("crls/"+name, crlInfo)
	if err != nil {
		return err
	}
	if err = storage.Put(ctx, entry); err != nil {
		return err
	}

	b.crls[name] = crlInfo
	return err
}

type CDPInfo struct {
	Url        string     `json:"url" structs:"url" mapstructure:"url"`
	ValidUntil time.Time  `json:"valid_until" structs:"valid_until" mapstructure:"valid_until"`
	fetchOnce  *sync.Once `json:"-" structs:"-" mapstructure:"-"`
}

type CRLInfo struct {
	CDP     *CDPInfo                     `json:"cdp" structs:"cdp" mapstructure:"cdp"`
	Serials map[string]RevokedSerialInfo `json:"serials" structs:"serials" mapstructure:"serials"`
}

type RevokedSerialInfo struct{}

const pathCRLsHelpSyn = `
Manage Certificate Revocation Lists checked during authentication.
`

const pathCRLsHelpDesc = `
This endpoint allows you to create, read, update, and delete the Certificate
Revocation Lists checked during authentication.

When any CRLs are in effect, any login will check the trust chains sent by a
client against the submitted CRLs. Any chain containing a serial number revoked
by one or more of the CRLs causes that chain to be marked as invalid for the
authentication attempt. Conversely, *any* valid chain -- that is, a chain
in which none of the serials are revoked by any CRL -- allows authentication.
This allows authentication to succeed when interim parts of one chain have been
revoked; for instance, if a certificate is signed by two intermediate CAs due to
one of them expiring.
`
