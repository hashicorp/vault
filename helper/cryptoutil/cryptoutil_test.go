package cryptoutil

import "testing"

func TestBlake2b256Hash(t *testing.T) {
	hashVal, err = Blake2b256Hash("sampletext")
	if err != nil {
		t.Fatal(err)
	}

	if string(hashVal) == "" || string(hashVal) == "sampletext" {
		t.Fatalf("failed to hash the text")
	}
}
