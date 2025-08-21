package keysutil

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"

	"github.com/cloudflare/circl/kem"
)

func generateKeyPairOrFatal(t *testing.T, box kyberBox) (kem.PublicKey, kem.PrivateKey) {
	t.Helper()
	publicKey, privateKey, err := box.s.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return publicKey, privateKey
}

func TestKyberBox_RoundTrip_WithAssociatedData(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("hello, kyber + aes-gcm 🚀")
	associatedData := []byte("context-bound-AD")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	decryptedPlaintext, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, associatedData)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(decryptedPlaintext, plaintext) {
		t.Fatalf("decrypted plaintext mismatch: got %q want %q", decryptedPlaintext, plaintext)
	}
}

func TestKyberBox_RoundTrip_EmptyPlaintext(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	var plaintext []byte
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	decryptedPlaintext, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, associatedData)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(decryptedPlaintext, plaintext) {
		t.Fatalf("decrypted plaintext mismatch: got %q want %q", decryptedPlaintext, plaintext)
	}
}

func TestKyberBox_AssociatedDataMismatch_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("secret")
	validAssociatedData := []byte("ctx")
	invalidAssociatedData := []byte("ctx2")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, validAssociatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, invalidAssociatedData); err == nil {
		t.Fatalf("Decrypt succeeded with wrong associated data; expected failure")
	}
}

func TestKyberBox_AssociatedData_NilEqualsEmpty(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("nil vs empty AD")
	var nilAssociatedData []byte
	emptyAssociatedData := []byte{}

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, nilAssociatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Decrypt using empty (non-nil) AD; should still succeed because the implementation
	// binds to H(AD), and H(nil) == H([]byte{}).
	decryptedPlaintext, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, emptyAssociatedData)
	if err != nil {
		t.Fatalf("Decrypt with empty AD failed; expected success: %v", err)
	}
	if !bytes.Equal(decryptedPlaintext, plaintext) {
		t.Fatalf("decrypted plaintext mismatch: got %q want %q", decryptedPlaintext, plaintext)
	}
}

func TestKyberBox_InvalidNonceLength(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("nonce length test")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(nonce) == 0 {
		t.Fatalf("unexpected zero-length nonce")
	}

	// Truncate nonce to force failure.
	truncatedNonce := nonce[:len(nonce)-1]
	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, truncatedNonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with invalid nonce length; expected error")
	}
}

func TestKyberBox_InvalidCapsuleLength(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("capsule length test")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(kemCapsule) == 0 {
		t.Fatalf("unexpected zero-length capsule")
	}

	// Truncate capsule to force failure at input validation.
	truncatedCapsule := kemCapsule[:len(kemCapsule)-1]
	if _, err := kyberBox.Decrypt(privateKey, truncatedCapsule, nonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with invalid capsule length; expected error")
	}
}

func TestKyberBox_TamperCiphertext_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("auth test")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Fatalf("unexpected zero-length ciphertext")
	}

	ciphertext[0] ^= 0x01 // flip a bit to simulate tampering
	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded after tampering; expected failure")
	}
}

func TestKyberBox_RandomizedOutputs(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("same plaintext, different outputs")
	associatedData := []byte("ctx")

	capsuleOne, nonceOne, ciphertextOne, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt #1: %v", err)
	}
	capsuleTwo, nonceTwo, ciphertextTwo, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt #2: %v", err)
	}

	if bytes.Equal(capsuleOne, capsuleTwo) && bytes.Equal(nonceOne, nonceTwo) && bytes.Equal(ciphertextOne, ciphertextTwo) {
		t.Fatalf("encryption outputs should be randomized; got identical capsule+nonce+ciphertext")
	}

	// Sanity: both decrypt fine
	if _, err := kyberBox.Decrypt(privateKey, capsuleOne, nonceOne, ciphertextOne, associatedData); err != nil {
		t.Fatalf("Decrypt #1: %v", err)
	}
	if _, err := kyberBox.Decrypt(privateKey, capsuleTwo, nonceTwo, ciphertextTwo, associatedData); err != nil {
		t.Fatalf("Decrypt #2: %v", err)
	}
}

