// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oidc

import (
	"errors"
)

var (
	ErrInvalidParameter           = errors.New("invalid parameter")
	ErrNilParameter               = errors.New("nil parameter")
	ErrInvalidCACert              = errors.New("invalid CA certificate")
	ErrInvalidIssuer              = errors.New("invalid issuer")
	ErrExpiredRequest             = errors.New("request is expired")
	ErrInvalidResponseState       = errors.New("invalid response state")
	ErrInvalidSignature           = errors.New("invalid signature")
	ErrInvalidSubject             = errors.New("invalid subject")
	ErrInvalidAudience            = errors.New("invalid audience")
	ErrInvalidNonce               = errors.New("invalid nonce")
	ErrInvalidNotBefore           = errors.New("invalid not before")
	ErrExpiredToken               = errors.New("token is expired")
	ErrInvalidJWKs                = errors.New("invalid jwks")
	ErrInvalidIssuedAt            = errors.New("invalid issued at (iat)")
	ErrInvalidAuthorizedParty     = errors.New("invalid authorized party (azp)")
	ErrInvalidAtHash              = errors.New("access_token hash does not match value in id_token")
	ErrInvalidCodeHash            = errors.New("authorization code hash does not match value in id_token")
	ErrTokenNotSigned             = errors.New("token is not signed")
	ErrMalformedToken             = errors.New("token malformed")
	ErrUnsupportedAlg             = errors.New("unsupported signing algorithm")
	ErrIDGeneratorFailed          = errors.New("id generation failed")
	ErrMissingIDToken             = errors.New("id_token is missing")
	ErrMissingAccessToken         = errors.New("access_token is missing")
	ErrIDTokenVerificationFailed  = errors.New("id_token verification failed")
	ErrNotFound                   = errors.New("not found")
	ErrLoginFailed                = errors.New("login failed")
	ErrUserInfoFailed             = errors.New("user info failed")
	ErrUnauthorizedRedirectURI    = errors.New("unauthorized redirect_uri")
	ErrInvalidFlow                = errors.New("invalid OIDC flow")
	ErrUnsupportedChallengeMethod = errors.New("unsupported PKCE challenge method")
	ErrExpiredAuthTime            = errors.New("expired auth_time")
	ErrMissingClaim               = errors.New("missing required claim")
)
