package credentials

// Deprecated: the RSA key pair is deprecated
type RsaKeyPairCredential struct {
	PrivateKey        string
	PublicKeyId       string
	SessionExpiration int
}

// Deprecated: the RSA key pair is deprecated
func NewRsaKeyPairCredential(privateKey, publicKeyId string, sessionExpiration int) *RsaKeyPairCredential {
	return &RsaKeyPairCredential{
		PrivateKey:        privateKey,
		PublicKeyId:       publicKeyId,
		SessionExpiration: sessionExpiration,
	}
}