func TestKyberBox_ConcurrentEncryptDecrypt(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	associatedData := []byte("ctx")
	type testCase struct {
		plaintext []byte
	}
	const numberOfWorkers = 16

	testVectors := make([]testCase, numberOfWorkers)
	for i := range testVectors {
		randomPlaintext := make([]byte, 64)
		if _, err := rand.Read(randomPlaintext); err != nil {
			t.Fatalf("rand.Read: %v", err)
		}
		testVectors[i] = testCase{plaintext: randomPlaintext}
	}

	errorCh := make(chan error, numberOfWorkers)
	for i := range testVectors {
		index := i
		go func() {
			kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, testVectors[index].plaintext, associatedData)
			if err != nil {
				errorCh <- err
				return
			}
			decryptedPlaintext, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, associatedData)
			if err != nil {
				errorCh <- err
				return
			}
			if !bytes.Equal(decryptedPlaintext, testVectors[index].plaintext) {
				errorCh <- errRoundTripMismatch("round-trip mismatch")
				return
			}
			errorCh <- nil
		}()
	}

	timeout := time.After(5 * time.Second)
	for range testVectors {
		select {
		case err := <-errorCh:
			if err != nil {
				t.Fatalf("concurrent round-trip error: %v", err)
			}
		case <-timeout:
			t.Fatalf("concurrent test timed out")
		}
	}
}

func TestKyberBox_DecryptionWithWrongPrivateKey_Fails(t *testing.T) {
	kyberBox := newKyberBox()

	// Correct keypair used for encryption
	publicKeyCorrect, _ := generateKeyPairOrFatal(t, kyberBox)
	// Different (wrong) private key used for decryption
	_, privateKeyWrong := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("decrypt with wrong key should fail")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKeyCorrect, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if _, err := kyberBox.Decrypt(privateKeyWrong, kemCapsule, nonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with wrong private key; expected failure")
	}
}

func TestKyberBox_TamperNonce_ValidLengthButBitFlipped_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("tamper nonce should fail")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Copy and flip one bit in the nonce (keep length valid)
	modifiedNonce := append([]byte(nil), nonce...)
	modifiedNonce[0] ^= 0x80

	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, modifiedNonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with tampered nonce; expected failure")
	}
}

func TestKyberBox_TruncatedCiphertext_TagRemoved_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("truncate ct should fail")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(ciphertext) < 2 {
		t.Fatalf("ciphertext unexpectedly short for truncation test")
	}

	truncatedCiphertext := ciphertext[:len(ciphertext)-1] // remove 1 byte (likely in the tag)

	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, truncatedCiphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with truncated ciphertext; expected failure")
	}
}

func TestKyberBox_NonceUniquenessAcrossEncryptions(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, _ := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("check nonce uniqueness")
	associatedData := []byte("ctx")

	_, nonceOne, _, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt #1: %v", err)
	}
	_, nonceTwo, _, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt #2: %v", err)
	}

	if bytes.Equal(nonceOne, nonceTwo) {
		t.Fatalf("nonces should differ across encryptions; got identical nonces")
	}
}

func TestKyberBox_RoundTrip_WithLargeAssociatedData(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("large AD round trip")
	largeAssociatedData := make([]byte, 64*1024) // 64 KiB AD
	if _, err := rand.Read(largeAssociatedData); err != nil {
		t.Fatalf("rand.Read(AD): %v", err)
	}

	capsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, largeAssociatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	decryptedPlaintext, err := kyberBox.Decrypt(privateKey, capsule, nonce, ciphertext, largeAssociatedData)
	if err != nil {
		t.Fatalf("Decrypt with large AD failed: %v", err)
	}
	if !bytes.Equal(decryptedPlaintext, plaintext) {
		t.Fatalf("decrypted plaintext mismatch with large AD: got %q want %q", decryptedPlaintext, plaintext)
	}
}

func TestKyberBox_TamperCapsule_BitFlipSameLength_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("capsule bit flip should fail")
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Flip one bit but keep capsule length identical.
	modifiedCapsule := append([]byte(nil), kemCapsule...)
	modifiedCapsule[len(modifiedCapsule)/2] ^= 0x01

	if _, err := kyberBox.Decrypt(privateKey, modifiedCapsule, nonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with tampered capsule; expected failure")
	}
}

