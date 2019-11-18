package seal

import (
	wrapping "github.com/hashicorp/go-kms-wrapping"
)

func NewTestSeal(secret []byte) *Access {
	return &Access{
		Wrapper: wrapping.NewTestWrapper(secret),
	}
}
