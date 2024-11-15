// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwtauth

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type JWKSPair struct {
	JWKSUrl   string `mapstructure:"jwks_url"`
	JWKSCAPEM string `mapstructure:"jwks_ca_pem"`
}

func NewJWKSPairsConfig(jc *jwtConfig) ([]*JWKSPair, error) {
	if len(jc.JWKSPairs) <= 0 {
		return nil, nil
	}

	pairs := make([]*JWKSPair, 0, len(jc.JWKSPairs))
	for i := 0; i < len(jc.JWKSPairs); i++ {
		pairsMap, ok := jc.JWKSPairs[i].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("jwks_pairs must be provided as a list of json objects with the fields jwks_url and jwks_ca_pem")
		}
		jp, err := Initialize(pairsMap)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, jp)
	}

	return pairs, nil
}

func Initialize(jp map[string]interface{}) (*JWKSPair, error) {
	var newJp JWKSPair
	if err := mapstructure.Decode(jp, &newJp); err != nil {
		return nil, err
	}

	return &newJp, nil
}
