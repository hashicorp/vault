package pac

import (
	"bytes"

	"github.com/jcmturner/gokrb5/v8/iana/chksumtype"
	"github.com/jcmturner/rpc/v2/mstypes"
)

/*
https://msdn.microsoft.com/en-us/library/cc237955.aspx

The Key Usage Value MUST be KERB_NON_KERB_CKSUM_SALT (17) [MS-KILE] (section 3.1.5.9).

Server Signature (SignatureType = 0x00000006)
https://msdn.microsoft.com/en-us/library/cc237957.aspx
The KDC will use the long-term key that the KDC shares with the server, so that the server can verify this signature on receiving a PAC.
The server signature is a keyed hash [RFC4757] of the entire PAC message, with the Signature fields of both PAC_SIGNATURE_DATA structures set to zero.
The key used to protect the ciphertext part of the response is used.
The checksum type corresponds to the key unless the key is DES, in which case the KERB_CHECKSUM_HMAC_MD5 key is used.
The resulting hash value is then placed in the Signature field of the server's PAC_SIGNATURE_DATA structure.

KDC Signature (SignatureType = 0x00000007)
https://msdn.microsoft.com/en-us/library/dd357117.aspx
The KDC will use KDC (krbtgt) key [RFC4120], so that other KDCs can verify this signature on receiving a PAC.
The KDC signature is a keyed hash [RFC4757] of the Server Signature field in the PAC message.
The cryptographic system that is used to calculate the checksum depends on which system the KDC supports, as defined below:
- Supports RC4-HMAC --> KERB_CHECKSUM_HMAC_MD5
- Does not support RC4-HMAC and supports AES256 --> HMAC_SHA1_96_AES256
- Does not support RC4-HMAC or AES256-CTS-HMAC-SHA1-96, and supports AES128-CTS-HMAC-SHA1-96 --> HMAC_SHA1_96_AES128
- Does not support RC4-HMAC, AES128-CTS-HMAC-SHA1-96 or AES256-CTS-HMAC-SHA1-96 -->  None. The checksum operation will fail.
*/

// SignatureData implements https://msdn.microsoft.com/en-us/library/cc237955.aspx
type SignatureData struct {
	SignatureType  uint32 // A 32-bit unsigned integer value in little-endian format that defines the cryptographic system used to calculate the checksum. This MUST be one of the following checksum types: KERB_CHECKSUM_HMAC_MD5 (signature size = 16), HMAC_SHA1_96_AES128 (signature size = 12), HMAC_SHA1_96_AES256 (signature size = 12).
	Signature      []byte // Size depends on the type. See comment above.
	RODCIdentifier uint16 // A 16-bit unsigned integer value in little-endian format that contains the first 16 bits of the key version number ([MS-KILE] section 3.1.5.8) when the KDC is an RODC. When the KDC is not an RODC, this field does not exist.
}

// Unmarshal bytes into the SignatureData struct
func (k *SignatureData) Unmarshal(b []byte) (rb []byte, err error) {
	r := mstypes.NewReader(bytes.NewReader(b))

	k.SignatureType, err = r.Uint32()
	if err != nil {
		return
	}

	var c int
	switch k.SignatureType {
	case chksumtype.KERB_CHECKSUM_HMAC_MD5_UNSIGNED:
		c = 16
	case uint32(chksumtype.HMAC_SHA1_96_AES128):
		c = 12
	case uint32(chksumtype.HMAC_SHA1_96_AES256):
		c = 12
	}
	k.Signature, err = r.ReadBytes(c)
	if err != nil {
		return
	}

	// When the KDC is not an Read Only Domain Controller (RODC), this field does not exist.
	if len(b) >= 4+c+2 {
		k.RODCIdentifier, err = r.Uint16()
		if err != nil {
			return
		}
	}

	// Create bytes with zeroed signature needed for checksum verification
	rb = make([]byte, len(b), len(b))
	copy(rb, b)
	z := make([]byte, len(b), len(b))
	copy(rb[4:4+c], z)

	return
}
