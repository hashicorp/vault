package kubernetes

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/serviceregistration/kubernetes/client"
	kubetest "github.com/hashicorp/vault/serviceregistration/kubernetes/testing"
)

func TestRetryHandlerSimple(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping because this test takes 10-15 seconds")
	}

	testState, testConf, closeFunc := kubetest.Server(t)
	defer closeFunc()

	client.Scheme = testConf.ClientScheme
	client.TokenFile = testConf.PathToTokenFile
	client.RootCAFile = testConf.PathToRootCAFile
	if err := os.Setenv(client.EnvVarKubernetesServiceHost, testConf.ServiceHost); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(client.EnvVarKubernetesServicePort, testConf.ServicePort); err != nil {
		t.Fatal(err)
	}

	logger := hclog.NewNullLogger()
	shutdownCh := make(chan struct{})
	wait := &sync.WaitGroup{}
	testPatch := &client.Patch{
		Operation: client.Add,
		Path:      "patch-path",
		Value:     "true",
	}

	c, err := client.New(logger, shutdownCh)
	if err != nil {
		t.Fatal(err)
	}

	r := &retryHandler{
		logger:         logger,
		namespace:      kubetest.ExpectedNamespace,
		podName:        kubetest.ExpectedPodName,
		client:         c,
		patchesToRetry: make(map[string]*client.Patch),
	}
	go r.Run(shutdownCh, wait)

	if testState.NumPatches() != 0 {
		t.Fatal("expected no current patches")
	}
	r.Add(testPatch)
	// Wait ample until the next try should have occurred.
	<-time.NewTimer(retryFreq * 2).C
	if testState.NumPatches() != 1 {
		t.Fatal("expected 1 patch")
	}
}

func TestRetryHandlerAdd(t *testing.T) {
	r := &retryHandler{
		logger:         hclog.NewNullLogger(),
		namespace:      "some-namespace",
		podName:        "some-pod-name",
		patchesToRetry: make(map[string]*client.Patch),
	}

	testPatch1 := &client.Patch{
		Operation: client.Add,
		Path:      "one",
		Value:     "true",
	}
	testPatch2 := &client.Patch{
		Operation: client.Add,
		Path:      "two",
		Value:     "true",
	}
	testPatch3 := &client.Patch{
		Operation: client.Add,
		Path:      "three",
		Value:     "true",
	}
	testPatch4 := &client.Patch{
		Operation: client.Add,
		Path:      "four",
		Value:     "true",
	}

	// Should be able to add all 4 patches.
	r.Add(testPatch1)
	if len(r.patchesToRetry) != 1 {
		t.Fatal("expected 1 patch")
	}

	r.Add(testPatch2)
	if len(r.patchesToRetry) != 2 {
		t.Fatal("expected 2 patches")
	}

	r.Add(testPatch3)
	if len(r.patchesToRetry) != 3 {
		t.Fatal("expected 3 patches")
	}

	r.Add(testPatch4)
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	// Adding a dupe should result in no change.
	r.Add(testPatch4)
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	// Adding a reversion should result in its twin being subtracted.
	r.Add(&client.Patch{
		Operation: client.Add,
		Path:      "four",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	r.Add(&client.Patch{
		Operation: client.Add,
		Path:      "three",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	r.Add(&client.Patch{
		Operation: client.Add,
		Path:      "two",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	r.Add(&client.Patch{
		Operation: client.Add,
		Path:      "one",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}
}

// This is meant to be run with the -race flag on.
func TestRetryHandlerRacesAndDeadlocks(t *testing.T) {
	_, testConf, closeFunc := kubetest.Server(t)
	defer closeFunc()

	client.Scheme = testConf.ClientScheme
	client.TokenFile = testConf.PathToTokenFile
	client.RootCAFile = testConf.PathToRootCAFile
	if err := os.Setenv(client.EnvVarKubernetesServiceHost, testConf.ServiceHost); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(client.EnvVarKubernetesServicePort, testConf.ServicePort); err != nil {
		t.Fatal(err)
	}

	logger := hclog.NewNullLogger()
	shutdownCh := make(chan struct{})
	wait := &sync.WaitGroup{}
	testPatch := &client.Patch{
		Operation: client.Add,
		Path:      "patch-path",
		Value:     "true",
	}

	c, err := client.New(logger, shutdownCh)
	if err != nil {
		t.Fatal(err)
	}

	r := &retryHandler{
		logger:         logger,
		namespace:      kubetest.ExpectedNamespace,
		podName:        kubetest.ExpectedPodName,
		client:         c,
		patchesToRetry: make(map[string]*client.Patch),
	}
	go r.Run(shutdownCh, wait)

	// Now hit it as quickly as possible to see if we can produce
	// races or deadlocks.
	start := make(chan struct{})
	done := make(chan bool)
	numRoutines := 100
	for i := 0; i < numRoutines; i++ {
		go func() {
			<-start
			r.Add(testPatch)
			done <- true
		}()
	}
	close(start)

	// Allow up to 5 seconds for everything to finish.
	timer := time.NewTimer(5 * time.Second)
	for i := 0; i < numRoutines; i++ {
		select {
		case <-timer.C:
			t.Fatal("test took too long to complete, check for deadlock")
		case <-done:
		}
	}
}
