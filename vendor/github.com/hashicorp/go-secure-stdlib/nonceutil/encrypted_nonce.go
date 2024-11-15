// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Nonce is a class for generating and validating nonces loosely based off
// the design of Let's Encrypt's Boulder nonce service here:
//
//   https://github.com/letsencrypt/boulder/blob/main/nonce/nonce.go
//
// We use an encrypted tuple of (expiry timestamp, counter value), allowing
// us to maintain a map of only unexpired counter values that have been
// redeemed. This means that issuing nonces involves updating counter values
// and creating only up to a fixed-amount of memory (the size of the validity
// period) in the maxIssued map, whereas the sync.Map potentially grows
// indefinitely when also coupled with the fact that sync.Map never releases
// memory back to the host when its size has shrunk.
//
// Redeeming a nonce thus only stores the used counter value (8 bytes)
// and other checks for delayed or reused nonces remain as fast as parsing
// and decrypting the token value.

package nonceutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Internal, versioned sentinel to make sure our base64 data is truly
	// a nonce-like value.
	nonceSentinel = "vault0"

	// Wire length of the nonce, excluding raw url base64 encoding:
	//  - 6 byte sentinel (above),
	//  - 8 byte AES-GCM IV
	//  - 16 byte encrypted (timestamp, counter) tuple (1 AES block)
	//  - 16 byte AES-GCM tag.
	nonceLength = len(nonceSentinel) + 8 + 16 + 16

	// Length of the decrypted plaintext underlying the nonce:
	// - 8 byte expiry timestamp, unix seconds
	// - 8 byte incrementing counter value (uint64)
	noncePlaintextLength = 8 + 8
)

type (
	ensTimestamp uint64
	ensCounter   uint64
)

type encryptedNonceService struct {
	// How long a nonce is valid for. This directly correlates to memory
	// usage (retention of redeemed nonces).
	validity time.Duration

	// Underlying cipher for minting tokens.
	crypt cipher.AEAD

	// The next counter value to use for issuing, _after_ calling Add(1)
	// on it.
	nextCounter *atomic.Uint64

	// The remaining fields are locked by this read-only mutex. During
	// issuing a nonce, we update maxIssued; during redeeming we update
	// minCounter (an atomic) and redeemedTokens, and during tidy, we
	// potentially update update all fields.
	//
	// By storing maxIssued, we can (from our tidy run) update the
	// minCounter value when nonces were not redeemed recently, to make
	// any later redemptions fast (within a time period).
	//
	// The outer map in redeemedTokens and maxIssued map are of fixed size,
	// around the size of validity (in seconds). However, the internal
	// redeemedTokens[timestamp] maps may grow unbounded (assuming a
	// sufficiently fast system that can mint tokens infinitely fast).
	// However, once this timestamp expires, we can fully delete all
	// references to that map, and thus free up a potentially significant
	// chunk of memory.
	issueLock      *sync.Mutex
	maxIssued      map[ensTimestamp]ensCounter
	minCounter     *atomic.Uint64
	redeemedTokens map[ensTimestamp]map[ensCounter]struct{}
}

func newEncryptedNonceService(validity time.Duration) *encryptedNonceService {
	return &encryptedNonceService{
		validity: validity,

		// nextCounter.Add(1) returns the _new_ value; by initializing to
		// zero, we guarantee that nextCounter = minCounter + 1 on the first
		// read; if it is redeemed right away, we then hold that invariant.
		nextCounter: new(atomic.Uint64),

		issueLock:      new(sync.Mutex),
		maxIssued:      make(map[ensTimestamp]ensCounter, validity/time.Second),
		minCounter:     new(atomic.Uint64),
		redeemedTokens: make(map[ensTimestamp]map[ensCounter]struct{}, validity/time.Second),
	}
}

func (ens *encryptedNonceService) Initialize() error {
	// On initialization, create a new AES key. This avoids having issues
	// with the number of encryptions we can do under this service.
	//
	// Note that the nonce service will panic if this is not created.
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return fmt.Errorf("failed to initialize AES key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to initialize AES cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to initialize AES-GCM: %w", err)
	}

	ens.crypt = aead
	return nil
}

// This nonce service is strict (prohibits reuse of nonces even within
// the validity period) but is not cross-node: there would need to be
// an external communication mechanism to map nonce->node and only
// check for redemption there. Additionally, each initialization
// creates new key material and thus nonces from other nodes would
// not validate.

func (ens *encryptedNonceService) IsStrict() bool    { return true }
func (ens *encryptedNonceService) IsCrossNode() bool { return false }

