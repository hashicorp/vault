package vault

import (
	"testing"
	"time"
)

func TestLease_Validate(t *testing.T) {
	l := &Lease{}
	if err := l.Validate(); err.Error() != "lease duration must be greater than zero" {
		t.Fatalf("err: %v", err)
	}

	l.Duration = time.Minute
	if err := l.Validate(); err.Error() != "maximum lease duration must be greater than zero" {
		t.Fatalf("err: %v", err)
	}

	l.MaxDuration = time.Second
	if err := l.Validate(); err.Error() != "lease duration cannot be greater than maximum lease duration" {
		t.Fatalf("err: %v", err)
	}

	l.MaxDuration = time.Minute
	l.MaxIncrement = -1 * time.Second
	if err := l.Validate(); err.Error() != "maximum lease increment cannot be negative" {
		t.Fatalf("err: %v", err)
	}

	l.MaxIncrement = time.Second
	if err := l.Validate(); err != nil {
		t.Fatalf("err: %v", err)
	}
}
