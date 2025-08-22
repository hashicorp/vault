package keysutil

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"
)

var kyberCases = []struct {
	name string
	kt   KeyType
}{
	{name: "kyber512", kt: KeyType_Kyber512},
	{name: "kyber768", kt: KeyType_Kyber768},
	{name: "kyber1024", kt: KeyType_Kyber1024},
}

func withEachKyberBox(t *testing.T, fn func(t *testing.T, box kyberBox)) {
	t.Helper()
	for _, c := range kyberCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			box, err := newKyberBox(c.kt)
			if err != nil {
				t.Fatalf("newKyberBox: %v", err)
			}
			fn(t, box)
		})
	}
}

func TestKyberBox_RoundTrip_WithAssociatedData(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		plaintext := []byte("hello, kyber + aes-gcm 🚀")
		ad := []byte("plaintext")
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, ad)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		pt, err := kyb.Decrypt(sk, capsule, nonce, ct, ad)
		if err != nil {
			t.Fatalf("Decrypt: %v", err)
		}
		if !bytes.Equal(pt, plaintext) {
			t.Fatalf("decrypted plaintext mismatch: got %q want %q", pt, plaintext)
		}
	})
}

func TestKyberBox_RoundTrip_EmptyPlaintext(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		var plaintext []byte
		ad := []byte("plaintext")
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, ad)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		pt, err := kyb.Decrypt(sk, capsule, nonce, ct, ad)
		if err != nil {
			t.Fatalf("Decrypt: %v", err)
		}
		if !bytes.Equal(pt, plaintext) {
			t.Fatalf("decrypted plaintext mismatch: got %q want %q", pt, plaintext)
		}
	})
}

func TestKyberBox_AssociatedDataMismatch_Fails(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		plaintext := []byte("secret")
		goodAD := []byte("plaintext")
		badAD := []byte("mismatch")
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, goodAD)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		if _, err := kyb.Decrypt(sk, capsule, nonce, ct, badAD); err == nil {
			t.Fatalf("Decrypt succeeded with mismatched associated data; expected failure")
		}
	})
}

func TestKyberBox_AssociatedData_NilEqualsEmpty(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		plaintext := []byte("nil vs empty AD")
		var adNil []byte
		adEmpty := []byte{}
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, adNil)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		pt, err := kyb.Decrypt(sk, capsule, nonce, ct, adEmpty)
		if err != nil {
			t.Fatalf("Decrypt: %v", err)
		}
		if !bytes.Equal(pt, plaintext) {
			t.Fatalf("decrypted plaintext mismatch: got %q want %q", pt, plaintext)
		}
	})
}

func TestKyberBox_InvalidNonceLength(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		plaintext := []byte("nonce length test")
		ad := []byte("plaintext")
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, ad)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		if len(nonce) == 0 {
			t.Fatalf("unexpected zero-length nonce")
		}
		trunc := nonce[:len(nonce)-1]
		if _, err := kyb.Decrypt(sk, capsule, trunc, ct, ad); err == nil {
			t.Fatalf("Decrypt succeeded with invalid nonce length; expected error")
		}
	})
}

func TestKyberBox_InvalidCapsuleLength(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		plaintext := []byte("capsule length test")
		ad := []byte("plaintext")
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, ad)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		if len(capsule) == 0 {
			t.Fatalf("unexpected zero-length capsule")
		}
		trunc := capsule[:len(capsule)-1]
		if _, err := kyb.Decrypt(sk, trunc, nonce, ct, ad); err == nil {
			t.Fatalf("Decrypt succeeded with invalid capsule length; expected error")
		}
	})
}

func TestKyberBox_TamperCiphertext_Fails(t *testing.T) {
	withEachKyberBox(t, func(t *testing.T, kyb kyberBox) {
		pk, sk, err := kyb.s.GenerateKeyPair()
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}
		plaintext := []byte("tamper check")
		ad := []byte("plaintext")
		capsule, nonce, ct, err := kyb.Encrypt(pk, plaintext, ad)
		if err != nil {
			t.Fatalf("Encrypt: %v", err)
		}
		if len(ct) == 0 {
			t.Fatalf("unexpected zero-length ciphertext")
		}
		ct[0] ^= 0x01
		if _, err := kyb.Decrypt(sk, capsule, nonce, ct, ad); err == nil {
			t.Fatalf("Decrypt succeeded on tampered ciphertext; expected error")
		}
	})
}

func TestKyberBox_WrongPrivateKey_Fails(t *testing.T) {
	cases := []struct {
		name string
		kt   KeyType
	}{
		{"kyber512", KeyType_Kyber512},
		{"kyber768", KeyType_Kyber768},
		{"kyber1024", KeyType_Kyber1024},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			box, err := newKyberBox(c.kt)
			if err != nil {
				t.Fatalf("newKyberBox: %v", err)
			}
			pk1, sk1, err := box.s.GenerateKeyPair()
			if err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}
			_, sk2, err := box.s.GenerateKeyPair()
			if err != nil {
				t.Fatalf("GenerateKeyPair(2): %v", err)
			}
			ad := []byte("plaintext")
			pt := []byte("wrong key")
			cap1, n1, ct1, err := box.Encrypt(pk1, pt, ad)
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}
			if _, err := box.Decrypt(sk2, cap1, n1, ct1, ad); err == nil {
				t.Fatalf("expected error with wrong private key")
			}
			if _, err := box.Decrypt(sk1, cap1, n1, ct1, ad); err != nil {
				t.Fatalf("decrypt with correct key failed: %v", err)
			}
		})
	}
}