func TestNewGCM_InvalidKeyLength_ReturnsError(t *testing.T) {
	// newGCM should reject any key not 16/24/32 bytes for AES.
	invalidKey := make([]byte, 31) // 31 bytes is invalid for AES
	if _, err := newGCM(invalidKey); err == nil {
		t.Fatalf("newGCM accepted invalid AES key length; expected error")
	}
}

func TestNewGCMWithNonce_InvalidKeyLength_ReturnsError(t *testing.T) {
	invalidKey := make([]byte, 0) // definitely invalid
	if _, _, err := newGCMWithNonce(invalidKey); err == nil {
		t.Fatalf("newGCMWithNonce accepted invalid AES key length; expected error")
	}
}

func TestKyberBox_CapsuleSwapBetweenTwoEncryptions_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintextA := []byte("message A")
	plaintextB := []byte("message B")
	associatedData := []byte("ctx")

	capsuleA, nonceA, ciphertextA, err := kyberBox.Encrypt(publicKey, plaintextA, associatedData)
	if err != nil {
		t.Fatalf("Encrypt A: %v", err)
	}
	capsuleB, _, ciphertextB, err := kyberBox.Encrypt(publicKey, plaintextB, associatedData)
	if err != nil {
		t.Fatalf("Encrypt B: %v", err)
	}

	// Swap capsules; key derivation binds to capsule+AD, so this must fail.
	if _, err := kyberBox.Decrypt(privateKey, capsuleB, nonceA, ciphertextA, associatedData); err == nil {
		t.Fatalf("Decrypt succeeded with swapped capsule; expected failure")
	}

	// Sanity checks still succeed.
	if _, err := kyberBox.Decrypt(privateKey, capsuleA, nonceA, ciphertextA, associatedData); err != nil {
		t.Fatalf("Sanity decrypt A failed: %v", err)
	}
	if _, err := kyberBox.Decrypt(privateKey, capsuleB, nonceA /* wrong nonce */, ciphertextB, associatedData); err == nil {
		// Note: also likely fails because nonce belongs to A.
		t.Fatalf("Decrypt unexpectedly succeeded with mismatched nonce")
	}
}

func TestKyberBox_RoundTrip_WithLargePlaintext_OneMiB(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	largePlaintext := make([]byte, 1<<20) // 1 MiB
	if _, err := rand.Read(largePlaintext); err != nil {
		t.Fatalf("rand.Read(largePlaintext): %v", err)
	}
	associatedData := []byte("ctx")

	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, largePlaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	decryptedPlaintext, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, ciphertext, associatedData)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(decryptedPlaintext, largePlaintext) {
		t.Fatalf("large round-trip mismatch: got %d bytes want %d bytes", len(decryptedPlaintext), len(largePlaintext))
	}
}

func TestKyberBox_Decrypt_WithEmptyCiphertext_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("will not be used")
	associatedData := []byte("ctx")
	kemCapsule, nonce, _, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, nonce, []byte{}, associatedData); err == nil {
		t.Fatalf("Decrypt unexpectedly succeeded with empty ciphertext")
	}
}

func TestKyberBox_Decrypt_WithTooLongNonce_Fails(t *testing.T) {
	kyberBox := newKyberBox()
	publicKey, privateKey := generateKeyPairOrFatal(t, kyberBox)

	plaintext := []byte("nonce size check")
	associatedData := []byte("ctx")
	kemCapsule, nonce, ciphertext, err := kyberBox.Encrypt(publicKey, plaintext, associatedData)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Extend nonce by 1 byte to violate aead.NonceSize().
	extendedNonce := append(append([]byte(nil), nonce...), 0x00)

	if _, err := kyberBox.Decrypt(privateKey, kemCapsule, extendedNonce, ciphertext, associatedData); err == nil {
		t.Fatalf("Decrypt unexpectedly succeeded with overlong nonce")
	}
}

type errRoundTripMismatch string

func (e errRoundTripMismatch) Error() string { return string(e) }
