package pki

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	keyPrefix    = "/config/key/"
	issuerPrefix = "/config/issuer/"
)

var (
	emptyKey        = key{}
	emptyIssuer     = issuer{}
	emptyCertBundle = certutil.CertBundle{}
)

type keyId string

func (p keyId) String() string {
	return string(p)
}

type issuerId string

func (p issuerId) String() string {
	return string(p)
}

type key struct {
	ID             keyId                   `json:"id" structs:"id" mapstructure:"id"`
	PrivateKeyType certutil.PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	PrivateKey     string                  `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

type issuer struct {
	ID           issuerId `json:"id" structs:"id" mapstructure:"id"`
	Name         string   `json:"name" structs:"name" mapstructure:"name"`
	KeyID        keyId    `json:"key_id" structs:"key_id" mapstructure:"key_id"`
	Certificate  string   `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	CAChain      []string `json:"ca_chain" structs:"ca_chain" mapstructure:"ca_chain"`
	SerialNumber string   `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

func fetchKeyById(ctx context.Context, s logical.Storage, keyId keyId) (key, error) {
	keyEntry, err := s.Get(ctx, keyPrefix+keyId.String())
	if err != nil {
		return emptyKey, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: %v", err)}
	}
	if keyEntry == nil {
		// FIXME: Dedicated/specific error for this?
		return emptyKey, errutil.UserError{Err: fmt.Sprintf("pki key id %s does not exist", keyId.String())}
	}

	var key key
	if err := keyEntry.DecodeJSON(&key); err != nil {
		return emptyKey, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki key with id %s: %v", keyId.String(), err)}
	}

	return key, nil
}

func writeKey(ctx context.Context, s logical.Storage, key key) error {
	keyId := key.ID

	json, err := logical.StorageEntryJSON(keyPrefix+keyId.String(), key)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func fetchIssuerById(ctx context.Context, s logical.Storage, issuerId issuerId) (issuer, error) {
	issuerEntry, err := s.Get(ctx, issuerPrefix+issuerId.String())
	if err != nil {
		return emptyIssuer, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
	}
	if issuerEntry == nil {
		// FIXME: Dedicated/specific error for this?
		return emptyIssuer, errutil.UserError{Err: fmt.Sprintf("pki issuer id %s does not exist", issuerId.String())}
	}

	var issuer issuer
	if err := issuerEntry.DecodeJSON(&issuer); err != nil {
		return emptyIssuer, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki issuer with id %s: %v", issuerId.String(), err)}
	}

	return issuer, nil
}

func writeIssuer(ctx context.Context, s logical.Storage, issuer issuer) error {
	issuerId := issuer.ID

	json, err := logical.StorageEntryJSON(issuerPrefix+issuerId.String(), issuer)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func fetchCertBundleByIssuerId(ctx context.Context, s logical.Storage, issuerId issuerId) (certutil.CertBundle, error) {
	issuer, err := fetchIssuerById(ctx, s, issuerId)
	if err != nil {
		return emptyCertBundle, err
	}

	if issuer.KeyID == "" {
		return emptyCertBundle, errutil.UserError{Err: fmt.Sprintf("requested a cert bundle for an issuer id %s that did not contain a key", issuerId.String())}
	}

	key, err := fetchKeyById(ctx, s, issuer.KeyID)
	if err != nil {
		return emptyCertBundle, err
	}
	return certutil.CertBundle{
		PrivateKeyType: key.PrivateKeyType,
		Certificate:    issuer.Certificate,
		IssuingCA:      issuer.CAChain[0],
		CAChain:        issuer.CAChain,
		PrivateKey:     key.PrivateKey,
		SerialNumber:   issuer.SerialNumber,
	}, nil
}