func (ens *encryptedNonceService) encryptNonce(counter uint64, expiry time.Time) (token string, err error) {
	// counter is an 8-byte value and expiry (as a unix timestamp) is
	// likewise, so we have exactly one block of data.
	//
	// We encode in argument order, i.e., (counter, expiry), and in
	// big endian format.

	// Like Let's Encrypt, we use a 12-byte nonce with the leading four
	// bytes as zero and the remaining 8 bytes as zeros. This gives us
	// 2^(8*8)/2 = 2^32 birthday paradox limit to reuse a nonce. However,
	// note that nonce reuse (in AES-GCM) doesn't leak the key, only the
	// XOR of the plaintext. Here, as long as they can't forge valid nonces,
	// we're fine.
	nonce := make([]byte, 12)
	for i := 0; i < 4; i++ {
		nonce[i] = 0
	}
	if _, err := io.ReadFull(rand.Reader, nonce[4:]); err != nil {
		return "", fmt.Errorf("failed to read AEAD nonce: %w", err)
	}

	plaintext := make([]byte, noncePlaintextLength)
	binary.BigEndian.PutUint64(plaintext[0:], counter)
	binary.BigEndian.PutUint64(plaintext[8:], uint64(expiry.Unix()))
	ciphertext := ens.crypt.Seal(nil, nonce, plaintext, nil)

	// Now, generate the wire format of the nonce. Use a prefix, the nonce,
	// and then the ciphertext.
	var wire []byte
	wire = append(wire, []byte(nonceSentinel)...)
	wire = append(wire, nonce[4:]...)
	wire = append(wire, ciphertext...)

	if len(wire) != nonceLength {
		return "", fmt.Errorf("expected nonce length of %v got %v", nonceLength, len(wire))
	}

	return base64.RawURLEncoding.EncodeToString(wire), nil
}

func (ens *encryptedNonceService) recordCounterForTime(counter uint64, expiry time.Time) {
	timestamp := ensTimestamp(expiry.Unix())
	value := ensCounter(counter)

	ens.issueLock.Lock()
	defer ens.issueLock.Unlock()

	// This allows us to update minCounter when a given timestamp expires, if
	// we haven't seen all of that timestamp's nonces redeemed. Otherwise, we
	// could potentially be stuck at a lower counter value, making it harder
	// for us to check if nonces are redeemed quickly.

	lastValue, ok := ens.maxIssued[timestamp]
	if !ok || lastValue < value {
		ens.maxIssued[timestamp] = value
	}
}

func (ens *encryptedNonceService) Get() (token string, expiry time.Time, err error) {
	counter := ens.nextCounter.Add(1)
	now := time.Now()
	then := now.Add(ens.validity)

	token, err = ens.encryptNonce(counter, then)
	if err != nil {
		return "", now, err
	}

	ens.recordCounterForTime(counter, then)
	return token, then, nil
}

func (ens *encryptedNonceService) decryptNonce(token string) (counter uint64, expiry time.Time, ok bool) {
	zero := time.Time{}

	wire, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, zero, false
	}

	if len(wire) != nonceLength {
		return 0, zero, false
	}

	data := wire

	sentinel := data[0:len(nonceSentinel)]
	data = data[len(nonceSentinel):]
	if subtle.ConstantTimeCompare([]byte(nonceSentinel), sentinel) != 1 {
		return 0, zero, false
	}

	nonce := make([]byte, 12)
	for i := 0; i < 4; i++ {
		nonce[i] = 0
	}
	copy(nonce[4:12], data[0:8])
	data = data[8:]

	ciphertext := data[:]

	plaintext, err := ens.crypt.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, zero, false
	}

	if len(plaintext) != noncePlaintextLength {
		return 0, zero, false
	}

	counter = binary.BigEndian.Uint64(plaintext[0:8])
	unix := binary.BigEndian.Uint64(plaintext[8:])
	expiry = time.Unix(int64(unix), 0)

	return counter, expiry, true
}

func (ens *encryptedNonceService) Redeem(token string) bool {
	now := time.Now()
	counter, expiry, ok := ens.decryptNonce(token)
	if !ok {
		return false
	}

	if expiry.Before(now) {
		return false
	}

	if counter <= ens.minCounter.Load() {
		return false
	}

	timestamp := ensTimestamp(expiry.Unix())
	counterValue := ensCounter(counter)

	// From here on out, we're doing the expensive checks. This _looks_
	// like a valid token, but now we want to verify the used-exactly-once
	// nature.
	ens.issueLock.Lock()
	defer ens.issueLock.Unlock()

	minCounter := ens.minCounter.Load()
	if counter <= minCounter {
		// Someone else redeemed this token or time has rolled over before we
		// grabbed this lock. Reject this token.
		return false
	}

	// Check if this has already been redeemed.
	timestampMap, present := ens.redeemedTokens[timestamp]
	if !present {
		// No tokens have been redeemed from this token. Provision the
		// timestamp-specific map, but wait to see if we need to add into
		// it.
		timestampMap = make(map[ensCounter]struct{})
		ens.redeemedTokens[timestamp] = timestampMap
	}

	_, present = timestampMap[counterValue]
	if present {
		// Token was already redeemed. Reject this request.
		return false
	}

	// From here on out, the token is valid. Let's start by seeing if we can
	// free any memory usage.
	minCounter = ens.tidyMemoryHoldingLock(now, minCounter)

	// Before we add to the map, we should see if we can save memory by just
	// incrementing the minimum accepted by one, instead of adding to the
	// timestamp for out of order redemption.
	if minCounter+1 == counter {
		minCounter = counter
	} else {
		// Otherwise, we've got to flag this counter as valid.
		timestampMap[counterValue] = struct{}{}
	}

	// Finally, update our value of minCounter because we held the lock.
	ens.minCounter.Store(minCounter)

	return true
}

