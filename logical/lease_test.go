package logical

import (
	"testing"
	"time"
)

func TestLeaseOptionsIncrementedLease(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Second
	l.LeaseIssue = time.Now().UTC()

	actual := l.IncrementedLease(1 * time.Second)
	if actual > 3*time.Second || actual < 1*time.Second {
		t.Fatalf("bad: %s", actual)
	}
}

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
	l.LeaseIssue = time.Now().UTC()

	actual := l.ExpirationTime()
	expected := l.LeaseIssue.Add(l.Lease)
	if !actual.Equal(expected) {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime_grace(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour
	l.LeaseGracePeriod = 30 * time.Minute
	l.LeaseIssue = time.Now().UTC()

	actual := l.ExpirationTime()
	expected := l.LeaseIssue.Add(l.Lease + l.LeaseGracePeriod)
	if !actual.Equal(expected) {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime_graceNegative(t *testing.T) {
	var l LeaseOptions
	l.Lease = 1 * time.Hour
	l.LeaseGracePeriod = -1 * 30 * time.Minute
	l.LeaseIssue = time.Now().UTC()

	actual := l.ExpirationTime()
	expected := l.LeaseIssue.Add(l.Lease)
	if !actual.Equal(expected) {
		t.Fatalf("bad: %s", actual)
	}
}

func TestLeaseOptionsExpirationTime_noLease(t *testing.T) {
	var l LeaseOptions
	if !l.ExpirationTime().IsZero() {
		t.Fatal("should be zero")
	}
}
