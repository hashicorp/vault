package pki

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pkiKeyPrefix    = "/config/key/"
	pkiIssuerPrefix = "/config/issuer/"
)

var (
	emptyPkiKey = pkiKey{}
	emptyIssuer = pkiIssuer{}
)

type pkiKeyId string

func (p pkiKeyId) String() string {
	return string(p)
}

type pkiIssuerId string

func (p pkiIssuerId) String() string {
	return string(p)
}

type pkiKey struct {
	ID             pkiKeyId                `json:"id" structs:"id" mapstructure:"id"`
	PrivateKeyType certutil.PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	PrivateKey     string                  `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

type pkiIssuer struct {
	ID           pkiIssuerId `json:"id" structs:"id" mapstructure:"id"`
	PKIKeyID     pkiKeyId    `json:"pki_key_id" structs:"pki_key_id" mapstructure:"pki_key_id"`
	Certificate  string      `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	CAChain      []string    `json:"ca_chain" structs:"ca_chain" mapstructure:"ca_chain"`
	SerialNumber string      `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

func fetchPKIKeyById(ctx context.Context, s logical.Storage, keyId pkiKeyId) (pkiKey, error) {
	keyEntry, err := s.Get(ctx, pkiKeyPrefix+keyId.String())
	if err != nil {
		return emptyPkiKey, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: %v", err)}
	}
	if keyEntry == nil {
		// FIXME: Dedicated/specific error for this?
		return emptyPkiKey, errutil.UserError{Err: fmt.Sprintf("pki key id %s does not exist", keyId.String())}
	}

	var pkiKey pkiKey
	if err := keyEntry.DecodeJSON(&pkiKey); err != nil {
		return emptyPkiKey, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki key with id %s: %v", keyId.String(), err)}
	}

	return pkiKey, nil
}

func writePKIKey(ctx context.Context, s logical.Storage, pkiKey pkiKey) error {
	keyId := pkiKey.ID

	json, err := logical.StorageEntryJSON(pkiKeyPrefix+keyId.String(), pkiKey)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func fetchPKIIssuerById(ctx context.Context, s logical.Storage, issuerId pkiIssuerId) (pkiIssuer, error) {
	issuerEntry, err := s.Get(ctx, pkiIssuerPrefix+issuerId.String())
	if err != nil {
		return emptyIssuer, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
	}
	if issuerEntry == nil {
		// FIXME: Dedicated/specific error for this?
		return emptyIssuer, errutil.UserError{Err: fmt.Sprintf("pki issuer id %s does not exist", issuerId.String())}
	}

	var pkiIssuer pkiIssuer
	if err := issuerEntry.DecodeJSON(&pkiIssuer); err != nil {
		return emptyIssuer, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki issuer with id %s: %v", issuerId.String(), err)}
	}

	return pkiIssuer, nil
}

func writePKIIssuer(ctx context.Context, s logical.Storage, pkiIssuer pkiIssuer) error {
	issuerId := pkiIssuer.ID

	json, err := logical.StorageEntryJSON(pkiIssuerPrefix+issuerId.String(), pkiIssuer)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}
