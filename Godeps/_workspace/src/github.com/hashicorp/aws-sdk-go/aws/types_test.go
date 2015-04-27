package aws_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/aws-sdk-go/aws"
)

func TestUnixTimestampSerialization(t *testing.T) {
	d := time.Date(2014, 12, 20, 14, 55, 30, 500000000, time.UTC)
	ts := aws.UnixTimestamp{Time: d}
	out, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}

	if v, want := string(out), `1419087330.5`; v != want {
		t.Errorf("Was %q but expected %q", v, want)
	}
}

func TestUnixTimestampDeserialization(t *testing.T) {
	var ts aws.UnixTimestamp
	if err := json.Unmarshal([]byte(`1419087330.5`), &ts); err != nil {
		t.Fatal(err)
	}

	if v, want := ts.Time.Format(time.RFC3339Nano), "2014-12-20T14:55:30.5Z"; v != want {
		t.Errorf("Was %s but expected %s", v, want)
	}
}
