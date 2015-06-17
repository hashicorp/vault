package logical

import (
	"testing"
	"time"
)

func TestLeaseOptionsLeaseTotal(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour

	actual := l.LeaseTotal()
	expected := l.Lease
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsLeaseTotal_grace(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour
	l.LeaseGracePeriod = 30 * time.Minute

	actual := l.LeaseTotal()
	expected := l.Lease + l.LeaseGracePeriod
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsLeaseTotal_negLease(t *testing.T) {
	var l LeaseOptions
	l.Lease = -1 * 1 * time.Hour
	l.LeaseGracePeriod = 30 * time.Minute

	actual := l.LeaseTotal()
	expected := time.Duration(0)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsLeaseTotal_negGrace(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour
	l.LeaseGracePeriod = -1 * 30 * time.Minute

	actual := l.LeaseTotal()
	expected := l.Lease
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour

	limit := time.Now().UTC().Add(time.Hour)
	exp := l.ExpirationTime()
	if exp.Before(limit) {
		t.Fatalf("bad: %s", exp)
	}
}

func TestLeaseOptionsExpirationTime_grace(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour
	l.LeaseGracePeriod = 30 * time.Minute

	limit := time.Now().UTC().Add(time.Hour + 30*time.Minute)
	actual := l.ExpirationTime()
	if actual.Before(limit) {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime_graceNegative(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour
	l.LeaseGracePeriod = -1 * 30 * time.Minute

	limit := time.Now().UTC().Add(time.Hour)
	actual := l.ExpirationTime()
	if actual.Before(limit) {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime_noLease(t *testing.T) {
	var l LeaseOptions
	if !l.ExpirationTime().IsZero() {
		t.Fatal("should be zero")
	}
}
