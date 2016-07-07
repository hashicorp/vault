package logical

import (
	"testing"
	"time"
)

func TestLeaseOptionsLeaseTotal(t *testing.T) {
	var l LeaseOptions
	l.TTL = 1 * time.Hour

	actual := l.LeaseTotal()
	expected := l.TTL
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsLeaseTotal_grace(t *testing.T) {
	var l LeaseOptions
	l.TTL = 1 * time.Hour

	actual := l.LeaseTotal()
	if actual != l.TTL {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsLeaseTotal_negLease(t *testing.T) {
	var l LeaseOptions
	l.TTL = -1 * 1 * time.Hour

	actual := l.LeaseTotal()
	expected := time.Duration(0)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime(t *testing.T) {
	var l LeaseOptions
	l.TTL = 1 * time.Hour

	limit := time.Now().Add(time.Hour)
	exp := l.ExpirationTime()
	if exp.Before(limit) {
		t.Fatalf("bad: %s", exp)
	}
}

func TestLeaseOptionsExpirationTime_noLease(t *testing.T) {
	var l LeaseOptions
	if !l.ExpirationTime().IsZero() {
		t.Fatal("should be zero")
	}
}
