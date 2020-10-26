package awsutil

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	multierror "github.com/hashicorp/go-multierror"
)

func Test_CheckAWSError(t *testing.T) {
	testCases := []struct {
		Name     string
		Err      error
		Expected error
	}{
		{
			Name: "Something not checked",
			Err:  fmt.Errorf("something"),
		},
		{
			Name:     "Upstream throttle error",
			Err:      awserr.New("Throttling", "", nil),
			Expected: ErrUpstreamRateLimited,
		},
		{
			Name:     "Upstream RequestLimitExceeded",
			Err:      awserr.New("RequestLimitExceeded", "Request rate limited", nil),
			Expected: ErrUpstreamRateLimited,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := CheckAWSError(tc.Err)
			if err == nil && tc.Expected != nil {
				t.Fatalf("expected non-nil error (%#v), got nil", tc.Expected)
			}
			if err != nil && tc.Expected == nil {
				t.Fatalf("expected nil error, got (%#v)", err)
			}
			if err != tc.Expected {
				t.Fatalf("expected error (%#v), got (%#v)", tc.Expected, err)
			}
		})
	}
}

func Test_AppendRateLimitedError(t *testing.T) {
	awsErr := awserr.New("Throttling", "", nil)
	testCases := []struct {
		Name     string
		Err      error
		Expected error
	}{
		{
			Name:     "Something not checked",
			Err:      fmt.Errorf("something"),
			Expected: fmt.Errorf("something"),
		},
		{
			Name:     "Upstream throttle error",
			Err:      awsErr,
			Expected: multierror.Append(awsErr, ErrUpstreamRateLimited),
		},
		{
			Name: "Nil",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := AppendAWSError(tc.Err)
			if err == nil && tc.Expected != nil {
				t.Fatalf("expected non-nil error (%#v), got nil", tc.Expected)
			}
			if err != nil && tc.Expected == nil {
				t.Fatalf("expected nil error, got (%#v)", err)
			}
			if err == nil && tc.Expected == nil {
				return
			}
			if err.Error() != tc.Expected.Error() {
				t.Fatalf("expected error (%#v), got (%#v)", tc.Expected.Error(), err.Error())
			}
		})
	}
}
