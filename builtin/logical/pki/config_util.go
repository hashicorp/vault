package pki

import (
	"fmt"
	"strings"
	"time"
)

func (sc *storageContext) isDefaultKeySet() (bool, error) {
	config, err := sc.getKeysConfig()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultKeyId.String()) != "", nil
}

func (sc *storageContext) isDefaultIssuerSet() (bool, error) {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultIssuerId.String()) != "", nil
}

func (sc *storageContext) updateDefaultKeyId(id keyID) error {
	config, err := sc.getKeysConfig()
	if err != nil {
		return err
	}

	if config.DefaultKeyId != id {
		return sc.setKeysConfig(&keyConfigEntry{
			DefaultKeyId: id,
		})
	}

	return nil
}

func (sc *storageContext) updateDefaultIssuerId(id issuerID) error {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return err
	}

	if config.DefaultIssuerId != id {
		oldDefault := config.DefaultIssuerId
		newDefault := id
		now := time.Now().UTC()

		err := sc.setIssuersConfig(&issuerConfigEntry{
			DefaultIssuerId: newDefault,
		})
		if err != nil {
			return err
		}

		// When the default issuer changes, we need to modify four
		// pieces of information:
		//
		// 1. The old default issuer's modification time, as it no
		//    longer works for the /cert/ca path.
		// 2. The new default issuer's modification time, as it now
		//    works for the /cert/ca path.
		// 3. & 4. Both issuer's CRLs, as they behave the same, under
		//    the /cert/crl path!
		for _, thisId := range []issuerID{oldDefault, newDefault} {
			if len(thisId) == 0 {
				continue
			}

			// 1 & 2 above.
			issuer, err := sc.fetchIssuerById(thisId)
			if err != nil {
				return fmt.Errorf("unable to update issuer (%v)'s modification time: error fetching issuer: %v", thisId, err)
			}

			issuer.LastModified = now
			err = sc.writeIssuer(issuer)
			if err != nil {
				return fmt.Errorf("unable to update issuer (%v)'s modification time: error persisting issuer: %v", thisId, err)
			}
		}

		// Fetch and update the localCRLConfigEntry (3&4).
		cfg, err := sc.getLocalCRLConfig()
		if err != nil {
			return fmt.Errorf("unable to update local CRL config's modification time: error fetching local CRL config: %v", err)
		}

		cfg.LastModified = now
		cfg.DeltaLastModified = now
		err = sc.setLocalCRLConfig(cfg)
		if err != nil {
			return fmt.Errorf("unable to update local CRL config's modification time: error persisting local CRL config: %v", err)
		}
	}

	return nil
}