func TestKyberBox_NonceRandomized(t *testing.T) {
	cases := []struct {
		name string
		kt   KeyType
	}{
		{"kyber512", KeyType_Kyber512},
		{"kyber768", KeyType_Kyber768},
		{"kyber1024", KeyType_Kyber1024},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			box, err := newKyberBox(c.kt)
			if err != nil {
				t.Fatalf("newKyberBox: %v", err)
			}
			pk, _, err := box.s.GenerateKeyPair()
			if err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}
			ad := []byte("plaintext")
			pt := []byte("same data")
			_, n1, ct1, err := box.Encrypt(pk, pt, ad)
			if err != nil {
				t.Fatalf("Encrypt(1): %v", err)
			}
			_, n2, ct2, err := box.Encrypt(pk, pt, ad)
			if err != nil {
				t.Fatalf("Encrypt(2): %v", err)
			}
			if bytes.Equal(n1, n2) {
				t.Fatalf("nonces should differ")
			}
			if bytes.Equal(ct1, ct2) {
				t.Fatalf("ciphertexts should differ")
			}
		})
	}
}

func TestKyberBox_LargePlaintext_RoundTrip(t *testing.T) {
	cases := []struct {
		name string
		kt   KeyType
	}{
		{"kyber512", KeyType_Kyber512},
		{"kyber768", KeyType_Kyber768},
		{"kyber1024", KeyType_Kyber1024},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			box, err := newKyberBox(c.kt)
			if err != nil {
				t.Fatalf("newKyberBox: %v", err)
			}
			pk, sk, err := box.s.GenerateKeyPair()
			if err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}
			buf := make([]byte, 4096)
			if _, err := rand.Read(buf); err != nil {
				t.Fatalf("rand: %v", err)
			}
			ad := []byte("plaintext")
			cap1, n1, ct1, err := box.Encrypt(pk, buf, ad)
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}
			out, err := box.Decrypt(sk, cap1, n1, ct1, ad)
			if err != nil {
				t.Fatalf("Decrypt: %v", err)
			}
			if !bytes.Equal(out, buf) {
				t.Fatalf("mismatch after round-trip")
			}
		})
	}
}

func TestNewKyberBox_InvalidType_Error(t *testing.T) {
	const invalid KeyType = 9999
	if _, err := newKyberBox(invalid); err == nil {
		t.Fatalf("expected error for invalid kyber type")
	}
}

func TestNewKyberBox_CiphertextSize_Order(t *testing.T) {
	b512, err := newKyberBox(KeyType_Kyber512)
	if err != nil {
		t.Fatalf("newKyberBox 512: %v", err)
	}
	b768, err := newKyberBox(KeyType_Kyber768)
	if err != nil {
		t.Fatalf("newKyberBox 768: %v", err)
	}
	b1024, err := newKyberBox(KeyType_Kyber1024)
	if err != nil {
		t.Fatalf("newKyberBox 1024: %v", err)
	}
	s1 := b512.s.CiphertextSize()
	s2 := b768.s.CiphertextSize()
	s3 := b1024.s.CiphertextSize()
	if !(s1 < s2 && s2 < s3) {
		t.Fatalf("unexpected ciphertext size order: %d, %d, %d", s1, s2, s3)
	}
}

func TestKyberBox_InvalidCapsuleTooLong(t *testing.T) {
	cases := []struct {
		name string
		kt   KeyType
	}{
		{"kyber512", KeyType_Kyber512},
		{"kyber768", KeyType_Kyber768},
		{"kyber1024", KeyType_Kyber1024},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			box, err := newKyberBox(c.kt)
			if err != nil {
				t.Fatalf("newKyberBox: %v", err)
			}
			pk, sk, err := box.s.GenerateKeyPair()
			if err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}
			ad := []byte("plaintext")
			pt := []byte("x")
			cap1, n1, ct1, err := box.Encrypt(pk, pt, ad)
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}
			badCap := append(cap1, 0x00)
			if _, err := box.Decrypt(sk, badCap, n1, ct1, ad); err == nil {
				t.Fatalf("expected error for too-long capsule")
			}
		})
	}
}

func TestKyberBox_Encapsulate_WrongScheme_Fails(t *testing.T) {
	cases := []struct {
		name  string
		encKT KeyType
		keyKT KeyType
	}{
		{"enc512_key768", KeyType_Kyber512, KeyType_Kyber768},
		{"enc512_key1024", KeyType_Kyber512, KeyType_Kyber1024},
		{"enc768_key512", KeyType_Kyber768, KeyType_Kyber512},
		{"enc768_key1024", KeyType_Kyber768, KeyType_Kyber1024},
		{"enc1024_key512", KeyType_Kyber1024, KeyType_Kyber512},
		{"enc1024_key768", KeyType_Kyber1024, KeyType_Kyber768},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			encBox, err := newKyberBox(c.encKT)
			if err != nil {
				t.Fatalf("newKyberBox(enc): %v", err)
			}
			keyBox, err := newKyberBox(c.keyKT)
			if err != nil {
				t.Fatalf("newKyberBox(key): %v", err)
			}

			pk, _, err := keyBox.s.GenerateKeyPair()
			if err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}

			_, _, _, err = encBox.Encrypt(pk, []byte("data"), []byte("plaintext"))
			if err == nil {
				t.Fatalf("expected encapsulate error with mismatched public key scheme")
			}
			if !strings.Contains(err.Error(), "encapsulate") {
				t.Fatalf("expected error to contain 'encapsulate', got: %v", err)
			}
		})
	}
}
