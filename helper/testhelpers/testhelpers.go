package testhelpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

// Generates a root token on the target cluster.
func GenerateRoot(t testing.T, cluster *vault.TestCluster, drToken bool) string {
	token, err := GenerateRootWithError(t, cluster, drToken)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRootWithError(t testing.T, cluster *vault.TestCluster, drToken bool) (string, error) {
	buf := make([]byte, 16)
	readLen, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	if readLen != 16 {
		return "", fmt.Errorf("wrong readlen: %d", readLen)
	}
	otp := base64.StdEncoding.EncodeToString(buf)

	// If recovery keys supported, use those to perform root token generation instead
	var keys [][]byte
	if cluster.Cores[0].SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}

	client := cluster.Cores[0].Client
	f := client.Sys().GenerateRootInit
	if drToken {
		f = client.Sys().GenerateDROperationTokenInit
	}
	status, err := f(otp, "")
	if err != nil {
		return "", err
	}

	if status.Required > len(keys) {
		return "", fmt.Errorf("need more keys than have, need %d have %d", status.Required, len(keys))
	}

	for i, key := range keys {
		if i >= status.Required {
			break
		}
		f := client.Sys().GenerateRootUpdate
		if drToken {
			f = client.Sys().GenerateDROperationTokenUpdate
		}
		status, err = f(base64.StdEncoding.EncodeToString(key), status.Nonce)
		if err != nil {
			return "", err
		}
	}
	if !status.Complete {
		return "", errors.New("generate root operation did not end successfully")
	}
	tokenBytes, err := xor.XORBase64(status.EncodedToken, otp)
	if err != nil {
		return "", err
	}
	token, err := uuid.FormatUUID(tokenBytes)
	if err != nil {
		return "", err
	}
	return token, nil
}