func (ens *encryptedNonceService) tidyMemoryHoldingLock(now time.Time, minCounter uint64) uint64 {
	// Quick and dirty tidy: any expired timestamps should be deleted, which
	// should free the most memory (relatively speaking, given a uniform
	// usage pattern). This also avoids an expensive iteration over all
	// redeemed counter values.
	//
	// First tidy the redeemed tokens, as that is the largest value.
	var deleteCandidates []ensTimestamp
	for candidate := range ens.redeemedTokens {
		if candidate < ensTimestamp(now.Unix()) {
			deleteCandidates = append(deleteCandidates, candidate)
		}
	}
	for _, candidate := range deleteCandidates {
		delete(ens.redeemedTokens, candidate)
	}

	// Then tidy the last used timestamp values. Here, any removed timestamps
	// have an expiry time before now, which means they cannot be used. This
	// means our minCounterValue, if it is
	deleteCandidates = nil
	for candidate, lastIssuedInTimestamp := range ens.maxIssued {
		if candidate < ensTimestamp(now.Unix()) {
			deleteCandidates = append(deleteCandidates, candidate)
			if lastIssuedInTimestamp > ensCounter(minCounter) {
				minCounter = uint64(lastIssuedInTimestamp)
			}
		}
	}
	for _, candidate := range deleteCandidates {
		delete(ens.maxIssued, candidate)
	}
	return minCounter
}

func (ens *encryptedNonceService) tidySequentialNonces(now time.Time, minCounter uint64) uint64 {
	// This potentially slow sequential tidy allows us to free up an
	// incremental amount of memory when out-of-order (common) redemption
	// occurs. The underlying map may not shrink, but it should have
	// additional capacity to handle additional redemptions of cohort
	// nonces without additional allocations _if_ this tidy works.
	//
	// This is made possible by updating the minCounter based on the
	// earlier maxIssued map and tries to maintain the fast-case invariant
	// described in newEncryptedNonceService(...).
	var timestamps []ensTimestamp
	for timestamp := range ens.redeemedTokens {
		timestamps = append(timestamps, timestamp)
	}

	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i] < timestamps[j] })
	var deleteCandidates []ensTimestamp
	for _, timestamp := range timestamps {
		counters := ens.redeemedTokens[timestamp]
		for len(counters) > 0 {
			_, present := counters[ensCounter(minCounter+1)]
			if !present {
				return minCounter
			}

			minCounter += 1
			delete(counters, ensCounter(minCounter))
		}

		if len(counters) == 0 {
			deleteCandidates = append(deleteCandidates, timestamp)
		}
	}
	for _, candidate := range deleteCandidates {
		delete(ens.redeemedTokens, candidate)
	}

	return minCounter
}

func (ens *encryptedNonceService) getMessage(lock time.Duration, memory time.Duration, sequential time.Duration) string {
	now := time.Now()
	var message string
	message += fmt.Sprintf("len(ens.maxIssued): %v\n", len(ens.maxIssued))
	message += fmt.Sprintf("len(ens.redeemedTokens): %v\n", len(ens.redeemedTokens))

	var total int
	for timestamp, counters := range ens.redeemedTokens {
		message += fmt.Sprintf("    ens.redeemedTokens[%v]: %v\n", timestamp, len(counters))
		total += len(counters)
	}
	build := time.Now()

	message += fmt.Sprintf("total redeemed tokens: %v\n", total)
	message += fmt.Sprintf("time to grab lock: %v\n", lock)
	message += fmt.Sprintf("time to tidy memory: %v\n", memory)
	message += fmt.Sprintf("time to tidy sequential: %v\n", sequential)
	message += fmt.Sprintf("time to build message: %v\n", build.Sub(now))
	return message
}

func (ens *encryptedNonceService) Tidy() *NonceStatus {
	lockStart := time.Now()
	ens.issueLock.Lock()
	defer ens.issueLock.Unlock()
	lockEnd := time.Now()

	minCounter := ens.minCounter.Load()

	now := time.Now()
	minCounter = ens.tidyMemoryHoldingLock(now, minCounter)
	memory := time.Now()
	minCounter = ens.tidySequentialNonces(now, minCounter)
	sequential := time.Now()
	ens.minCounter.Store(minCounter)

	issued := ens.nextCounter.Load()
	return &NonceStatus{
		Issued:      issued,
		Outstanding: issued - minCounter,
		Message:     ens.getMessage(lockEnd.Sub(lockStart), memory.Sub(now), sequential.Sub(memory)),
	}
}
