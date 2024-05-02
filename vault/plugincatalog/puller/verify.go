//

package puller

import (
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func verifySignature(signedFile, signature []byte) error {
	key, err := crypto.NewKeyFromArmored(hcPubKey)
	if err != nil {
		return err
	}
	ring, err := crypto.NewKeyRing(key)
	if err != nil {
		return err
	}
	sigBin := crypto.NewPGPSignature(signature)
	err = ring.VerifyDetached(crypto.NewPlainMessage(signedFile), sigBin, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}
