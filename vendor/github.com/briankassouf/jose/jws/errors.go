package jws

import "errors"

var (

	// ErrNotEnoughMethods is returned if New was called _or_ the Flat/Compact
	// methods were called with 0 SigningMethods.
	ErrNotEnoughMethods = errors.New("not enough methods provided")

	// ErrCouldNotUnmarshal is returned when Parse's json.Unmarshaler
	// parameter returns an error.
	ErrCouldNotUnmarshal = errors.New("custom unmarshal failed")

	// ErrNotCompact signals that the provided potential JWS is not
	// in its compact representation.
	ErrNotCompact = errors.New("not a compact JWS")

	// ErrDuplicateHeaderParameter signals that there are duplicate parameters
	// in the provided Headers.
	ErrDuplicateHeaderParameter = errors.New("duplicate parameters in the JOSE Header")

	// ErrTwoEmptyHeaders is returned if both Headers are empty.
	ErrTwoEmptyHeaders = errors.New("both headers cannot be empty")

	// ErrNotEnoughKeys is returned when not enough keys are provided for
	// the given SigningMethods.
	ErrNotEnoughKeys = errors.New("not enough keys (for given methods)")

	// ErrDidNotValidate means the given JWT did not properly validate
	ErrDidNotValidate = errors.New("did not validate")

	// ErrNoAlgorithm means no algorithm ("alg") was found in the Protected
	// Header.
	ErrNoAlgorithm = errors.New("no algorithm found")

	// ErrAlgorithmDoesntExist means the algorithm asked for cannot be
	// found inside the signingMethod cache.
	ErrAlgorithmDoesntExist = errors.New("algorithm doesn't exist")

	// ErrMismatchedAlgorithms means the algorithm inside the JWT was
	// different than the algorithm the caller wanted to use.
	ErrMismatchedAlgorithms = errors.New("mismatched algorithms")

	// ErrCannotValidate means the JWS cannot be validated for various
	// reasons. For example, if there aren't any signatures/payloads/headers
	// to actually validate.
	ErrCannotValidate = errors.New("cannot validate")

	// ErrIsNotJWT means the given JWS is not a JWT.
	ErrIsNotJWT = errors.New("JWS is not a JWT")

	// ErrHoldsJWE means the given JWS holds a JWE inside its payload.
	ErrHoldsJWE = errors.New("JWS holds JWE")

	// ErrNotEnoughValidSignatures means the JWS did not meet the required
	// number of signatures.
	ErrNotEnoughValidSignatures = errors.New("not enough valid signatures in the JWS")

	// ErrNoTokenInRequest means there's no token present inside the *http.Request.
	ErrNoTokenInRequest = errors.New("no token present in request")
)
