// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubernetes

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	sr "github.com/hashicorp/vault/serviceregistration"
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

	c, err := client.New(logger)
	if err != nil {
		t.Fatal(err)
	}

	r := &retryHandler{
		logger:         logger,
		namespace:      kubetest.ExpectedNamespace,
		podName:        kubetest.ExpectedPodName,
		patchesToRetry: make(map[string]*client.Patch),
		client:         c,
		initialState:   sr.State{},
	}
	r.Run(shutdownCh, wait)

	// Initial number of patches upon Run from setting the initial state
	initStatePatches := testState.NumPatches()
	if initStatePatches == 0 {
		t.Fatalf("expected number of states patches after initial patches to be non-zero")
	}

	// Send a new patch
	testPatch := &client.Patch{
		Operation: client.Add,
		Path:      "patch-path",
		Value:     "true",
	}
	r.Notify(testPatch)

	// Wait ample until the next try should have occurred.
	<-time.NewTimer(retryFreq * 2).C

	if testState.NumPatches() != initStatePatches+1 {
		t.Fatalf("expected 1 patch, got: %d", testState.NumPatches())
	}
}

func TestRetryHandlerAdd(t *testing.T) {
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
	c, err := client.New(logger)
	if err != nil {
		t.Fatal(err)
	}

	r := &retryHandler{
		logger:         hclog.NewNullLogger(),
		namespace:      "some-namespace",
		podName:        "some-pod-name",
		patchesToRetry: make(map[string]*client.Patch),
		client:         c,
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
	r.Notify(testPatch1)
	if len(r.patchesToRetry) != 1 {
		t.Fatal("expected 1 patch")
	}

	r.Notify(testPatch2)
	if len(r.patchesToRetry) != 2 {
		t.Fatal("expected 2 patches")
	}

	r.Notify(testPatch3)
	if len(r.patchesToRetry) != 3 {
		t.Fatal("expected 3 patches")
	}

	r.Notify(testPatch4)
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	// Adding a dupe should result in no change.
	r.Notify(testPatch4)
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	// Adding a reversion should result in its twin being subtracted.
	r.Notify(&client.Patch{
		Operation: client.Add,
		Path:      "four",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	r.Notify(&client.Patch{
		Operation: client.Add,
		Path:      "three",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	r.Notify(&client.Patch{
		Operation: client.Add,
		Path:      "two",
		Value:     "false",
	})
	if len(r.patchesToRetry) != 4 {
		t.Fatal("expected 4 patches")
	}

	r.Notify(&client.Patch{
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

	c, err := client.New(logger)
	if err != nil {
		t.Fatal(err)
	}

	r := &retryHandler{
		logger:         logger,
		namespace:      kubetest.ExpectedNamespace,
		podName:        kubetest.ExpectedPodName,
		patchesToRetry: make(map[string]*client.Patch),
		initialState:   sr.State{},
		client:         c,
	}

	// Now hit it as quickly as possible to see if we can produce
	// races or deadlocks.
	start := make(chan struct{})
	done := make(chan bool)
	numRoutines := 100
	for i := 0; i < numRoutines; i++ {
		go func() {
			<-start
			r.Notify(testPatch)
			done <- true
		}()
		go func() {
			<-start
			r.Run(shutdownCh, wait)
			done <- true
		}()
	}
	close(start)

	// Allow up to 5 seconds for everything to finish.
	timer := time.NewTimer(5 * time.Second)
	for i := 0; i < numRoutines*2; i++ {
		select {
		case <-timer.C:
			t.Fatal("test took too long to complete, check for deadlock")
		case <-done:
		}
	}
}

// In this test, the API server sends bad responses for 5 seconds,
// then sends good responses, and we make sure we get the expected behavior.
func TestRetryHandlerAPIConnectivityProblemsInitialState(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testState, testConf, closeFunc := kubetest.Server(t)
	defer closeFunc()
	kubetest.ReturnGatewayTimeouts.Store(true)

	client.Scheme = testConf.ClientScheme
	client.TokenFile = testConf.PathToTokenFile
	client.RootCAFile = testConf.PathToRootCAFile
	client.RetryMax = 0
	if err := os.Setenv(client.EnvVarKubernetesServiceHost, testConf.ServiceHost); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(client.EnvVarKubernetesServicePort, testConf.ServicePort); err != nil {
		t.Fatal(err)
	}

	shutdownCh := make(chan struct{})
	wait := &sync.WaitGroup{}
	reg, err := NewServiceRegistration(map[string]string{
		"namespace": kubetest.ExpectedNamespace,
		"pod_name":  kubetest.ExpectedPodName,
	}, hclog.NewNullLogger(), sr.State{
		VaultVersion:         "vault-version",
		IsInitialized:        true,
		IsSealed:             true,
		IsActive:             true,
		IsPerformanceStandby: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := reg.Run(shutdownCh, wait, ""); err != nil {
		t.Fatal(err)
	}

	// At this point, since the initial state can't be set,
	// remotely we should have false for all these labels.
	patch := testState.Get(pathToLabels + labelVaultVersion)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelActive)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelSealed)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelPerfStandby)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelInitialized)
	if patch != nil {
		t.Fatal("expected no value")
	}

	kubetest.ReturnGatewayTimeouts.Store(false)

	// Now we need to wait to give the retry handler
	// a chance to update these values.
	time.Sleep(retryFreq + time.Second)
	val := testState.Get(pathToLabels + labelVaultVersion)["value"]
	if val != "vault-version" {
		t.Fatal("expected vault-version")
	}
	val = testState.Get(pathToLabels + labelActive)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
	val = testState.Get(pathToLabels + labelSealed)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
	val = testState.Get(pathToLabels + labelPerfStandby)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
	val = testState.Get(pathToLabels + labelInitialized)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
}

// In this test, the API server sends bad responses for 5 seconds,
// then sends good responses, and we make sure we get the expected behavior.
func TestRetryHandlerAPIConnectivityProblemsNotifications(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testState, testConf, closeFunc := kubetest.Server(t)
	defer closeFunc()
	kubetest.ReturnGatewayTimeouts.Store(true)

	client.Scheme = testConf.ClientScheme
	client.TokenFile = testConf.PathToTokenFile
	client.RootCAFile = testConf.PathToRootCAFile
	client.RetryMax = 0
	if err := os.Setenv(client.EnvVarKubernetesServiceHost, testConf.ServiceHost); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(client.EnvVarKubernetesServicePort, testConf.ServicePort); err != nil {
		t.Fatal(err)
	}

	shutdownCh := make(chan struct{})
	wait := &sync.WaitGroup{}
	reg, err := NewServiceRegistration(map[string]string{
		"namespace": kubetest.ExpectedNamespace,
		"pod_name":  kubetest.ExpectedPodName,
	}, hclog.NewNullLogger(), sr.State{
		VaultVersion:         "vault-version",
		IsInitialized:        false,
		IsSealed:             false,
		IsActive:             false,
		IsPerformanceStandby: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := reg.NotifyActiveStateChange(true); err != nil {
		t.Fatal(err)
	}
	if err := reg.NotifyInitializedStateChange(true); err != nil {
		t.Fatal(err)
	}
	if err := reg.NotifyPerformanceStandbyStateChange(true); err != nil {
		t.Fatal(err)
	}
	if err := reg.NotifySealedStateChange(true); err != nil {
		t.Fatal(err)
	}

	if err := reg.Run(shutdownCh, wait, ""); err != nil {
		t.Fatal(err)
	}

	// At this point, since the initial state can't be set,
	// remotely we should have false for all these labels.
	patch := testState.Get(pathToLabels + labelVaultVersion)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelActive)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelSealed)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelPerfStandby)
	if patch != nil {
		t.Fatal("expected no value")
	}
	patch = testState.Get(pathToLabels + labelInitialized)
	if patch != nil {
		t.Fatal("expected no value")
	}

	kubetest.ReturnGatewayTimeouts.Store(false)

	// Now we need to wait to give the retry handler
	// a chance to update these values.
	time.Sleep(retryFreq + time.Second)

	// They should be "true" if the Notifications were set after the
	// initial state.
	val := testState.Get(pathToLabels + labelVaultVersion)["value"]
	if val != "vault-version" {
		t.Fatal("expected vault-version")
	}
	val = testState.Get(pathToLabels + labelActive)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
	val = testState.Get(pathToLabels + labelSealed)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
	val = testState.Get(pathToLabels + labelPerfStandby)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
	val = testState.Get(pathToLabels + labelInitialized)["value"]
	if val != "true" {
		t.Fatal("expected true")
	}
}
