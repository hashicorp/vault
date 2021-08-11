package api

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
)

func TestRenewer_NewRenewer(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		i    *RenewerInput
		e    *Renewer
		err  bool
	}{
		{
			"nil",
			nil,
			nil,
			true,
		},
		{
			"missing_secret",
			&RenewerInput{
				Secret: nil,
			},
			nil,
			true,
		},
		{
			"default_grace",
			&RenewerInput{
				Secret: &Secret{},
			},
			&Renewer{
				secret: &Secret{},
			},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := client.NewRenewer(tc.i)
			if (err != nil) != tc.err {
				t.Fatal(err)
			}

			if v == nil {
				return
			}

			// Zero-out channels because reflect
			v.client = nil
			v.random = nil
			v.doneCh = nil
			v.renewCh = nil
			v.stopCh = nil

			if diff := deep.Equal(tc.e, v); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestLifetimeWatcher(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	// Note that doRenewWithOptions starts its loop with an initial renewal.
	// This has a big impact on the particulars of the following cases.

	renewedSecret := &Secret{}
	var caseOneErrorCount int
	var caseManyErrorsCount int
	cases := []struct {
		maxTestTime          time.Duration
		name                 string
		leaseDurationSeconds int
		incrementSeconds     int
		renew                renewFunc
		expectError          error
		expectRenewal        bool
	}{
		{
			time.Second,
			"no_error",
			60,
			60,
			func(_ string, _ int) (*Secret, error) {
				return renewedSecret, nil
			},
			nil,
			true,
		},
		{
			5 * time.Second,
			"one_error",
			15,
			15,
			func(_ string, _ int) (*Secret, error) {
				if caseOneErrorCount == 0 {
					caseOneErrorCount++
					return nil, fmt.Errorf("renew failure")
				}
				return renewedSecret, nil
			},
			nil,
			true,
		},
		{
			15 * time.Second,
			"many_errors",
			15,
			15,
			func(_ string, _ int) (*Secret, error) {
				if caseManyErrorsCount == 3 {
					return renewedSecret, nil
				}
				caseManyErrorsCount++
				return nil, fmt.Errorf("renew failure")
			},
			nil,
			true,
		},
		{
			15 * time.Second,
			"only_errors",
			15,
			15,
			func(_ string, _ int) (*Secret, error) {
				return nil, fmt.Errorf("renew failure")
			},
			nil,
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
				Secret: &Secret{
					LeaseDuration: tc.leaseDurationSeconds,
				},
				Increment: tc.incrementSeconds,
			})
			if err != nil {
				t.Fatal(err)
			}

			go func() {
				v.doneCh <- v.doRenewWithOptions(false, false, tc.leaseDurationSeconds, "myleaseID", tc.renew, time.Second)
			}()
			defer v.Stop()

			select {
			case <-time.After(tc.maxTestTime):
				t.Fatalf("renewal didn't happen")
			case r := <-v.RenewCh():
				if !tc.expectRenewal {
					t.Fatal("expected no renewals")
				}
				if r.Secret != renewedSecret {
					t.Fatalf("expected secret %v, got %v", renewedSecret, r.Secret)
				}
			case err := <-v.DoneCh():
				if tc.expectError != nil && !errors.Is(err, tc.expectError) {
					t.Fatalf("expected error %q, got: %v", tc.expectError, err)
				}
				if tc.expectRenewal {
					t.Fatal("expected at least one renewal")
				}
			}
		})
	}
}
