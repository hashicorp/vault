package cryptoutil

import "testing"

func TestBlake2b256Hash(t *testing.T) {
	hashVal := Blake2b256Hash("sampletext")

	if string(hashVal) == "" || string(hashVal) == "sampletext" {
		t.Fatalf("failed to hash the text")
	}
}

func TestHMACSHA256Hash(t *testing.T) {
	// Test on empty value
	if _, err := HMACSHA256Hash("", "bar"); err == nil {
		t.Fatal("expected error when an empty key is provided")
	}

	hashVal, err := HMACSHA256Hash("foo", "bar")
	if err != nil {
		t.Fatal(err)
	}

	if hashVal == "" || hashVal == "bar" {
		t.Fatalf("expected result to return a hashed value, got: %v", hashVal)
	}
}
