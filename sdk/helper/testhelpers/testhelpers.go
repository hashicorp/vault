package testhelpers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// RandomWithPrefix is used to generate a unique name with a prefix, for
// randomizing names in acceptance tests
func RandomWithPrefix(name string) string {
	return fmt.Sprintf("%s-%d", name, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

// RetryUntilAtCadence runs f until it returns a nil result or the timeout is reached.
// If a nil result hasn't been obtained by timeout, calls t.Fatal.
func RetryUntilAtCadence(t testing.TB, timeout, sleepTime time.Duration, f func() error) {
	t.Helper()
	fail := func(err error) {
		t.Helper()
		t.Fatalf("did not complete before deadline, err: %v", err)
	}
	RetryUntilAtCadenceWithHandler(t, timeout, sleepTime, fail, f)
}

// RetryUntilAtCadenceWithHandler runs f until it returns a nil result or the timeout is reached.
// If a nil result hasn't been obtained by timeout, onFailure is called.
func RetryUntilAtCadenceWithHandler(t testing.TB, timeout, sleepTime time.Duration, onFailure func(error), f func() error) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var err error
	for time.Now().Before(deadline) {
		if err = f(); err == nil {
			return
		}
		time.Sleep(sleepTime)
	}
	onFailure(err)
}

// RetryUntil runs f with a 100ms pause between calls, until f returns a nil result
// or the timeout is reached.
// If a nil result hasn't been obtained by timeout, calls t.Fatal.
// NOTE: See RetryUntilAtCadence if you want to specify a different wait/sleep
// duration between calls.
func RetryUntil(t testing.TB, timeout time.Duration, f func() error) {
	t.Helper()
	RetryUntilAtCadence(t, timeout, 100*time.Millisecond, f)
}
