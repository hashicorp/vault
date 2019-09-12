package logical

import (
	"encoding/json"
	"testing"
)

func TestJSONSerialization(t *testing.T) {
	tt := TokenTypeDefaultBatch
	s, err := json.Marshal(tt)
	if err != nil {
		t.Fatal(err)
	}

	var utt TokenType
	err = json.Unmarshal(s, &utt)
	if err != nil {
		t.Fatal(err)
	}

	if tt != utt {
		t.Fatalf("expected %v, got %v", tt, utt)
	}

	utt = TokenTypeDefault
	err = json.Unmarshal([]byte(`"default-batch"`), &utt)
	if err != nil {
		t.Fatal(err)
	}
	if tt != utt {
		t.Fatalf("expected %v, got %v", tt, utt)
	}

	// Test on an empty value, which should unmarshal into TokenTypeDefault
	tt = TokenTypeDefault
	err = json.Unmarshal([]byte(`""`), &utt)
	if err != nil {
		t.Fatal(err)
	}
	if tt != utt {
		t.Fatalf("expected %v, got %v", tt, utt)
	}
}
