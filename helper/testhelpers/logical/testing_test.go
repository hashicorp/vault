package testing

import (
	"os"
	"testing"
)

func init() {
	testTesting = true

	if err := os.Setenv(TestEnvVar, "1"); err != nil {
		panic(err)
	}
}

func TestTest_noEnv(t *testing.T) {
	// Unset the variable
	if err := os.Setenv(TestEnvVar, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Setenv(TestEnvVar, "1")

	mt := new(mockT)
	Test(mt, TestCase{
		AcceptanceTest: true,
	})

	if !mt.SkipCalled {
		t.Fatal("skip not called")
	}
}

func TestTest_preCheck(t *testing.T) {
	called := false

	mt := new(mockT)
	Test(mt, TestCase{
		PreCheck: func() { called = true },
	})

	if !called {
		t.Fatal("precheck should be called")
	}
}

// mockT implements TestT for testing
type mockT struct {
	ErrorCalled bool
	ErrorArgs   []any
	FatalCalled bool
	FatalArgs   []any
	SkipCalled  bool
	SkipArgs    []any

	f bool
}

func (t *mockT) Error(args ...any) {
	t.ErrorCalled = true
	t.ErrorArgs = args
	t.f = true
}

func (t *mockT) Fatal(args ...any) {
	t.FatalCalled = true
	t.FatalArgs = args
	t.f = true
}

func (t *mockT) Skip(args ...any) {
	t.SkipCalled = true
	t.SkipArgs = args
	t.f = true
}

func (t *mockT) failed() bool {
	return t.f
}

func (t *mockT) failMessage() string {
	if t.FatalCalled {
		return t.FatalArgs[0].(string)
	} else if t.ErrorCalled {
		return t.ErrorArgs[0].(string)
	} else if t.SkipCalled {
		return t.SkipArgs[0].(string)
	}

	return "unknown"
}
