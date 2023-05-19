// Nonce is a class for generating and validating nonces loosely based off
// the design of Let's Encrypt's Boulder nonce service here:
//
//   https://github.com/letsencrypt/boulder/blob/main/nonce/nonce.go

package nonce

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
	nonceSentinel        = "vault0"
	nonceLength          = 6 + 8 + 16 + 16
	noncePlaintextLength = 16
)

type (
	ensTimestamp uint64
	ensCounter   uint64
)

type encryptedNonceService struct {
	validity time.Duration

	crypt cipher.AEAD

	nextCounter *atomic.Uint64

	issueLock      *sync.Mutex
	maxIssued      map[ensTimestamp]ensCounter
	minCounter     *atomic.Uint64
	redeemedTokens map[ensTimestamp]map[ensCounter]struct{}
}

func newEncryptedNonceService(validity time.Duration) (*encryptedNonceService, error) {
	// On startup, create a new AES key. This avoids having issues with the
	// number of encryptions we can do under this service.
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to initialize AES key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AES cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AES-GCM: %w", err)
	}

	return &encryptedNonceService{
		validity: validity,
		crypt:    aead,

		// nextCounter.Add(1) returns the _new_ value; by initializing to
		// zero, we guarantee that nextCounter = minCounter + 1 on the first
		// read; if it is redeemed right away, we then hold that invariant.
		nextCounter: new(atomic.Uint64),

		issueLock:      new(sync.Mutex),
		maxIssued:      make(map[ensTimestamp]ensCounter, validity/time.Second),
		minCounter:     new(atomic.Uint64),
		redeemedTokens: make(map[ensTimestamp]map[ensCounter]struct{}, validity/time.Second),
	}, nil
}

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
