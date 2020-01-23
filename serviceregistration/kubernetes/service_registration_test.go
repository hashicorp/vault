package kubernetes

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/serviceregistration/kubernetes/client"
)

var testVersion = "version 1"

func TestServiceRegistration(t *testing.T) {
	currentPatches, closeFunc := client.TestServer(t)
	defer closeFunc()

	if len(currentPatches) != 0 {
		t.Fatalf("expected 0 patches but have %d: %+v", len(currentPatches), currentPatches)
	}
	shutdownCh := make(chan struct{})
	config := map[string]string{
		"namespace": client.TestNamespace,
		"pod_name":  client.TestPodname,
	}
	logger := hclog.NewNullLogger()
	state := &sr.State{
		VaultVersion:         testVersion,
		IsInitialized:        true,
		IsSealed:             true,
		IsActive:             true,
		IsPerformanceStandby: true,
	}
	reg, err := NewServiceRegistration(shutdownCh, config, logger, state, "")
	if err != nil {
		t.Fatal(err)
	}

	// Test initial state.
	if len(currentPatches) != 5 {
		t.Fatalf("expected 5 current labels but have %d: %+v", len(currentPatches), currentPatches)
	}
	if currentPatches[pathToLabels+labelVaultVersion].Value != testVersion {
		t.Fatalf("expected %q but received %q", testVersion, currentPatches[pathToLabels+labelVaultVersion])
	}
	if currentPatches[pathToLabels+labelActive].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelVaultVersion])
	}
	if currentPatches[pathToLabels+labelSealed].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelVaultVersion])
	}
	if currentPatches[pathToLabels+labelPerfStandby].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelVaultVersion])
	}
	if currentPatches[pathToLabels+labelInitialized].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelVaultVersion])
	}

	// Test NotifyActiveStateChange.
	if err := reg.NotifyActiveStateChange(false); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelActive].Value != toString(false) {
		t.Fatalf("expected %q but received %q", toString(false), currentPatches[pathToLabels+labelActive])
	}
	if err := reg.NotifyActiveStateChange(true); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelActive].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelActive])
	}

	// Test NotifySealedStateChange.
	if err := reg.NotifySealedStateChange(false); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelSealed].Value != toString(false) {
		t.Fatalf("expected %q but received %q", toString(false), currentPatches[pathToLabels+labelSealed])
	}
	if err := reg.NotifySealedStateChange(true); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelSealed].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelSealed])
	}

	// Test NotifyPerformanceStandbyStateChange.
	if err := reg.NotifyPerformanceStandbyStateChange(false); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelPerfStandby].Value != toString(false) {
		t.Fatalf("expected %q but received %q", toString(false), currentPatches[pathToLabels+labelPerfStandby])
	}
	if err := reg.NotifyPerformanceStandbyStateChange(true); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelPerfStandby].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelPerfStandby])
	}

	// Test NotifyInitializedStateChange.
	if err := reg.NotifyInitializedStateChange(false); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelInitialized].Value != toString(false) {
		t.Fatalf("expected %q but received %q", toString(false), currentPatches[pathToLabels+labelInitialized])
	}
	if err := reg.NotifyInitializedStateChange(true); err != nil {
		t.Fatal(err)
	}
	if currentPatches[pathToLabels+labelInitialized].Value != toString(true) {
		t.Fatalf("expected %q but received %q", toString(true), currentPatches[pathToLabels+labelInitialized])
	}
}
