/**
 *  Copyright 2014 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package otp

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"image"
	"net/url"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// Error when attempting to convert the secret from base32 to raw bytes.
var ErrValidateSecretInvalidBase32 = errors.New("Decoding of secret as base32 failed.")

// The user provided passcode length was not expected.
var ErrValidateInputInvalidLength = errors.New("Input length unexpected")

// When generating a Key, the Issuer must be set.
var ErrGenerateMissingIssuer = errors.New("Issuer must be set")

// When generating a Key, the Account Name must be set.
var ErrGenerateMissingAccountName = errors.New("AccountName must be set")

// Key represents an TOTP or HTOP key.
type Key struct {
	orig string
	url  *url.URL
}

// NewKeyFromURL creates a new Key from an TOTP or HOTP url.
//
// The URL format is documented here:
//   https://github.com/google/google-authenticator/wiki/Key-Uri-Format
//
func NewKeyFromURL(orig string) (*Key, error) {
	s := strings.TrimSpace(orig)

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return &Key{
		orig: s,
		url:  u,
	}, nil
}

func (k *Key) String() string {
	return k.orig
}

// Image returns an QR-Code image of the specified width and height,
// suitable for use by many clients like Google-Authenricator
// to enroll a user's TOTP/HOTP key.
func (k *Key) Image(width int, height int) (image.Image, error) {
	b, err := qr.Encode(k.orig, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	b, err = barcode.Scale(b, width, height)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// Type returns "hotp" or "totp".
func (k *Key) Type() string {
	return k.url.Host
}

// Issuer returns the name of the issuing organization.
func (k *Key) Issuer() string {
	q := k.url.Query()

	issuer := q.Get("issuer")

	if issuer != "" {
		return issuer
	}

	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return ""
	}

	return p[:i]
}

// AccountName returns the name of the user's account.
func (k *Key) AccountName() string {
	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return p
	}

	return p[i+1:]
}

// Secret returns the opaque secret for this Key.
func (k *Key) Secret() string {
	q := k.url.Query()

	return q.Get("secret")
}

// Period returns a tiny int representing the rotation time in seconds.
func (k *Key) Period() uint64 {
	q := k.url.Query()

	if u, err := strconv.ParseUint(q.Get("period"), 10, 64); err == nil {
		return u
	}

	// If no period is defined 30 seconds is the default per (rfc6238)
	return 30
}

// Digits returns a tiny int representing the number of OTP digits.
func (k *Key) Digits() Digits {
	q := k.url.Query()

	if u, err := strconv.ParseUint(q.Get("digits"), 10, 64); err == nil {
		switch u {
		case 8:
			return DigitsEight
		default:
			return DigitsSix
		}
	}

	// Six is the most common value.
	return DigitsSix
}

// Algorithm returns the algorithm used or the default (SHA1).
func (k *Key) Algorithm() Algorithm {
	q := k.url.Query()

	a := strings.ToLower(q.Get("algorithm"))
	switch a {
	case "md5":
		return AlgorithmMD5
	case "sha256":
		return AlgorithmSHA256
	case "sha512":
		return AlgorithmSHA512
	default:
		return AlgorithmSHA1
	}
}

// URL returns the OTP URL as a string
func (k *Key) URL() string {
	return k.url.String()
}

// Algorithm represents the hashing function to use in the HMAC
// operation needed for OTPs.
type Algorithm int

const (
	// AlgorithmSHA1 should be used for compatibility with Google Authenticator.
	//
	// See https://github.com/pquerna/otp/issues/55 for additional details.
	AlgorithmSHA1 Algorithm = iota
	AlgorithmSHA256
	AlgorithmSHA512
	AlgorithmMD5
)

func (a Algorithm) String() string {
	switch a {
	case AlgorithmSHA1:
		return "SHA1"
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmSHA512:
		return "SHA512"
	case AlgorithmMD5:
		return "MD5"
	}
	panic("unreached")
}

func (a Algorithm) Hash() hash.Hash {
	switch a {
	case AlgorithmSHA1:
		return sha1.New()
	case AlgorithmSHA256:
		return sha256.New()
	case AlgorithmSHA512:
		return sha512.New()
	case AlgorithmMD5:
		return md5.New()
	}
	panic("unreached")
}

// Digits represents the number of digits present in the
// user's OTP passcode. Six and Eight are the most common values.
type Digits int

const (
	DigitsSix   Digits = 6
	DigitsEight Digits = 8
)

// Format converts an integer into the zero-filled size for this Digits.
func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

// Length returns the number of characters for this Digits.
func (d Digits) Length() int {
	return int(d)
}

func (d Digits) String() string {
	return fmt.Sprintf("%d", d)
}
