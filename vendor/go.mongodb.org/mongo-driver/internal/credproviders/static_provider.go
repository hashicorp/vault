// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package credproviders

import (
	"errors"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)

// staticProviderName provides a name of Static provider
const staticProviderName = "StaticProvider"

// A StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct {
	credentials.Value

	verified bool
	err      error
}

func verify(v credentials.Value) error {
	if !v.HasKeys() {
		return errors.New("failed to retrieve ACCESS_KEY_ID and SECRET_ACCESS_KEY")
	}
	if v.AccessKeyID != "" && v.SecretAccessKey == "" {
		return errors.New("ACCESS_KEY_ID is set, but SECRET_ACCESS_KEY is missing")
	}
	if v.AccessKeyID == "" && v.SecretAccessKey != "" {
		return errors.New("SECRET_ACCESS_KEY is set, but ACCESS_KEY_ID is missing")
	}
	if v.AccessKeyID == "" && v.SecretAccessKey == "" && v.SessionToken != "" {
		return errors.New("AWS_SESSION_TOKEN is set, but ACCESS_KEY_ID and SECRET_ACCESS_KEY are missing")
	}
	return nil

}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s *StaticProvider) Retrieve() (credentials.Value, error) {
	if !s.verified {
		s.err = verify(s.Value)
		s.Value.ProviderName = staticProviderName
		s.verified = true
	}
	return s.Value, s.err
}

// IsExpired returns if the credentials are expired.
//
// For StaticProvider, the credentials never expired.
func (s *StaticProvider) IsExpired() bool {
	return false
}
