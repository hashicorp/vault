package hostutil

import (
	"testing"
)

func TestCollectHostInfo(t *testing.T) {
	info, err := CollectHostInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info.Timestamp.IsZero() {
		t.Fatal("expected non-zero Timestamp")
	}
	if info.CPU == nil {
		t.Fatal("expected non-nil CPU value")
	}
	if info.CPUTimes == nil {
		t.Fatal("expected non-nil CPUTimes value")
	}
	if info.Disk == nil {
		t.Fatal("expected non-nil Disk value")
	}
	if info.Host == nil {
		t.Fatal("expected non-nil Host value")
	}
	if info.Memory == nil {
		t.Fatal("expected non-nil Memory value")
	}
}
