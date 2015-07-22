package ssh

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretOTPType = "secret_otp_type"

func secretOTP(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretOTPType,
		Fields: map[string]*framework.FieldSchema{
			"otp": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "One time password",
			},
		},
		DefaultDuration:    1 * time.Hour,
		DefaultGracePeriod: 10 * time.Minute,
		Renew:              b.secretOTPRenew,
		Revoke:             b.secretOTPRevoke,
	}
}

func (b *backend) secretOTPRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}
	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, false)
	return f(req, d)
}

func (b *backend) secretOTPRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	otpRaw, ok := req.Secret.InternalData["otp"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	otp, ok := otpRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}

	otpSalted := b.salt.SaltID(otp)

	otpPath := fmt.Sprintf("otp/%s", otpSalted)
	err := req.Storage.Delete(otpPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
