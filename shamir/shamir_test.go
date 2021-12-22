package shamir

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

const MaxAttempts = 10
const MaxWaitInBetweenRetries = time.Second * 2

func TestSplit_invalid(t *testing.T) {
	secret := []byte("test")

	if _, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, secret, 0, 0); err == nil {
		t.Fatalf("expect error")
	}

	if _, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, secret, 2, 3); err == nil {
		t.Fatalf("expect error")
	}

	if _, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, secret, 1000, 3); err == nil {
		t.Fatalf("expect error")
	}

	if _, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, secret, 10, 1); err == nil {
		t.Fatalf("expect error")
	}

	if _, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, nil, 3, 2); err == nil {
		t.Fatalf("expect error")
	}
}

func TestSplit(t *testing.T) {
	secret := []byte("test")

	out, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, secret, 5, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(out) != 5 {
		t.Fatalf("bad: %v", out)
	}

	for _, share := range out {
		if len(share) != len(secret)+1 {
			t.Fatalf("bad: %v", out)
		}
	}
}

func TestCombine_invalid(t *testing.T) {
	// Not enough parts
	if _, err := Combine(nil); err == nil {
		t.Fatalf("should err")
	}

	// Mis-match in length
	parts := [][]byte{
		[]byte("foo"),
		[]byte("ba"),
	}
	if _, err := Combine(parts); err == nil {
		t.Fatalf("should err")
	}

	// Too short
	parts = [][]byte{
		[]byte("f"),
		[]byte("b"),
	}
	if _, err := Combine(parts); err == nil {
		t.Fatalf("should err")
	}

	parts = [][]byte{
		[]byte("foo"),
		[]byte("foo"),
	}
	if _, err := Combine(parts); err == nil {
		t.Fatalf("should err")
	}
}

func TestCombine(t *testing.T) {
	secret := []byte("test")

	out, err := retrySplit(MaxAttempts, MaxWaitInBetweenRetries, secret, 5, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// There is 5*4*3 possible choices,
	// we will just brute force try them all
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 5; k++ {
				if k == i || k == j {
					continue
				}
				parts := [][]byte{out[i], out[j], out[k]}
				recomb, err := Combine(parts)
				if err != nil {
					t.Fatalf("err: %v", err)
				}

				if !bytes.Equal(recomb, secret) {
					t.Errorf("parts: (i:%d, j:%d, k:%d) %v", i, j, k, parts)
					t.Fatalf("bad: %v %v", recomb, secret)
				}
			}
		}
	}
}

func TestField_Add(t *testing.T) {
	if out := add(16, 16); out != 0 {
		t.Fatalf("Bad: %v 16", out)
	}

	if out := add(3, 4); out != 7 {
		t.Fatalf("Bad: %v 7", out)
	}
}

func TestField_Mult(t *testing.T) {
	if out := mult(3, 7); out != 9 {
		t.Fatalf("Bad: %v 9", out)
	}

	if out := mult(3, 0); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}

	if out := mult(0, 3); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}
}

func TestField_Divide(t *testing.T) {
	if out := div(0, 7); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}

	if out := div(3, 3); out != 1 {
		t.Fatalf("Bad: %v 1", out)
	}

	if out := div(6, 3); out != 2 {
		t.Fatalf("Bad: %v 2", out)
	}
}

func TestPolynomial_Random(t *testing.T) {
	p, err := retryMakePolynomial(MaxAttempts, MaxWaitInBetweenRetries, 42, 2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.coefficients[0] != 42 {
		t.Fatalf("bad: %v", p.coefficients)
	}
}

func TestPolynomial_Eval(t *testing.T) {
	p, err := retryMakePolynomial(MaxAttempts, MaxWaitInBetweenRetries, 42, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if out := p.evaluate(0); out != 42 {
		t.Fatalf("bad: %v", out)
	}

	out := p.evaluate(1)
	exp := add(42, mult(1, p.coefficients[1]))
	if out != exp {
		t.Fatalf("bad: %v %v %v", out, exp, p.coefficients)
	}
}

func TestInterpolate_Rand(t *testing.T) {
	for i := 0; i < 256; i++ {
		p, err := retryMakePolynomial(MaxAttempts, MaxWaitInBetweenRetries, uint8(i), 2)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		x_vals := []uint8{1, 2, 3}
		y_vals := []uint8{p.evaluate(1), p.evaluate(2), p.evaluate(3)}
		out := interpolatePolynomial(x_vals, y_vals, 0)
		if out != uint8(i) {
			t.Fatalf("Bad: %v %d", out, i)
		}
	}
}

// retryMakePolynomial retries makePolynomial() for maxAttempts and sleeping for sleep duration.
// Need to handcode this for test to pass. Other solution is to import Retry packages from
// github.com/GoogleCloudPlatform/golang-samples/internal/testutil. These packages bring in a lot of external dependencies
// Hence handcoding our two retry functions for makePolynomial() and Split()
// go generics is badly needed here
func retryMakePolynomial(maxAttempts int, sleep time.Duration, intercept, degree uint8) (polynomial, error) {
	var p polynomial
	var err error
	for i := 0; i < maxAttempts; i++ {
		p, err = makePolynomial(intercept, degree)
		switch err {
		case ErrZeroShare, ErrNonUniqueShare:
			time.Sleep(sleep)
		case nil:
			return p, nil
		default:
			return p, err
		}
	}
	return p, fmt.Errorf("max retries %d exceed for", maxAttempts)
}

func retrySplit(maxAttempts int, sleep time.Duration, secret []byte, parts, threshold int) ([][]byte, error) {
	var splits [][]byte
	var err error
	for i := 0; i < maxAttempts; i++ {
		splits, err = Split(secret, parts, threshold)
		switch err {
		case ErrZeroShare, ErrNonUniqueShare:
			time.Sleep(sleep)
		case nil:
			return splits, nil
		default:
			return splits, err
		}
	}
	return splits, fmt.Errorf("max retries %d exceed for", maxAttempts)
}
