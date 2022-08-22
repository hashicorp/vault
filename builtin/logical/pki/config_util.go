package pki

import (
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
		err := sc.setIssuersConfig(&issuerConfigEntry{
			DefaultIssuerId: id,
		})
		if err != nil {
			return err
		}

		issuer, err := sc.fetchIssuerById(id)
		if err != nil {
			return err
		}

		issuer.LastModified = time.Now().In(time.FixedZone("GMT", 0))
		// err = sc.writeIssuer(issuer)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}
